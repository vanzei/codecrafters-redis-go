package builtin

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

type StreamID struct {
	Ms  uint64
	Seq uint64
}

type StreamEntry struct {
	ID     StreamID
	Fields map[string]string
}

type Stream struct {
	Entries []StreamEntry
	LastID  StreamID
}

func HandleXadd(args []string) (string, error) {
	if len(args) < 4 || !strings.EqualFold(args[0], "xadd") {
		return "", fmt.Errorf("ERR wrong number of arguments for 'xadd'")
	}
	key := args[1]

	trim, idTok, fields, err := parseXAddArgs(args[2:])
	if err != nil {
		return "", err
	}

	stream, err := getOrCreateStream(key, true)
	if err != nil {
		return "", err
	}

	id, err := nextID(stream, idTok)
	if err != nil {
		return "", err
	}

	appendEntry(stream, id, fields)
	applyTrim(stream, trim)

	//return fmt.Sprintf("%d-%d", id.Ms, id.Seq), nil
	idStr := fmt.Sprintf("%d-%d", id.Ms, id.Seq)
	return fmt.Sprintf("$%d\r\n%s\r\n", len(idStr), idStr), nil
}

func getOrCreateStream(key string, create bool) (*Stream, error) {
	if v, ok := database[key]; ok {
		if v.Type != "stream" || v.Stream == nil {
			return nil, fmt.Errorf("WRONGTYPE Operation against the key holding thewrong kind of value")
		}
		return v.Stream, nil
	}
	if !create {
		return nil, fmt.Errorf("No stream created (NOMKSTREAM)")
	}

	s := &Stream{LastID: StreamID{Ms: 0, Seq: 0}}
	database[key] = Value{Type: "stream", Stream: s}
	return s, nil
}

type trimSpec struct {
	kind      string
	approx    bool
	limit     int
	threshold uint64
	threshSeq uint64
}

func parseXAddArgs(args []string) (trim trimSpec, idToken string, fields map[string]string, err error) {
	i := 0
	if i+1 < len(args) && (strings.EqualFold(args[i], "maxlen") || strings.EqualFold(args[i], "minid")) {
		trim.kind = strings.ToLower(args[i])
		i++
		// ~ flag: mark approximate trimming
		if i < len(args) && args[i] == "~" {
			trim.approx = true
			i++
		}
		if i >= len(args) {
			return trim, "", nil, fmt.Errorf("syntax error")
		}
		if trim.kind == "minid" {
			parts := strings.Split(args[i], "-")
			if len(parts) != 2 {
				return trim, "", nil, fmt.Errorf("syntax error")
			}
			trim.threshold, err = strconv.ParseUint(parts[0], 10, 64)
			if err != nil {
				return trim, "", nil, err
			}
			trim.threshSeq, err = strconv.ParseUint(parts[1], 10, 64)
			if err != nil {
				return trim, "", nil, err
			}
		} else {
			trim.threshold, err = strconv.ParseUint(args[i], 10, 64)
			if err != nil {
				return trim, "", nil, err
			}
		}
		i++
		if i < len(args) && strings.EqualFold(args[i], "limit") {
			i++
			if i >= len(args) {
				return trim, "", nil, fmt.Errorf("syntax error")
			}
			trim.limit, err = strconv.Atoi(args[i])
			if err != nil {
				return trim, "", nil, err
			}
			i++
		}
	}
	if i >= len(args) {
		return trim, "", nil, fmt.Errorf("missing id")
	}
	idToken = args[i]
	i++

	if (len(args)-i)%2 != 0 {
		return trim, "", nil, fmt.Errorf("wrong number of fields")
	}
	fields = map[string]string{}
	for i < len(args) {
		fields[args[i]] = args[i+1]
		i += 2
	}
	return
}

func parseExplicitID(tok string) (StreamID, error) {

	parts := strings.Split(tok, "-")
	if len(parts) != 2 {
		return StreamID{}, fmt.Errorf("ERR invalid stream id specified as stream command argument")
	}
	ms, err := strconv.ParseUint(parts[0], 10, 64)
	if err != nil {
		return StreamID{}, err
	}
	seq, err := strconv.ParseUint(parts[1], 10, 64)
	if err != nil {
		return StreamID{}, err
	}
	return StreamID{Ms: ms, Seq: seq}, nil
}

func compareID(a, b StreamID) int {
	if a.Ms < b.Ms {
		return -1
	} else if a.Ms > b.Ms {
		return 1
	}

	if a.Seq < b.Seq {
		return -1
	} else if a.Seq > b.Seq {
		return 1
	}
	return 0

}

func nextID(stream *Stream, token string) (StreamID, error) {
	last := stream.LastID
	switch token {
	case "*":
		now := uint64(time.Now().UnixNano() / 1e6)
		id := StreamID{Ms: now, Seq: 0}
		if now == last.Ms && (last.Ms != 0 || last.Seq != 0) {
			id.Seq = last.Seq + 1
		}
		if compareID(id, last) <= 0 {
			id.Seq = last.Seq + 1
		}
		return id, nil
	default:
		id, err := parseExplicitID(token)
		if err != nil {
			return StreamID{}, err
		}
		// Reject 0-0 explicitly; stream IDs must be greater than 0-0 even on an empty stream.
		if id.Ms == 0 && id.Seq == 0 {
			return StreamID{}, fmt.Errorf("ERR The ID specified in XADD must be greater than 0-0")
		}
		if compareID(id, last) <= 0 {
			return StreamID{}, fmt.Errorf("ERR The ID specified in XADD is equal or smaller than the target stream top item")
		}
		return id, nil
	}
}
func appendEntry(stream *Stream, id StreamID, fields map[string]string) {
	stream.Entries = append(stream.Entries, StreamEntry{ID: id, Fields: fields})
	stream.LastID = id
}

// applyTrim honors trim.approx (~): trims in batches instead of exact sizing
func applyTrim(stream *Stream, spec trimSpec) {
	if spec.kind == "" {
		return
	}

	switch spec.kind {
	case "maxlen":
		if spec.threshold == 0 {
			return
		}
		// If ~ given, only trim when we exceed threshold + slack, then drop a batch
		if spec.approx {
			slack := 32 // small slack window; tweak as desired
			if len(stream.Entries) <= int(spec.threshold)+slack {
				return
			}
			cut := len(stream.Entries) - int(spec.threshold)
			// limit cap if provided
			if spec.limit > 0 && cut > spec.limit {
				cut = spec.limit
			}
			// drop cut items from head
			stream.Entries = stream.Entries[cut:]
			return
		}
		// exact trim: ensure length <= threshold, respecting limit
		if len(stream.Entries) > int(spec.threshold) {
			cut := len(stream.Entries) - int(spec.threshold)
			if spec.limit > 0 && cut > spec.limit {
				cut = spec.limit
			}
			stream.Entries = stream.Entries[cut:]
		}

	case "minid":
		threshold := StreamID{Ms: spec.threshold, Seq: spec.threshSeq}
		if spec.approx {
			// Approximate: drop in batches from the head while head < threshold
			batch := spec.limit
			if batch <= 0 {
				batch = 64 // default batch size
			}
			for len(stream.Entries) > 0 {
				head := stream.Entries[0].ID
				if compareID(head, threshold) >= 0 {
					break
				}
				// drop up to batch items
				drop := batch
				if drop > len(stream.Entries) {
					drop = len(stream.Entries)
				}
				stream.Entries = stream.Entries[drop:]
			}
			return
		}
		// exact: keep entries >= threshold
		keep := stream.Entries[:0]
		for _, e := range stream.Entries {
			if compareID(e.ID, threshold) >= 0 {
				keep = append(keep, e)
			}
		}
		stream.Entries = keep
	}
}

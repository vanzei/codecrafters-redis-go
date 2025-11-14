package builtin

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

func HandleSet(elements []string) (string, error) {
	if len(elements) < 3 || strings.ToUpper(elements[0]) != "SET" {
		return "", fmt.Errorf("Invalid SET command")
	}

	key, value := elements[1], elements[2]
	database[key] = value

	delete(expiry, key)

	if len(elements) > 3 {
		if len(elements) != 5 {
			return "", fmt.Errorf("Invalid SET arguments")
		}
		option := strings.ToUpper(elements[3])
		ttlSting := elements[4]

		ttl, err := strconv.ParseInt(ttlSting, 10, 64)
		if err != nil || ttl <= 0 {
			return "", fmt.Errorf("TTL provide not a number")
		}
		switch option {
		case "EX":
			expiry[key] = time.Now().Add(time.Duration(ttl) * time.Second)
		case "PX":
			expiry[key] = time.Now().Add(time.Duration(ttl) * time.Millisecond)
		default:
			return "", fmt.Errorf("Unknown SET option")

		}

	}
	return "+OK\r\n", nil

}

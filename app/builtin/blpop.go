package builtin

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

func HandleBlpop(elements []string) (string, error) {
	if len(elements) < 3 || strings.ToUpper(elements[0]) != "BLPOP" {
		return "", fmt.Errorf("Invalid BLPOP command")
	}

	lists := elements[1 : len(elements)-1]
	timeoutStr := elements[len(elements)-1]
	// timeoutSec, err := strconv.Atoi(timeoutStr)
	timeoutSec, err := strconv.ParseFloat(timeoutStr, 64)
	if err != nil || timeoutSec < 0 {
		return "", fmt.Errorf("Invalid timeout")
	}

	//Immediate pop path

	for _, key := range lists {
		val := database[key]
		if len(val.List) > 0 {
			resp, err := HandleLpop([]string{"LPOP", key})
			if err != nil {
				return "", err
			}
			return fmt.Sprintf("*2\r\n$%d\r\n%s\r\n%s", len(key), key, resp), nil

		}
	}

	var timer *time.Timer
	var timeout <-chan time.Time
	if timeoutSec > 0 {
		timer = time.NewTimer(time.Duration(timeoutSec * float64(time.Second)))
		timeout = timer.C
	}

	req := &blockRequest{
		clientID: uuid.NewString(),
		lists:    lists,
		result:   make(chan string, 1),
		timeout:  timer,
		// timeout:  time.NewTimer(time.Duration(timeoutSec) * time.Second),
	}

	addWaiter(req)
	defer removeWaiter(req)

	select {
	case resp := <-req.result:
		return resp, nil
	case <-timeout:
		return "*-1\r\n", nil
	}

}

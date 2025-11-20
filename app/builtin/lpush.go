package builtin

import (
	"fmt"
	"strings"
)

func wakeFirstWaiter(key string) bool {
	queue := waiters[key]
	if len(queue) == 0 {
		return false
	}

	req := queue[0]
	removeWaiter(req)

	resp, _ := HandleLpop([]string{"LPOP", key})
	select {
	case req.result <- fmt.Sprintf("*2\r\n$%d\r\n%s\r\n%s", len(key), key, resp):
	default:
	}
	if req.timeout != nil {
		req.timeout.Stop()
	} else {
		// req.timeout.Stop()
	}
	return true

}

func HandleLpush(elements []string) (string, error) {
	if len(elements) < 3 || strings.ToUpper(elements[0]) != "LPUSH" {
		return "", fmt.Errorf("Invalid LPUSH command")
	}

	key := elements[1]
	val := database[key]
	if val.Type != "" && val.Type != "list" {
		return "", fmt.Errorf("WRONGTYPE Operation")
	}

	val.Type = "list"
	list := val.List
	for i := 2; i < len(elements); i++ {
		fmt.Println(elements[i], list)
		list = append([]string{elements[i]}, list...)
	}
	val.List = list

	database[key] = val
	wakeFirstWaiter(key)
	return fmt.Sprintf(":%d\r\n", len(val.List)), nil
}

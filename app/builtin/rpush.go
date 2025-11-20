package builtin

import (
	"fmt"
	"strings"
)

func HandleRpush(elements []string) (string, error) {
	if len(elements) < 3 || strings.ToUpper(elements[0]) != "RPUSH" {
		return "", fmt.Errorf("Invalid RPUSH command")
	}

	key := elements[1]
	val := database[key]
	if val.Type != "" && val.Type != "list" {
		return "", fmt.Errorf("WRONGTYPE Operation")
	}

	val.Type = "list"
	val.List = append(val.List, elements[2:]...)

	database[key] = val
	wakeFirstWaiter(key)
	return fmt.Sprintf(":%d\r\n", len(val.List)), nil
}

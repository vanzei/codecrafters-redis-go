package builtin

import (
	"fmt"
	"strings"
)

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
	return fmt.Sprintf(":%d\r\n", len(val.List)), nil
}

package builtin

import (
	"fmt"
	"strings"
)

func HandleGet(elements []string) (string, error) {
	if len(elements) != 2 || strings.ToUpper(elements[0]) != "GET" {
		return "", fmt.Errorf("Invalid GET command")
	}
	key := elements[1]
	val, ok := database[key]
	if !ok {
		return "$-1\r\n", nil
	}
	return fmt.Sprintf("$%d\r\n%s\r\n", len(val), val), nil
}

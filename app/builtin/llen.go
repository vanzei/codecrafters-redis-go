package builtin

import (
	"fmt"
	"strings"
)

func HandleLlen(elements []string) (string, error) {
	if len(elements) != 2 || strings.ToUpper(elements[0]) != "LLEN" {
		return "", fmt.Errorf("Invalid LLEN command")
	}
	key := elements[1]

	// check expiration

	val, ok := database[key]
	if !ok {
		return ":0\r\n", nil
	}
	result := len(val.List)
	return fmt.Sprintf(":%d\r\n", result), nil
}

package builtin

import (
	"fmt"
	"strings"
)

func HandlePing(elements []string) (string, error) {
	if len(elements) != 1 || strings.ToUpper(elements[0]) != "PING" {
		return "", fmt.Errorf("Invalid Ping format")
	}
	return "+PONG\r\n", nil
}

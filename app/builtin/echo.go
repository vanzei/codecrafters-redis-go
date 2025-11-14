package builtin

import (
	"fmt"
	"strings"
)

// Handle ECHO
// Expects elements like ["ECHO", "arg"]
// return RESP-formated response

func HandleEcho(elements []string) (string, error) {
	if len(elements) != 2 || strings.ToUpper(elements[0]) != "ECHO" {
		return "", fmt.Errorf("Invalid ECHO command")
	}
	message := elements[1]

	return fmt.Sprintf("$%d\r\n%s\r\n", len(message), message), nil
}

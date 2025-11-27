package builtin

import (
	"fmt"
	"strings"
	"time"
)

func HandleType(elements []string) (string, error) {
	if len(elements) != 2 || strings.ToUpper(elements[0]) != "TYPE" {
		return "", fmt.Errorf("Invalid TYPE command")
	}
	key := elements[1]

	// expire if needed
	if exp, ok := expiry[key]; ok && time.Now().After(exp) {
		delete(database, key)
		delete(expiry, key)
		return "+none\r\n", nil
	}

	val, ok := database[key]
	if !ok || val.Type == "" {
		return "+none\r\n", nil
	}
	return "+" + val.Type + "\r\n", nil
}

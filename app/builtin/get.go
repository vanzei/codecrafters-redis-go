package builtin

import (
	"fmt"
	"strings"
	"time"
)

func HandleGet(elements []string) (string, error) {
	if len(elements) != 2 || strings.ToUpper(elements[0]) != "GET" {
		return "", fmt.Errorf("Invalid GET command")
	}
	key := elements[1]

	// check expiration

	if exp, ok := expiry[key]; ok {
		if time.Now().After(exp) {
			delete(database, key)
			delete(expiry, key)
			return "$-1\r\n", nil
		}
	}

	val, ok := database[key]
	if !ok {
		return "$-1\r\n", nil
	}
	valStr := val.Str
	if val.Type != "string" {
		return "", fmt.Errorf("WRONGTYPE, not a string")
	}
	return fmt.Sprintf("$%d\r\n%s\r\n", len(valStr), valStr), nil
}

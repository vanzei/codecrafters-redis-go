package builtin

import (
	"fmt"
	"strings"
)

func HandleLpop(elements []string) (string, error) {
	if len(elements) < 2 || strings.ToUpper(elements[0]) != "LPOP" {
		return "", fmt.Errorf("Invalid LPOP command")
	}

	key := elements[1]

	values := database[key].List
	value := values[0]

	delete(expiry, value)

	tempVal := database[key]
	tempVal.List = tempVal.List[1:len(tempVal.List)]

	database[key] = tempVal
	return fmt.Sprintf("$%d\r\n%s\r\n", len(value), value), nil
}

package builtin

import (
	"fmt"
	"strconv"
	"strings"
)

func HandleLpop(elements []string) (string, error) {
	if len(elements) < 2 || strings.ToUpper(elements[0]) != "LPOP" {
		return "", fmt.Errorf("Invalid LPOP command")
	}

	key := elements[1]
	val := database[key]
	if val.Type != "" && val.Type != "list" {
		return "", fmt.Errorf("WRONGTYPE Operation")
	}

	if len(val.List) == 0 {
		return "$-1\r\n", nil
	}

	count := 1
	if len(elements) == 3 {
		n, err := strconv.Atoi(elements[2])
		if err != nil {
			return "", fmt.Errorf("Not valid number to LPOP")
		}
		count = n

	}

	// Clamping to len list to pop
	if count > len(val.List) {
		count = len(val.List)
	}
	popped := val.List[:count]
	val.List = val.List[count:]

	if len(val.List) == 0 {
		val.Type = ""
	}
	database[key] = val

	if len(popped) == 1 {
		delete(expiry, key)
	}

	var b strings.Builder

	if len(elements) == 3 {
		b.WriteString(fmt.Sprintf("*%d\r\n", len(popped)))
	}
	for _, v := range popped {
		b.WriteString(fmt.Sprintf("$%d\r\n%s\r\n", len(v), v))
	}
	return b.String(), nil
}

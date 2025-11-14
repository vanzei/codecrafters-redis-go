package builtin

import (
	"fmt"
	"strings"
)

func HandleSet(elements []string) (string, error) {
	if len(elements) < 3 || strings.ToUpper(elements[0]) != "SET" {
		return "", fmt.Errorf("Invalid SET command")
	}
	key, value := elements[1], elements[2]
	database[key] = value
	return "+OK\r\n", nil

}

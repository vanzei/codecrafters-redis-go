package builtin

import (
	"fmt"
	"strconv"
	"strings"
)

func HandleLrange(elements []string) (string, error) {
	if len(elements) != 4 || strings.ToUpper(elements[0]) != "LRANGE" {
		return "", fmt.Errorf("Invalid LRANGE command")
	}
	listKey := elements[1]
	dbK := database[listKey]
	startIdx, err := strconv.Atoi(elements[2])
	if err != nil {
		return "", fmt.Errorf("Start Index not a number")
	}
	endIdx, err := strconv.Atoi(elements[3])
	if err != nil {
		return "", fmt.Errorf("End Index not a number")
	}
	if (endIdx) > len(dbK.List) {
		endIdx = len(dbK.List) - 1
	}

	result := dbK.List[startIdx : endIdx+1]

	var b strings.Builder
	b.WriteString(fmt.Sprintf("*%d\r\n", len(result)))
	for _, v := range result {
		b.WriteString(fmt.Sprintf("$%d\r\n%s\r\n", len(v), v))
	}
	return b.String(), nil

}

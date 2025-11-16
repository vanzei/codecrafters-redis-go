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
	if startIdx < 0 {
		startIdx += len(dbK.List)
	}
	if startIdx < 0 {
		startIdx = 0
	}
	endIdx, err := strconv.Atoi(elements[3])
	if err != nil {
		return "", fmt.Errorf("End Index not a number")
	}
	if endIdx < 0 {
		endIdx += len(dbK.List)
	}
	if (endIdx) > len(dbK.List) {
		endIdx = len(dbK.List) - 1
	}

	if startIdx > endIdx || startIdx >= len(dbK.List) {
		return "*0\r\n", nil
	}
	fmt.Println(startIdx, endIdx)

	result := dbK.List[startIdx : endIdx+1]

	var b strings.Builder
	b.WriteString(fmt.Sprintf("*%d\r\n", len(result)))
	for _, v := range result {
		b.WriteString(fmt.Sprintf("$%d\r\n%s\r\n", len(v), v))
	}
	return b.String(), nil

}

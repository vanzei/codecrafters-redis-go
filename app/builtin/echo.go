package builtin

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

func ParseRESPArray(r io.Reader) ([]string, error) {
	reader := bufio.NewReader(r)
	firstByte, err := reader.ReadByte()
	if err != nil || firstByte != '*' {
		return nil, fmt.Errorf("invalid RESP array start")
	}

	// length of the array
	arrayLen, err := reader.ReadString('\n')
	if err != nil {
		return nil, err
	}

	arrayLenString := strings.TrimSuffix(arrayLen, "\r\n")
	lenght, err := strconv.Atoi(arrayLenString)
	if err != nil {
		return nil, err
	}

	// Parse each element (assuming bulk strings for simplicity)
	var elements []string
	for i := 0; i < lenght; i++ {
		element, err := parseBulkString(reader)
		if err != nil {
			return nil, err
		}
		elements = append(elements, element)

	}
	return elements, nil

}

// parseBulkString parses a RESP bulk string (e.g., $5\r\nhello\r\n)
func parseBulkString(r *bufio.Reader) (string, error) {
	// read $
	dolar, err := r.ReadByte()
	if err != nil || dolar != '$' {
		return "", fmt.Errorf("missing $, invalid bulk string")
	}
	// length of the string
	strlen, err := r.ReadString('\n')
	if err != nil {
		return "", err
	}

	strlen = strings.TrimSuffix(strlen, "\r\n")
	lenght, err := strconv.Atoi(strlen)
	if err != nil {
		return "", err
	}

	data := make([]byte, lenght)
	_, err = io.ReadFull(r, data)
	if err != nil {
		return "", err
	}

	//read final
	crlf := make([]byte, 2)
	_, err = io.ReadFull(r, crlf)
	if err != nil || string(crlf) != "\r\n" {
		return "", fmt.Errorf("Invalid bulk string end")
	}

	return string(data), nil

}

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

package main

import (
	"fmt"
	"strings"

	builtinPck "github.com/codecrafters-io/redis-starter-go/app/builtin"
)

type CommandHandler func([]string) (string, error)

var SupportedCommands = map[string]CommandHandler{
	"PING":   builtinPck.HandlePing,
	"ECHO":   builtinPck.HandleEcho,
	"GET":    builtinPck.HandleGet,
	"SET":    builtinPck.HandleSet,
	"RPUSH":  builtinPck.HandleRpush,
	"LPUSH":  builtinPck.HandleLpush,
	"LRANGE": builtinPck.HandleLrange,
	"LLEN":   builtinPck.HandleLlen,
	"LPOP":   builtinPck.HandleLpop,
	"BLPOP":  builtinPck.HandleBlpop,
	"TYPE":   builtinPck.HandleType,
	"XADD":   builtinPck.HandleXadd,
}

func ProcessCommand(elements []string) (string, error) {
	if len(elements) == 0 {
		return "", fmt.Errorf("Empty command")
	}

	command := strings.ToUpper(elements[0])
	handler, exists := SupportedCommands[command]
	if !exists {
		return "-ERR unknown command\r\n", nil
	}
	return handler(elements)
}

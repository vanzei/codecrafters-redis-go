package builtin

import "time"

type Value struct {
	Type string
	Str  string
	List []string
}

var database = make(map[string]Value)
var expiry = make(map[string]time.Time)

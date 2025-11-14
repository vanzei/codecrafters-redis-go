package builtin

import "time"

var database = make(map[string]string)

var expiry = make(map[string]time.Time)

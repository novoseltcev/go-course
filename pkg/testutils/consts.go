package testutils

import (
	"errors"
)

const (
	STRING  = "string"
	INT     = 10
	FLAOT   = 10.123
	JSON    = `{"ping": "pong"}`
	URL     = "http://localhost"
	UNKNOWN = "unknown"
)

// nolint: gochecknoglobals
var (
	Err   = errors.New("test error")
	Bytes = []byte{1, 2, 3}
)

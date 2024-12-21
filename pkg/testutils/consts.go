package testutils

import (
	"errors"
)

const (
	STRING = "string"
	INT    = 10
	FLAOT  = 10.123
	JSON   = `{"ping": "pong"}`
	URL    = "http://localhost"
)

var Err = errors.New("test error")

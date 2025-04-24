package core

import "errors"

var ErrBadArguments = errors.New("arguments are not acceptable")
var ErrAlreadyRunning = errors.New("already running")

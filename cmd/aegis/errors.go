package main

import "errors"

// ErrArgumentCount is called when a cmdFunc is called with an invalid argc.
var ErrArgumentCount = errors.New("invalid argument count")

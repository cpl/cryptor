package noise

import "errors"

var (
	// ErrBadHandshakeState ...
	ErrBadHandshakeState = errors.New("bad handshake state")

	// ErrBadHandshakeRole ...
	ErrBadHandshakeRole = errors.New("bad handshake role")
)

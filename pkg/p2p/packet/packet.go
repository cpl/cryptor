package packet

import (
	"net"
)

const (
	// TypeTransport ...
	TypeTransport byte = iota

	// TypeHandshakeInitializer ...
	TypeHandshakeInitializer

	// TypeHandshakeResponder ...
	TypeHandshakeResponder
)

// Packet contains the necessary info for sending and receiving messages over
// the network and handling them internally.
type Packet struct {
	Type    byte
	Payload []byte
	Address *net.UDPAddr
}

package packet

import (
	"net"
)

const (
	// TypeTransport ...
	TypeTransport byte = 0

	// TypeHandshakeInitializer ...
	TypeHandshakeInitializer byte = 1

	// TypeHandshakeResponder ...
	TypeHandshakeResponder byte = 2
)

// Packet contains the necessary info for sending and receiving messages over
// the network and handling them internally.
type Packet struct {
	Type    byte
	Payload []byte
	Address *net.UDPAddr
}

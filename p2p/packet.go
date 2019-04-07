package p2p

import (
	"net"
)

// TODO Redefine packet structure and handling

const (
	packetTypeTransport  byte = 0
	packetTypeHandshakeI byte = 1
	packetTypeHandshakeR byte = 2
)

// Packet contains the necessary info for sending and receiveing messages over
// the network and handling them internally.
type Packet struct {
	Type    byte
	Payload []byte
	Address *net.UDPAddr
}

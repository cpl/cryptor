package p2p

import "net"

// UDPPacketSize is the static expected size of a single UDP packet.
const UDPPacketSize = 1024

// UDPPacket holds the origin address of the packet and the data of the packet.
type UDPPacket struct {
	data []byte
	addr *net.UDPAddr
}

// NewUDPPacket is a test function for generating UDP packets.
func NewUDPPacket(data []byte, addr *net.UDPAddr) UDPPacket {
	return UDPPacket{
		data: data,
		addr: addr,
	}
}

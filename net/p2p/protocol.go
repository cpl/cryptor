package p2p

type protocolFunc func(...UDPPacket) UDPPacket

// Protocol ...
type Protocol struct {
	Name    string
	Version uint
	Run     protocolFunc
}

var (
	keyExchange = Protocol{
		Name:    "Key Exchange",
		Version: 0,
	}
	peerExchange = Protocol{
		Name:    "Peer Exchange",
		Version: 0,
	}
	packetFwd = Protocol{
		Name:    "Package Forwarding",
		Version: 0,
	}
)

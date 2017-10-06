package p2p

import (
	"net"

	"github.com/thee-engineer/cryptor/network"
)

// Peer ...
type Peer struct {
	pubk string // TODO: Change this to ECDSA key? (maybe)
	addr *net.UDPAddr
}

// NewPeer ...
func NewPeer(ip string, port int) *Peer {
	return &Peer{
		addr: network.IPPToUDP(ip, port),
	}
}

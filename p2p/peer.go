package p2p

import (
	"net"
)

// Peer represents a foreign node connected to the local node.
type Peer struct {
	udpAddr *net.UDPAddr
	tcpAddr *net.TCPAddr
}

// NewPeer creates a peer object given an IP:PORT pair. This is for testing and
// debugging purposes. Not to be used in production.
func NewPeer(ip string, tcp, udp int) *Peer {
	return &Peer{
		udpAddr: &net.UDPAddr{Port: udp},
		tcpAddr: &net.TCPAddr{Port: tcp, IP: net.ParseIP(ip)},
	}
}

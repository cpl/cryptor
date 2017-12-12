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

// AddPeer adds a given peer to the peer memory map.
func (n *Node) AddPeer(peer *Peer) {
	select {
	case n.addp <- peer:
	}
}

// RemovePeer sends a signal to the node to remove a peer
func (n *Node) RemovePeer(peer *Peer) {
	select {
	case n.remp <- peer:
	}
}

// Peers returns a list of all peers related to this Node.
func (n *Node) Peers() []*Peer {
	var peerList []*Peer

	select {
	case n.pops <- func(peers map[string]*Peer) {
		for _, p := range peers {
			peerList = append(peerList, p)
		}
	}:
		<-n.popd
	}

	return peerList
}

// PeerCount returns the number of related peers for this Node.
func (n *Node) PeerCount() int {
	var count int

	select {
	case n.pops <- func(peerList map[string]*Peer) { count = len(peerList) }:
		<-n.popd
	case <-n.quit:
	}

	return count
}

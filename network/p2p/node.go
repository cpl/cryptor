package p2p

import (
	"net"

	"github.com/thee-engineer/cryptor/network"
)

// Node ..
type Node struct {
	NodeConfig // Static configuration generated at node creation

	addr *net.UDPAddr  // Node UDP address
	quit chan struct{} // Stops node from running when it receives
	errc chan *error   // Channell for transmiting errors
	addp chan *Peer    // Add peer request channel
	pops chan peerFunc // Peer count and peer list operations
	popd chan struct{} // Peer operation done

	peers map[string]*Peer // Memory map with key/value peer pairs
}

// Function for peer list and peer count
type peerFunc func(map[string]*Peer)

// NewNode ...
func NewNode(ip string, port int, quit chan struct{}) *Node {
	return &Node{
		addr:  network.IPPToUDP(ip, port),
		quit:  quit,
		addp:  make(chan *Peer),
		pops:  make(chan peerFunc),
		popd:  make(chan struct{}),
		peers: make(map[string]*Peer),
	}
}

// Start ...
func (n *Node) Start() {
	for {
		select {
		case <-n.quit:
			return
		case peer := <-n.addp:
			n.peers[peer.addr.String()] = peer
		case operation := <-n.pops:
			operation(n.peers)
			n.popd <- struct{}{}
		}
	}
}

// Stop ...
func (n *Node) Stop() {
	close(n.quit)
}

// AddPeer ...
func (n *Node) AddPeer(ip string, port int) {
	select {
	case <-n.quit:
	case n.addp <- NewPeer(ip, port):
	}
}

// Peers ...
func (n *Node) Peers() []*Peer {
	var peerList []*Peer

	select {
	case n.pops <- func(peers map[string]*Peer) {
		for _, p := range peers {
			peerList = append(peerList, p)
		}
	}:
		<-n.popd
	case <-n.quit:

	}

	return peerList
}

// PeerCount ...
func (n *Node) PeerCount() int {
	var count int

	select {
	case n.pops <- func(peerList map[string]*Peer) { count = len(peerList) }:
		<-n.popd
	case <-n.quit:
	}

	return count
}

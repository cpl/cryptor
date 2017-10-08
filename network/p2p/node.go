package p2p

import (
	"fmt"
	"net"

	"github.com/thee-engineer/cryptor/network"
)

// Node ..
type Node struct {
	NodeConfig // Static configuration generated at node creation

	addr *net.UDPAddr  // Node UDP address
	quit chan struct{} // Stops node from running when it receives
	errc chan error    // Channell for transmiting errors

	addp chan *Peer    // Add peer request channel
	pops chan peerFunc // Peer count and peer list operations
	popd chan struct{} // Peer operation done

	peers map[string]*Peer  // Memory map with key/value peer pairs
	token map[string][]byte // List of tokens used in requests

	inbound  chan Packet
	outbound chan Packet
}

// Function for peer list and peer count
type peerFunc func(map[string]*Peer)

// NewNode ...
func NewNode(ip string, port int, quit chan struct{}) *Node {
	return &Node{
		addr:     network.IPPToUDP(ip, port),
		quit:     quit,
		addp:     make(chan *Peer),
		errc:     make(chan error),
		pops:     make(chan peerFunc),
		popd:     make(chan struct{}),
		peers:    make(map[string]*Peer),
		inbound:  make(chan Packet, 1024),
		outbound: make(chan Packet, 1024),
	}
}

// Start ...
func (n *Node) Start() {

	go n.connect(n.addr.Port)

	for {
		select {
		case err := <-n.errc:
			fmt.Println("err:", err) // DEBUG
		case <-n.quit:
			return
		case peer := <-n.addp:
			n.peers[peer.addr.String()] = peer
		case operation := <-n.pops:
			operation(n.peers)
			n.popd <- struct{}{}
		case packet := <-n.inbound:
			fmt.Println(
				"| packet from |", packet.addr.String(),
				"| containing  |", string(packet.data))
			n.outbound <- Packet{data: []byte("ok\n"), addr: packet.addr}
		}
	}
}

// Stop ...
func (n *Node) Stop() {
	close(n.quit)
}

// AddPeer ...
func (n *Node) AddPeer(peer *Peer) {
	select {
	case <-n.quit:
	case n.addp <- peer:
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

// Packet ...
type Packet struct {
	data []byte
	addr *net.UDPAddr
}

func (n *Node) connect(port int) {

	conn, err := net.ListenUDP("udp", &net.UDPAddr{Port: port})
	if err != nil {
		n.errc <- err
		return
	}

	for {
		select {
		case packet := <-n.outbound:
			_, err := conn.WriteToUDP(packet.data, packet.addr)
			if err != nil {
				n.errc <- err
				continue
			}
		default:
			buffer := make([]byte, 1024)
			r, addr, err := conn.ReadFromUDP(buffer)
			if err != nil {
				n.errc <- err
				continue
			}
			n.inbound <- Packet{buffer[:r], addr}
		}
	}
}

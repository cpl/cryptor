package p2p

import (
	"errors"
	"log"
	"net"
	"sync"
)

// Node is a representation of the machine on the Cryptor network. A node
// listens both UDP (Peer discovery and requests) and TCP (Chunk sharing and
// handshakes).
type Node struct {
	config NodeConfig // Static configuration generated at node creation

	lock    sync.Mutex // Protects the running state
	running bool       // Node state

	udpAddr *net.UDPAddr // Node UDP address
	tcpAddr *net.TCPAddr // Node TCP address

	quit chan struct{} // Stops node from running when it receives
	errc chan error    // Channel for transmiting errors

	addp chan *Peer    // Add peer request channel
	pops chan peerFunc // Peer count and peer list operations
	popd chan struct{} // Peer operation done

	peers map[string]*Peer  // Memory map with key/value peer pairs
	token map[string][]byte // List of tokens used in requests

	udpIncoming chan UDPPacket // Incoming UDP channel
	udpOutgoing chan UDPPacket // Outgoing UDP channel
}

// Function for peer list and peer count
type peerFunc func(map[string]*Peer)

// NewNode returns a Node attached to the given IP:PORT pair, and controlled
// using the given quit channel. This is more for testing and debugging than
// actual production.
func NewNode(ip string, tcp, udp int, quit chan struct{}) *Node {
	// Create node configuration
	config := NodeConfig{
		IP:  net.ParseIP(ip),
		TCP: tcp,
		UDP: udp,
	}

	return &Node{
		config:      config,
		udpAddr:     &net.UDPAddr{Port: config.UDP},
		tcpAddr:     &net.TCPAddr{IP: config.IP, Port: config.TCP},
		quit:        quit,
		addp:        make(chan *Peer),
		errc:        make(chan error),
		pops:        make(chan peerFunc),
		popd:        make(chan struct{}),
		peers:       make(map[string]*Peer),
		udpIncoming: make(chan UDPPacket, UDPPacketSize),
		udpOutgoing: make(chan UDPPacket, UDPPacketSize),
	}
}

// Start begins the listening process for the Node on the network and all
// operations/requests handling.
func (n *Node) Start() {
	if n.running {
		n.errc <- errors.New("node: already running")
		return
	}
	n.lock.Lock()
	n.running = true
	n.lock.Unlock()

	go n.listen()

	for {
		select {
		case err := <-n.errc:
			log.Println("err:", err) // DEBUG
		case <-n.quit:
			n.running = false
			return
		case peer := <-n.addp:
			n.peers[peer.tcpAddr.String()] = peer
		case operation := <-n.pops:
			operation(n.peers)
			n.popd <- struct{}{}
		case packet := <-n.udpIncoming:
			// DEBUG
			log.Println(
				"| packet from |", packet.addr.String(),
				"| containing  |", string(packet.data))
			// DEBUG
			// n.udpOutgoing <- UDPPacket{
			// 	data: []byte("ok\n"),
			// 	addr: packet.addr}
		}
	}
}

func (n *Node) listen() {
	// Listen for UDP on the given node address port
	conn, err := net.ListenUDP("udp4", n.udpAddr)
	if err != nil {
		n.errc <- err
		return
	}

	go n.receive(conn) // Pass incoming packets to incoming channel
	go n.send(conn)    // Pass outgoing packets to outgoing channel
}

func (n *Node) receive(conn *net.UDPConn) {
	buffer := make([]byte, UDPPacketSize)
	for {
		r, addr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			n.errc <- err
			continue
		}
		n.udpIncoming <- UDPPacket{buffer[:r], addr}
	}
}

func (n *Node) send(conn *net.UDPConn) {
	for packet := range n.udpOutgoing {
		_, err := conn.WriteToUDP(packet.data, packet.addr)
		if err != nil {
			n.errc <- err
			continue
		}
	}
}

// Send is currently a test method that sends one packet.
func (n *Node) Send(packet UDPPacket) {
	n.udpOutgoing <- packet
}

// UDPAddr is just for testing for now.
func (n *Node) UDPAddr() *net.UDPAddr {
	return n.udpAddr
}

// Stop closes the quit channel of the Node.
func (n *Node) Stop() {
	n.lock.Lock()
	defer n.lock.Unlock()
	if !n.running {
		return
	}
	n.running = false

	close(n.quit)
}

// AddPeer adds a given peer to the peer memory map.
func (n *Node) AddPeer(peer *Peer) {
	select {
	case n.addp <- peer:
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

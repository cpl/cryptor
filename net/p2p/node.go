package p2p

import (
	"errors"
	"log"
	"net"
	"sync"
)

// NodeConfig is a static structure used for creating node configuration
// presets used before node startup.
type NodeConfig struct {
	TrustedPeers []*Peer

	quitChan chan interface{}
	discChan chan interface{}
}

// Node can run in offline or online mode. In offline mode only local settings
// can be changed and peers managed. Once online the node will become a peer
// in the cryptor network. The node will send and receive requests. Peers
// will be discovered after connecting to the network unless running in trust
// mode where only trusted peers can connect. The local node will act as both
// a reciver and sender.
type Node struct {
	address string     // Host address (IPv4/IPv6/hostname)
	port    string     // Port for listening
	config  NodeConfig // Creation configuration

	lock sync.Mutex // Mutex lock for critical code section

	isRunning   bool
	isConnected bool

	peers map[string]*Peer // List of known peers

	incoming chan string // Incoming packets buffer
	outgoing chan string // Outgoing packets buffer

	errChan chan error       // Channel for logging error
	logChan chan interface{} // Channel for loggin messages

	// Run functions on peers using peerOp chan
	peerOp     chan peerFunc
	peerOpDone chan interface{}

	quit       chan interface{} // Stops the node from running
	disconnect chan interface{} // Disconnects the node from the network
}

// NewNode constructs a node with default configurations
func NewNode(addr, port string, config *NodeConfig) Node {

	var qc, dc chan interface{}

	if config != nil {
		qc = config.quitChan
		dc = config.discChan
	} else {
		qc = make(chan interface{})
		dc = make(chan interface{})
	}

	return Node{
		address: addr,
		port:    port,

		peers: make(map[string]*Peer),

		logChan: make(chan interface{}),
		errChan: make(chan error),

		incoming: make(chan string),
		outgoing: make(chan string),

		peerOp:     make(chan peerFunc),
		peerOpDone: make(chan interface{}),

		quit:       qc,
		disconnect: dc,
	}
}

// Start allows the node to recive commands from a control interface.
func (n *Node) Start() {
	// Check if the node is already running
	n.lock.Lock()
	if n.isRunning {
		n.errChan <- errors.New("already started")
		n.lock.Unlock()
		return
	}
	n.isRunning = true
	n.lock.Unlock()

	// Start listening on local channels
	log.Println("node log: started")
	for {
		select {
		case err := <-n.errChan:
			log.Println("node err:", err)
		case msg := <-n.logChan:
			log.Println("node log:", msg)
		case operation := <-n.peerOp:
			operation(n.peers)
			n.peerOpDone <- nil
		case <-n.quit:
			log.Println("node log: stopping")
			n.isRunning = false
			return
		}
	}
}

// Stop the node from running (also closes any active connections).
func (n *Node) Stop() {
	// Check if the node is not running
	if !checkRunning(n) {
		return
	}
	n.quit <- nil
}

// Connect joins the cryptor network and begins communication with its peers
func (n *Node) Connect() {
	if !checkRunning(n) {
		return
	}

	// Check if the node is already connected
	n.lock.Lock()
	if n.isConnected {
		n.errChan <- errors.New("already connected")
		n.lock.Unlock()
		return
	}
	n.isConnected = true
	n.lock.Unlock()

	// Bind to given address and port
	listener, err := net.Listen("tcp", n.address+":"+n.port)
	if err != nil {
		n.errChan <- err
	}
	n.logChan <- "listening on " + listener.Addr().String()

	// Start listening and handling connections
	for {
		conn, err := listener.Accept()
		if err != nil {
			n.errChan <- err
		}
		n.logChan <- "connection from " + conn.RemoteAddr().String()
	}
}

// Disconnect stops the node from sending or receiving on the network. Keeps
// the node running for control operations.
func (n *Node) Disconnect() {
	if !checkRunning(n) {
		return
	}

	// Check if the node is not connected
	n.lock.Lock()
	defer n.lock.Unlock()
	if !n.isConnected {
		n.errChan <- errors.New("not connected")
		return
	}

	n.isConnected = false
}

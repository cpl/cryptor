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
	trustedPeers []*Peer
}

// Node can run in offline or online mode. In offline mode only local settings
// can be changed and peers managed. Once online the node will become a peer
// in the cryptor network. The node will send and receive requests. Peers
// will be discovered after connecting to the network unless running in trust
// mode where only trusted peers can connect. The local node will act as both
// a reciver and sender.
type Node struct {
	config NodeConfig // Configuration

	address string // Host address (IPv4/IPv6/hostname)
	port    string // Port for listening

	udpConn net.PacketConn // UDP Listener

	lock sync.Mutex // Mutex lock for critical code section

	isRunning   bool
	isConnected bool

	peers peerMap // List of known peers

	incoming chan []byte // Incoming packets buffer
	outgoing chan []byte // Outgoing packets buffer

	errChan chan error       // Channel for logging error
	logChan chan interface{} // Channel for loggin messages

	// Run functions on peers using peerOp chan
	peerOp     chan peerFunc
	peerOpDone chan interface{}

	quit       chan interface{} // Stops the node from running
	disconnect chan interface{} // Disconnects the node from the network
}

// NewNode constructs a node with default configurations
func NewNode(addr, port string, config *NodeConfig) *Node {
	return &Node{
		address: addr,
		port:    port,

		peers: make(peerMap),

		logChan: make(chan interface{}),
		errChan: make(chan error),

		incoming: make(chan []byte, 10),
		outgoing: make(chan []byte, 10),

		peerOp:     make(chan peerFunc),
		peerOpDone: make(chan interface{}),

		quit:       make(chan interface{}),
		disconnect: make(chan interface{}),
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

	// Start goroutine
	go n.run()
	n.logChan <- "started"
}

// Stop the node from running (also closes any active connections).
func (n *Node) Stop() {
	// Check if the node is not running
	if !checkRunning(n) {
		return
	}

	// Disconnect if connected
	if n.isConnected {
		n.Disconnect()
	}

	n.quit <- nil
	n.isRunning = false

	log.Println("node log: stopped")
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

	// Bind to given address and port
	var err error
	n.udpConn, err = net.ListenPacket("udp", n.address+":"+n.port)
	if err != nil {
		n.errChan <- err
		n.lock.Unlock()
		return
	}

	// Connection is done
	n.isConnected = true
	n.lock.Unlock()

	// Read incoming
	go n.listen()

	// Confirm connection
	n.logChan <- "listening on " + n.udpConn.LocalAddr().String()
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

	// Stop connection
	err := n.udpConn.Close()
	if err != nil {
		n.errChan <- err
	}

	// Send disconnect signal
	n.disconnect <- nil
	n.isConnected = false

	n.logChan <- "disconnected"
}

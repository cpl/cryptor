/*
Package node defines the structure of a node (UDP listener and receiver). A
node is capable of running in offline mode where it may take commands and
configurations.
*/
package node // import "cpl.li/go/cryptor/p2p/node"

import (
	"log"
	"net"
	"os"
	"sync"

	"cpl.li/go/cryptor/crypt/ppk"
	"cpl.li/go/cryptor/p2p/peer"
	"cpl.li/go/cryptor/p2p/proto"
)

// Node represents the local machine running and/or connected to the Cryptor
// network. Other nodes are represented as peers.
type Node struct {

	// a custom logger
	logger *log.Logger

	// network aspect of a node
	net struct {
		sync.RWMutex

		addr *net.UDPAddr // address for listening and sending
		conn *net.UDPConn // listening connection

		recv chan *proto.Packet // receive network traffic
		send chan *proto.Packet // send network responses
	}

	// state of the node
	state struct {
		sync.RWMutex

		isRunning   bool // if running, the node is open to taking commands
		isConnected bool // if online, the node is available on the network
	}

	// static identity of a node
	identity struct {
		sync.RWMutex

		privateKey ppk.PrivateKey // static private key
		publicKey  ppk.PublicKey  // static public key (public identifier)
	}

	// maps which allow lookup of peers based on diffrent keys
	lookup struct {
		sync.RWMutex

		// a key,value map of publicKey/peer
		peers map[ppk.PublicKey]*peer.Peer

		// a key,value map of address/peer
		addres map[string]*peer.Peer
	}

	// communication covers the channels used by the node to send information
	// across concurrent routines
	comm struct {
		err chan error       // passing errors from other routines
		exi chan interface{} // passing the exit/stop signal
		dis chan interface{} // passing the disconnect signal
	}

	// TODO firewall
	// firewall determines which hosts and ip ranges are allowed to connect
	// if whitelist is used, only hosts defined in the list will be allowed
	// if blacklist is used, all hosts that are NOT in the list are allowed
	// only use one list, the other should be empty, whitelist has priority
	firewall struct {
		whitelist []string // list of allowed hosts
		blacklist []string // list of banned  hosts
	}
}

// New creates a node running on the local machine. The default starting
// state is NOT RUNNING and OFFLINE. Allowing the Node to be further configured
// before starting and connecting to the Cryptor Network.
func New(name string) *Node {
	n := new(Node)

	// initialize logger
	n.logger = log.New(os.Stdout, name+": ", log.Ldate|log.Ltime)

	// initialize communication channels
	n.comm.err = make(chan error)
	n.comm.exi = make(chan interface{})
	n.comm.dis = make(chan interface{})

	// initialize network forwarding channels
	n.net.recv = make(chan *proto.Packet)
	n.net.send = make(chan *proto.Packet)

	// default state
	n.state.isRunning = false
	n.state.isConnected = false

	// generate node keys
	// TODO Implement PBKDF2 for this
	n.identity.privateKey, _ = ppk.NewPrivateKey()
	n.identity.publicKey = n.identity.privateKey.PublicKey()

	// initialize lookup maps
	n.lookup.peers = make(map[ppk.PublicKey]*peer.Peer)
	n.lookup.addres = make(map[string]*peer.Peer)

	n.logger.Println("created")

	return n
}

// Start enable the node to receive and parse commands locally.
func (n *Node) Start() {
	// lock state
	n.state.Lock()
	defer n.state.Unlock()

	// ignore if node is already running
	if n.state.isRunning {
		n.logger.Println("node is already running")
		return
	}

	// set node as running
	n.state.isRunning = true
	n.logger.Println("started")

	// start service
	go n.run()
}

// Stop attempts to safely shutdown the Cryptor Node and any active tasks.
func (n *Node) Stop() {
	// lock state
	n.state.Lock()
	defer n.state.Unlock()

	// ignore if node is not running
	if !n.state.isRunning {
		return
	}

	// disconnect node if connected before continuing
	if n.state.isConnected {
		n.disconnect()
	}

	// set node as stopped
	n.state.isRunning = false
	n.logger.Println("stopped")

	// send stop signal
	n.comm.exi <- nil
}

// Connect binds the current node to the machine network and Cryptor network.
func (n *Node) Connect() {
	// lock state
	n.state.Lock()
	defer n.state.Unlock()

	// ignore if node is already connected
	if n.state.isConnected {
		return
	}

	// ignore if node is not running
	if !n.state.isRunning {
		n.logger.Println("can't connect to network, node not running")
		return
	}

	// create network bind
	if err := n.connect(); err != nil {
		n.comm.err <- err
		return
	}

	// set node as connected
	n.state.isConnected = true
	n.logger.Printf("connected on %s\n", n.net.conn.LocalAddr().String())

	// start forwarding service
	go n.forward()

	// start network listening
	go n.listen()
}

// Disconnect stops the network communication capabilities of the node. This
// does not stop the node from running and receiving local commands and config
// updates.
func (n *Node) Disconnect() {
	// lock state
	n.state.Lock()
	defer n.state.Unlock()

	n.disconnect()
}

// PublicKey returns the node static public key.
func (n *Node) PublicKey() ppk.PublicKey {
	return n.identity.publicKey
}

// Addr returns the node listening address.
func (n *Node) Addr() string {
	// if not connected, return addr
	if !n.state.isConnected {
		return n.net.addr.String()
	}

	return n.net.conn.LocalAddr().String()
}

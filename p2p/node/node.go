/*
Package node offers the representation of the local machine identity on the
Cryptor Network. A node is in charge of its peers, keys, connections, etc.
*/
package node // import "cpl.li/go/cryptor/p2p/node"

import (
	"errors"
	"log"
	"net"
	"os"
	"sync"

	"cpl.li/go/cryptor/crypt/ppk"
	"cpl.li/go/cryptor/p2p"
	"cpl.li/go/cryptor/p2p/packet"
	"cpl.li/go/cryptor/p2p/peer"
)

// Node represents the local machine running and/or connected to the Cryptor
// network. Other nodes are represented as peers.
type Node struct {
	// a custom logger
	logger *log.Logger

	// network aspect of a node
	net struct {
		sync.RWMutex

		addr *net.UDPAddr // address for listening and receiving
		conn *net.UDPConn // bind/connection for listening and receiving

		recv chan *packet.Packet // receive network traffic
		send chan *packet.Packet // send network responses
	}

	// state of the node
	state struct {
		sync.RWMutex

		isRunning   bool // if running, the node is open to taking commands
		isConnected bool // if online, the node is available on the network
	}

	// static identity of a node
	identity struct {
		privateKey ppk.PrivateKey // static private key
		publicKey  ppk.PublicKey  // static public key (public identifier)
	}

	// maps which allow lookup of peers based on different keys
	lookup struct {
		sync.RWMutex

		// total number of peers
		count int

		// a key,value map of publicKey/peer
		peers map[ppk.PublicKey]*peer.Peer

		// a key,value map of address/peer
		table map[uint64]*peer.Peer
	}

	// communication covers the channels used by the node to send information
	// across concurrent routines
	comm struct {
		err chan error       // passing errors from other routines
		exi chan interface{} // passing the exit/stop signal
		dis chan interface{} // passing the disconnect signal
	}

	// metadata
	meta struct {
		errCount int
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

// NewNode creates a node running on the local machine. The default starting
// state is NOT RUNNING and OFFLINE. Allowing the Node to be further configured
// before starting and connecting to the Cryptor Network.
func NewNode(name string, key ppk.PrivateKey) *Node {
	n := new(Node)

	// initialize logger
	// TODO Add logger configuration
	n.logger = log.New(os.Stdout, name+": ", log.Ldate|log.Ltime)

	// initialize communication channels
	n.comm.err = make(chan error)
	n.comm.exi = make(chan interface{})
	n.comm.dis = make(chan interface{})

	// initialize network forwarding channels
	n.net.recv = make(chan *packet.Packet)
	n.net.send = make(chan *packet.Packet)

	// default state
	n.state.isRunning = false
	n.state.isConnected = false

	// assign node keys
	n.identity.privateKey = key
	n.identity.publicKey = key.PublicKey()

	// initialize lookup maps
	n.lookup.peers = make(map[ppk.PublicKey]*peer.Peer)
	n.lookup.table = make(map[uint64]*peer.Peer)

	n.logger.Println("created")

	return n
}

// Start enable the node to receive and parse commands locally.
func (n *Node) Start() error {
	// lock state
	n.state.Lock()
	defer n.state.Unlock()

	// ignore if node is already running
	if n.state.isRunning {
		n.meta.errCount++
		return errors.New("can't start, node already running")
	}

	// set node as running
	n.state.isRunning = true

	// start service
	go n.run()
	n.logger.Println("started")

	return nil
}

// Stop attempts to safely shutdown the Cryptor Node and any active tasks. If
// the Node is connected to the network, a Disconnect will run before stopping.
func (n *Node) Stop() error {
	// lock state
	n.state.Lock()
	defer n.state.Unlock()

	// ignore if node is not running
	if !n.state.isRunning {
		n.meta.errCount++
		return errors.New("can't stop, node not running")
	}

	// disconnect node if connected before continuing
	if n.state.isConnected {
		if err := n.disconnect(); err != nil {
			// on error, cancel stopping routine
			n.meta.errCount++
			return err
		}
	}

	// set node as stopped
	n.comm.exi <- nil
	n.state.isRunning = false

	n.logger.Println("stopped")

	return nil
}

// Connect binds the current node to the machine network and Cryptor network.
func (n *Node) Connect() error {
	// lock state
	n.state.Lock()
	defer n.state.Unlock()

	// ignore if node is not running
	if !n.state.isRunning {
		n.meta.errCount++
		return errors.New("can't connect, node not running")
	}

	// ignore if node is already connected
	if n.state.isConnected {
		n.meta.errCount++
		return errors.New("can't connect to network, already connected")
	}

	// create network bind
	if err := n.connect(); err != nil {
		return err
	}

	// set node as connected
	n.state.isConnected = true

	// start network listening
	go n.listen()
	n.logger.Printf("connected on %s\n", n.net.conn.LocalAddr().String())

	return nil
}

// Disconnect stops the network communication capabilities of the node. This
// does not stop the node from running and receiving local commands and config
// updates.
func (n *Node) Disconnect() error {
	// lock state
	n.state.Lock()
	defer n.state.Unlock()

	// ignore if node is not running
	if !n.state.isRunning {
		n.meta.errCount++
		return errors.New("can't disconnect, node not running")
	}

	// ignore if node is not connected
	if !n.state.isConnected {
		n.meta.errCount++
		return errors.New("can't disconnect, node not connected")
	}

	// attempt disconnect
	if err := n.disconnect(); err != nil {
		n.meta.errCount++
		return err
	}

	return nil
}

// PublicKey returns the node static public key.
func (n *Node) PublicKey() ppk.PublicKey {
	return n.identity.publicKey
}

// Addr returns the node listening address. If the node is running we return the
// live connection address, otherwise the pre-configured address.
func (n *Node) Addr() string {
	// lock
	n.state.RLock()
	defer n.state.RUnlock()

	// if not connected, return addr
	if !n.state.isConnected {
		return n.net.addr.String()
	}

	return n.net.conn.LocalAddr().String()
}

// SetAddr sets the nodes listening address and port. The expected address
// must be "host:port".
func (n *Node) SetAddr(addr string) (err error) {
	// lock node state
	n.state.RLock()
	defer n.state.RUnlock()

	// ignore if node is connected
	if n.state.isConnected {
		n.meta.errCount++
		return errors.New("can't change address while node is connected")
	}

	// lock network
	n.net.Lock()
	defer n.net.Unlock()

	// resolve address (no change yet, validating and checking errors first)
	newaddr, err := net.ResolveUDPAddr(p2p.Network, addr)
	if err != nil {
		n.meta.errCount++
		return err
	}

	// reject change if address matches
	if newaddr.String() == n.net.addr.String() {
		n.meta.errCount++
		return errors.New("can't change address, same as current address")
	}

	// announce address change
	n.logger.Printf("change address from %s to %s\n",
		n.net.addr.String(), newaddr.String())

	// set new address
	n.net.addr = newaddr

	return nil
}

// ErrCount returns the number of errors which ocurred during runtime.
func (n *Node) ErrCount() int {
	return n.meta.errCount
}

// StateNode defines an alias of the underlaying type used for enumerating
// the possible node states.
type StateNode byte

const (
	// StateStopped is the default starting state of a Node, or the state
	// reached after calling node.Stop(). During this state node operations
	// are limited. A node also can't be connected during this state.
	StateStopped StateNode = 0

	// StateRunning is the state after successfully calling node.Start() or
	// node.Disconnect(). During this state the node allows for some
	// configurations or actions to be taken.
	StateRunning StateNode = 1

	// StateConnected is entered after calling node.Connect() on a running node.
	// Node must be running in order to also be connected to the network.
	StateConnected StateNode = 2
)

// a simple lookup for state "ID" to state string, allows for nice displaying
var stateNames = map[StateNode]string{
	StateStopped:   "STOPPED",
	StateRunning:   "RUNNING",
	StateConnected: "CONNECTED",
}

// String will return the name of the current state as a string from stateNames.
func (s StateNode) String() string {
	return stateNames[s]
}

// other possible future states: Crashed, Restarting, Blocked, Idle

// State returns a node state/status. This state can be one of the
// following constant pre-defined byte values:
// - Stopped
// - Running
// - Connected
// Other values may be added in the future.
func (n *Node) State() StateNode {
	// lock states for reading
	n.state.RLock()
	defer n.state.RUnlock()

	// check state and return
	switch {
	case n.state.isConnected:
		return StateConnected
	case n.state.isRunning:
		return StateRunning
	default:
		return StateStopped
	}
}

package node

import (
	"errors"
	"log"
	"net"
	"os"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"

	"cpl.li/go/cryptor/crypt/mwords"
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

	// node name
	name string

	// network aspect of a node
	net struct {
		sync.RWMutex

		addr *net.UDPAddr // address for listening and receiving
		conn *net.UDPConn // bind/connection for listening and receiving
	}

	// state of the node
	state struct {
		sync.RWMutex

		starting sync.WaitGroup
		stopping sync.WaitGroup

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

		recv chan *packet.Packet // receive network traffic
		send chan *packet.Packet // send network responses
	}

	// metadata
	meta struct {
		errCount uint32
		cpuCount int
	}
}

// NewNode creates a node running on the local machine. The default starting
// state is NOT RUNNING and OFFLINE. Allowing the Node to be further configured
// before starting and connecting to the Cryptor Network.
func NewNode(name string, key ppk.PrivateKey) *Node {
	n := new(Node)

	// initialize logger
	// TODO Add logger configuration
	if name == "" {
		name = strings.Join(mwords.RandomWords(3), "-")
	}
	n.name = name
	n.logger = log.New(os.Stdout, name+": ", log.Ldate|log.Ltime)

	// ! DEBUG, buffer sizes, will add configuration later

	// initialize communication channels
	n.comm.err = make(chan error, 20)
	n.comm.exi = make(chan interface{})
	n.comm.dis = make(chan interface{})

	// initialize network forwarding channels
	n.comm.recv = make(chan *packet.Packet, 20)
	n.comm.send = make(chan *packet.Packet, 20)

	// default state
	n.state.isRunning = false
	n.state.isConnected = false

	// assign node keys
	n.identity.privateKey = key
	n.identity.publicKey = key.PublicKey()

	// initialize lookup maps
	n.lookup.peers = make(map[ppk.PublicKey]*peer.Peer)
	n.lookup.table = make(map[uint64]*peer.Peer)

	// meta
	n.meta.cpuCount = runtime.NumCPU()*2 + 1

	n.logger.Printf("created node [%s - %d]\n", name, n.meta.cpuCount)

	return n
}

// Start enable the node to receive and parse commands locally.
func (n *Node) Start() error {
	// lock state
	n.state.Lock()
	defer n.state.Unlock()

	// ignore if node is already running
	if n.state.isRunning {
		atomic.AddUint32(&n.meta.errCount, 1)
		return errors.New("can't start, node already running")
	}

	// start service routines
	n.state.starting.Add(n.meta.cpuCount)
	go n.run()
	for i := 0; i < runtime.NumCPU(); i++ {
		go n.recv()
		go n.send()
	}

	// wait for all routines to start
	n.state.starting.Wait()

	// set node as running
	n.state.isRunning = true
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
		atomic.AddUint32(&n.meta.errCount, 1)
		return errors.New("can't stop, node not running")
	}

	// disconnect node if connected before continuing
	if n.state.isConnected {
		if err := n.disconnect(); err != nil {
			// on error, cancel stopping routine
			atomic.AddUint32(&n.meta.errCount, 1)
			return err
		}
	}

	// stop service routines
	n.state.stopping.Add(n.meta.cpuCount)
	for i := 0; i < n.meta.cpuCount; i++ {
		n.comm.exi <- nil
	}

	// wait for all routines to stop
	n.state.stopping.Wait()

	// set node as stopped
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
		atomic.AddUint32(&n.meta.errCount, 1)
		return errors.New("can't connect to network, node not running")
	}

	// ignore if node is already connected
	if n.state.isConnected {
		atomic.AddUint32(&n.meta.errCount, 1)
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
		atomic.AddUint32(&n.meta.errCount, 1)
		return errors.New("can't disconnect, node not running")
	}

	// ignore if node is not connected
	if !n.state.isConnected {
		atomic.AddUint32(&n.meta.errCount, 1)
		return errors.New("can't disconnect, node not connected")
	}

	// attempt disconnect
	if err := n.disconnect(); err != nil {
		atomic.AddUint32(&n.meta.errCount, 1)
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
		atomic.AddUint32(&n.meta.errCount, 1)
		return errors.New("can't change address while node is connected")
	}

	// lock network
	n.net.Lock()
	defer n.net.Unlock()

	// resolve address (no change yet, validating and checking errors first)
	newaddr, err := net.ResolveUDPAddr(p2p.Network, addr)
	if err != nil {
		atomic.AddUint32(&n.meta.errCount, 1)
		return err
	}

	// reject change if address matches
	if newaddr.String() == n.net.addr.String() {
		atomic.AddUint32(&n.meta.errCount, 1)
		return errors.New("can't change address, same as current address")
	}

	// announce address change
	n.logger.Printf("change address from %s to %s\n",
		n.net.addr.String(), newaddr.String())

	// set new address
	n.net.addr = newaddr

	return nil
}

// ErrCount returns the number of errors which occurred during runtime.
func (n *Node) ErrCount() uint32 {
	return atomic.LoadUint32(&n.meta.errCount)
}

// Wait will block until any state operations are done.
func (n *Node) Wait() {
	n.state.RLock()
	n.state.starting.Wait()
	n.state.stopping.Wait()
	n.state.RUnlock()
}

// Name will return the node name assigned at creation (either given or random).
func (n *Node) Name() string {
	return n.name
}

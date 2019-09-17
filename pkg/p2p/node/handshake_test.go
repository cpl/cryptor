package node_test

import (
	"net"
	"testing"

	"cpl.li/go/cryptor/pkg/crypt/ppk"
	"cpl.li/go/cryptor/pkg/p2p"
	"cpl.li/go/cryptor/pkg/p2p/node"
	"cpl.li/go/cryptor/pkg/p2p/noise"
	"cpl.li/go/cryptor/pkg/p2p/peer"
	"github.com/stretchr/testify/assert"
)

func TestNodeHandshake(t *testing.T) {
	// create node keys
	peerPrivate, _ := ppk.NewPrivateKey()
	peerPublic := peerPrivate.PublicKey()

	// create nodes
	n := node.NewNode("test", zeroKey)

	// start node
	assert.Nil(t, n.Start())

	// set address
	assert.Nil(t, n.SetAddr("localhost:"))

	// add external peer
	p, err := n.PeerAdd(peerPublic, "", 1)
	assert.Nil(t, err)

	// connect
	assert.Nil(t, n.Connect())

	// check peer count
	assertPeerCount(t, n, 1)

	// set peer addr
	assert.Nil(t, p.SetAddr("localhost:12345"))

	// simulate peer connection
	pConn, err := net.ListenUDP(p2p.Network, p.AddrUDP())
	assert.Nil(t, err)

	// attempt handshake with peer
	assert.Nil(t, n.Handshake(p))

	// read handshake
	buffer := make([]byte, p2p.MaxPayloadSize)
	r, err := pConn.Read(buffer)
	assert.Nil(t, err)

	// check size
	assert.Equal(t, r, noise.SizeMessageInitializer,
		"unexpected initializer message size")

	assert.Nil(t, n.Stop())
}

func TestNodeHandshakeInvalid(t *testing.T) {
	n := node.NewNode("test", zeroKey)
	p := peer.NewPeer(zeroKey.PublicKey(), "")

	// node not running or connected
	assert.NotNil(t, n.Handshake(p),
		"handshake while node not connected") // err 1

	// start and connect
	assert.Nil(t, n.Start())
	assert.Nil(t, n.Connect())

	// invalid peer
	assert.NotNil(t, n.Handshake(nil), "peer is nil")       // err 2
	assert.NotNil(t, n.Handshake(p), "peer address is nil") // err 3

	// valid peer
	p = peer.NewPeer(zeroKey.PublicKey(), "localhost:")

	// fake Handshake
	p.Handshake = new(noise.Handshake)
	assert.NotNil(t, n.Handshake(p), "peer has handshake") // err 4
	p.Handshake = nil

	// attempt handshake
	assert.Nil(t, n.Handshake(p))

	// stop
	assert.Nil(t, n.Stop())

	// check error count
	assertErrCount(t, n, 4)
}

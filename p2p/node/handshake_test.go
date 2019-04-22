package node_test

import (
	"net"
	"testing"

	"cpl.li/go/cryptor/crypt/ppk"
	"cpl.li/go/cryptor/p2p"
	"cpl.li/go/cryptor/p2p/node"
	"cpl.li/go/cryptor/p2p/noise"
	"cpl.li/go/cryptor/p2p/peer"
	"cpl.li/go/cryptor/tests"
)

func TestNodeHandshake(t *testing.T) {
	// create node keys
	peerPrivate, _ := ppk.NewPrivateKey()
	peerPublic := peerPrivate.PublicKey()

	// create nodes
	n := node.NewNode("test", zeroKey)

	// start node
	tests.AssertNil(t, n.Start())

	// set address
	tests.AssertNil(t, n.SetAddr("localhost:"))

	// add external peer
	p, err := n.PeerAdd(peerPublic, "", 1)
	tests.AssertNil(t, err)

	// connect
	tests.AssertNil(t, n.Connect())

	// check peer count
	if count := n.PeerCount(); count != 1 {
		t.Fatalf("node peer count != 1, got %d\n", count)
	}

	// set peer addr
	tests.AssertNil(t, p.SetAddr("localhost:12345"))

	// simulate peer connection
	pConn, err := net.ListenUDP(p2p.Network, p.AddrUDP())
	tests.AssertNil(t, err)

	// attempt handshake with peer
	tests.AssertNil(t, n.Handshake(p))

	// read handshake
	buffer := make([]byte, p2p.MaxPayloadSize)
	r, err := pConn.Read(buffer)
	tests.AssertNil(t, err)

	// check size
	if r != noise.SizeMessageInitializer {
		t.Fatalf("expected message size %d, got %d\n",
			noise.SizeMessageInitializer, r)
	}

	tests.AssertNil(t, n.Stop())
}

func TestNodeHandshakeInvalid(t *testing.T) {
	n := node.NewNode("test", zeroKey)
	p := peer.NewPeer(zeroKey.PublicKey(), "")

	// node not running or connected
	tests.AssertNotNil(t, n.Handshake(p),
		"handshake while node not connected") // err 1

	// start and connect
	tests.AssertNil(t, n.Start())
	tests.AssertNil(t, n.Connect())

	// invalid peer
	tests.AssertNotNil(t, n.Handshake(nil), "peer is nil")       // err 2
	tests.AssertNotNil(t, n.Handshake(p), "peer address is nil") // err 3

	// valid peer
	p = peer.NewPeer(zeroKey.PublicKey(), "localhost:")

	// fake Handshake
	p.Handshake = new(noise.Handshake)
	tests.AssertNotNil(t, n.Handshake(p), "peer has handshake") // err 4
	p.Handshake = nil

	// attempt handshake
	tests.AssertNil(t, n.Handshake(p))

	// stop
	tests.AssertNil(t, n.Stop())

	// check error count
	assertErrCount(t, n, 4)
}

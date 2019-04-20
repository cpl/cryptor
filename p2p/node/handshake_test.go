package node_test

import (
	"testing"

	"cpl.li/go/cryptor/crypt/ppk"
	"cpl.li/go/cryptor/p2p/node"
	"cpl.li/go/cryptor/p2p/noise"
	"cpl.li/go/cryptor/p2p/peer"
	"cpl.li/go/cryptor/tests"
)

func TestNodeHandshakeSimple(t *testing.T) {
	// create node keys
	sigmaKey, _ := ppk.NewPrivateKey()
	omegaKey, _ := ppk.NewPrivateKey()

	// create nodes
	sigma := node.NewNode("sigma", sigmaKey)
	omega := node.NewNode("omega", omegaKey)

	// start nodes
	tests.AssertNil(t, sigma.Start())
	defer sigma.Stop()
	tests.AssertNil(t, omega.Start())
	defer omega.Stop()

	tests.AssertNil(t, sigma.SetAddr("localhost:45000"))
	tests.AssertNil(t, omega.SetAddr("localhost:45001"))

	// add external peer omega
	p, err := sigma.PeerAdd(omega.PublicKey(), omega.Addr())
	tests.AssertNil(t, err)

	// list peers
	tests.AssertNil(t, sigma.PeerList()) // 1 peers
	tests.AssertNil(t, omega.PeerList()) // 0 peers

	// check peer count
	if count := sigma.PeerCount(); count != 1 {
		t.Fatalf("node sigma peer count != 1, got %d\n", count)
	}

	// check peer count
	if count := omega.PeerCount(); count != 0 {
		t.Fatalf("node omega peer count != 0, got %d\n", count)
	}

	// connect
	tests.AssertNil(t, sigma.Connect())
	tests.AssertNil(t, omega.Connect())

	// attempt handshake with peer
	tests.AssertNil(t, sigma.Handshake(p))

	for omega.PeerCount() != 1 {
	}

	// list peers
	sigma.PeerList()
	omega.PeerList()

	// check peer count
	if count := sigma.PeerCount(); count != 1 {
		t.Fatalf("node sigma peer count != 1, got %d\n", count)
	}
	// check peer count
	if count := omega.PeerCount(); count != 1 {
		t.Fatalf("node omega peer count != 1, got %d\n", count)
	}
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

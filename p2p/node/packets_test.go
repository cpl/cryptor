package node_test

import (
	"net"
	"testing"
	"time"

	"cpl.li/go/cryptor/crypt"
	"cpl.li/go/cryptor/crypt/ppk"
	"cpl.li/go/cryptor/p2p"
	"cpl.li/go/cryptor/p2p/node"
	"cpl.li/go/cryptor/p2p/noise"
	"cpl.li/go/cryptor/tests"
)

func testSendPacket(t *testing.T, conn *net.UDPConn, size uint) {
	_, err := conn.Write(crypt.RandomBytes(size))
	tests.AssertNil(t, err)
}

func testDialUDP(t *testing.T, n *node.Node) *net.UDPConn {
	// dial node
	nAddr, err := net.ResolveUDPAddr(p2p.Network, n.Addr())
	tests.AssertNil(t, err)
	nConn, err := net.DialUDP(p2p.Network, nil, nAddr)
	tests.AssertNil(t, err)

	return nConn
}

func testSetupNode(t *testing.T) (*node.Node, *net.UDPConn) {
	// create and start node
	n := node.NewNode("test", zeroKey)
	tests.AssertNil(t, n.Start())
	tests.AssertNil(t, n.SetAddr("localhost:"))
	tests.AssertNil(t, n.Connect())

	// dial node
	nConn := testDialUDP(t, n)

	// check error count
	assertErrCount(t, n, 0)

	return n, nConn
}

/*

TODO

! While writing this test, realize that the protocol is lacking when an address
! change detection mechanism. This will result in a diffrent approach to
! developing the p2p package and accelerate the implementation of the Cryptor
! transport layer protocol.

func TestNodePayloadKnownValid(t *testing.T) {
	// test setup
	n, nConn := testSetupNode(t)

	// initialize peer keys
	secret, _ := ppk.NewPrivateKey()
	public := secret.PublicKey()

	// add peer to node and initialize handshake
	p, err := n.PeerAdd(public, nConn.LocalAddr().String())
	tests.AssertNil(t, err)
	var msg *noise.MessageInitializer
	p.Handshake, msg = noise.Initialize(n.PublicKey(), public)

	// attempt handshake bypass with invalid packets
	testSendPacket(t, nConn, p2p.MinPayloadSize-1)
	testSendPacket(t, nConn, 0)
	testSendPacket(t, nConn, 1)
	testSendPacket(t, nConn, p2p.MaxPayloadSize*2)
	testSendPacket(t, nConn, p2p.MinPayloadSize+1)

	// check error count
	assertErrCount(t, n, 0)

	// wrap things up
	tests.AssertNil(t, nConn.Close())
	tests.AssertNil(t, n.Stop())

	// check error count
	assertErrCount(t, n, 0)
}
*/

func TestNodePayloadUnknownValid(t *testing.T) {
	// test setup
	n, nConn := testSetupNode(t)

	// initialize handshake message
	secret, _ := ppk.NewPrivateKey()
	public := secret.PublicKey()
	_, msg := noise.Initialize(public, n.PublicKey())

	// send msg
	payload, err := msg.MarshalBinary()
	tests.AssertNil(t, err)
	if _, err := nConn.Write(payload); err != nil {
		t.Fatal(err)
	}

	// wait for response
	readBuffer := make([]byte, 1024)
	r, _, err := nConn.ReadFrom(readBuffer)
	tests.AssertNil(t, err)

	// check response
	if r != noise.SizeMessageResponder {
		t.Fatalf("invalid handshake response size, expected %d, got %d\n",
			noise.SizeMessageResponder, r)
	}

	// check peer count of node
	if count := n.PeerCount(); count != 1 {
		t.Fatalf("invalid peer count, expected %d, got %d\n", 1, count)
	}

	// check peer
	p := n.PeerGet(public)
	if p == nil {
		t.Fatal("peer not found")
	}

	// check peer key
	if !p.PublicKey().Equals(public) {
		t.Fatal("peer public key does not match")
	}

	// check error count
	assertErrCount(t, n, 0)

	// attempt to duplicate peer from diffrent address
	nConn = testDialUDP(t, n)
	// send msg
	if _, err := nConn.Write(payload); err != nil {
		t.Fatal(err)
	}

	// wait
	time.Sleep(2 * time.Second)
	assertErrCount(t, n, 1)

	// wrap things up
	tests.AssertNil(t, nConn.Close())
	tests.AssertNil(t, n.Stop())

	// check error count
	assertErrCount(t, n, 1)
}

func TestNodePayloadUnknownInvalid(t *testing.T) {
	// test setup
	n, nConn := testSetupNode(t)

	// send invalid packets (not counted as errors)
	testSendPacket(t, nConn, p2p.MinPayloadSize-1)
	testSendPacket(t, nConn, 0)
	testSendPacket(t, nConn, 1)
	testSendPacket(t, nConn, p2p.MaxPayloadSize*2)
	testSendPacket(t, nConn, p2p.MinPayloadSize+1)

	// send valid packet with invalid payloads
	for i := 0; i < 5; i++ {
		testSendPacket(t, nConn, noise.SizeMessageInitializer) // 5x err
	}

	time.Sleep(2 * time.Second)

	// check error count
	assertErrCount(t, n, 5)

	// wrap things up
	tests.AssertNil(t, nConn.Close())
	tests.AssertNil(t, n.Stop())

	// check error count
	assertErrCount(t, n, 5)
}

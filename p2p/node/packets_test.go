package node_test

import (
	"net"
	"testing"
	"time"

	"cpl.li/go/cryptor/crypt/ppk"

	"cpl.li/go/cryptor/crypt"
	"cpl.li/go/cryptor/p2p"
	"cpl.li/go/cryptor/p2p/node"
	"cpl.li/go/cryptor/p2p/noise"
	"cpl.li/go/cryptor/tests"
)

func testRandomSendPacket(t *testing.T, conn *net.UDPConn, size uint) {
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

func TestPacketInvalidSimple(t *testing.T) {
	n, nConn := testSetupNode(t)

	// send invalid packets (not counted as errors)
	testRandomSendPacket(t, nConn, p2p.MinPayloadSize-1)
	testRandomSendPacket(t, nConn, 0)
	testRandomSendPacket(t, nConn, 1)
	testRandomSendPacket(t, nConn, p2p.MaxPayloadSize*2)
	testRandomSendPacket(t, nConn, p2p.MinPayloadSize+1)

	// send valid Initializer packets with invalid payloads
	for i := 0; i < 5; i++ {
		testRandomSendPacket(t, nConn, noise.SizeMessageInitializer) // 5x err
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

func TestPacketInvalidComplex(t *testing.T) {
	n, nConn := testSetupNode(t)

	// fake initializer
	iPrivate, _ := ppk.NewPrivateKey()
	iPublic := iPrivate.PublicKey()
	_, iMsg := noise.Initialize(iPublic, n.PublicKey())

	// send initializer message (id is 0) (invalid)
	iMsg.PeerID = 0
	iData0, _ := iMsg.MarshalBinary()
	nConn.Write(iData0)
	time.Sleep(2 * time.Second)

	// send initializer message (id is 1) (valid)
	iMsg.PeerID = 1
	iData1, _ := iMsg.MarshalBinary()
	nConn.Write(iData1)
	time.Sleep(2 * time.Second)

	// re-send same key (invalid)
	iMsg.PeerID = 2
	iData2, _ := iMsg.MarshalBinary()
	nConn.Write(iData2)
	time.Sleep(2 * time.Second)

	// check error count
	assertErrCount(t, n, 2)

	// wrap things up
	tests.AssertNil(t, nConn.Close())
	tests.AssertNil(t, n.Stop())

	// check error count
	assertErrCount(t, n, 2)
}

func TestPacketValidSimple(t *testing.T) {
	// setup test
	n, nConn := testSetupNode(t)
	addr, err := net.ResolveUDPAddr(p2p.Network, "localhost:")
	tests.AssertNil(t, err)
	pConn, err := net.ListenUDP(p2p.Network, addr)
	tests.AssertNil(t, err)
	buffer := make([]byte, p2p.MaxPayloadSize)

	// fake responder
	rPrivate, _ := ppk.NewPrivateKey()
	rPublic := rPrivate.PublicKey()

	// add peer to node
	p, err := n.PeerAdd(rPublic, pConn.LocalAddr().String(), 1)
	tests.AssertNil(t, err)

	// perform handshake initialization
	tests.AssertNil(t, n.Handshake(p))

	// receive initializer message
	r, err := pConn.Read(buffer)
	tests.AssertNil(t, err)

	// unpack initializer message
	iMsg := new(noise.MessageInitializer)
	tests.AssertNil(t, iMsg.UnmarshalBinary(buffer[:r]))

	// validate message
	if iMsg.PeerID != 1 {
		t.Fatalf("unexpected peer ID, expected %d, got %d\n", 1, iMsg.PeerID)
	}

	// perform response
	_, iPublic, rMsg, err := noise.Respond(iMsg, rPrivate)
	tests.AssertNil(t, err)

	// check matching initializer public key
	if !iPublic.Equals(n.PublicKey()) {
		t.Fatal("mismatch public key")
	}

	// assign ID
	rMsg.PeerID = 1

	// send message
	rData, _ := rMsg.MarshalBinary()
	if _, err := nConn.Write(rData); err != nil {
		t.Fatal(err)
	}

	time.Sleep(2 * time.Second)

	// check error count
	assertErrCount(t, n, 0)

	// wrap things up
	tests.AssertNil(t, nConn.Close())
	tests.AssertNil(t, n.Stop())

	// check error count
	assertErrCount(t, n, 0)
}

func TestPacketInvalidResponder(t *testing.T) {
	// setup test
	n, nConn := testSetupNode(t)
	addr, err := net.ResolveUDPAddr(p2p.Network, "localhost:")
	tests.AssertNil(t, err)
	pConn, err := net.ListenUDP(p2p.Network, addr)
	tests.AssertNil(t, err)
	buffer := make([]byte, p2p.MaxPayloadSize)

	// fake responder
	rPrivate, _ := ppk.NewPrivateKey()
	rPublic := rPrivate.PublicKey()

	// add peer to node
	p, err := n.PeerAdd(rPublic, pConn.LocalAddr().String(), 1)
	tests.AssertNil(t, err)

	// perform handshake initialization
	tests.AssertNil(t, n.Handshake(p))

	// receive initializer message
	r, err := pConn.Read(buffer)
	tests.AssertNil(t, err)

	// unpack initializer message
	iMsg := new(noise.MessageInitializer)
	tests.AssertNil(t, iMsg.UnmarshalBinary(buffer[:r]))

	// validate message
	if iMsg.PeerID != 1 {
		t.Fatalf("unexpected peer ID, expected %d, got %d\n", 1, iMsg.PeerID)
	}

	// perform bad response
	rMsg := new(noise.MessageResponder)
	rMsg.PeerID = 1

	// send message
	rData, _ := rMsg.MarshalBinary()
	if _, err := nConn.Write(rData); err != nil {
		t.Fatal(err)
	}
	time.Sleep(2 * time.Second)

	// invalid sizes messages
	invalid0 := make([]byte, noise.SizeMessageResponder-1)
	invalid0[0] = 1
	invalid1 := make([]byte, noise.SizeMessageResponder+1)
	invalid1[0] = 1

	// send messages
	if _, err := nConn.Write(invalid0); err != nil {
		t.Fatal(err)
	}
	if _, err := nConn.Write(invalid1); err != nil {
		t.Fatal(err)
	}

	time.Sleep(2 * time.Second)

	// check error count
	assertErrCount(t, n, 1)

	// wrap things up
	tests.AssertNil(t, nConn.Close())
	tests.AssertNil(t, n.Stop())

	// check error count
	assertErrCount(t, n, 1)
}

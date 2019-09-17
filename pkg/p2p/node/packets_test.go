package node_test

import (
	"net"
	"testing"
	"time"

	"cpl.li/go/cryptor/pkg/crypt/ppk"

	"cpl.li/go/cryptor/pkg/crypt"
	"cpl.li/go/cryptor/pkg/p2p"
	"cpl.li/go/cryptor/pkg/p2p/node"
	"cpl.li/go/cryptor/pkg/p2p/noise"
	"github.com/stretchr/testify/assert"
)

func testRandomSendPacket(t *testing.T, conn *net.UDPConn, size uint) {
	_, err := conn.Write(crypt.RandomBytes(size))
	assert.Nil(t, err)
}

func testDialUDP(t *testing.T, n *node.Node) *net.UDPConn {
	// dial node
	nAddr, err := net.ResolveUDPAddr(p2p.Network, n.Addr())
	assert.Nil(t, err)
	nConn, err := net.DialUDP(p2p.Network, nil, nAddr)
	assert.Nil(t, err)

	return nConn
}

func testSetupNode(t *testing.T) (*node.Node, *net.UDPConn) {
	// create and start node
	n := node.NewNode("test", zeroKey)
	assert.Nil(t, n.Start())
	assert.Nil(t, n.SetAddr("localhost:"))
	assert.Nil(t, n.Connect())

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
	assert.Nil(t, nConn.Close())
	assert.Nil(t, n.Stop())

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
	assert.Nil(t, nConn.Close())
	assert.Nil(t, n.Stop())

	// check error count
	assertErrCount(t, n, 2)
}

func assertConnWrite(t *testing.T, conn *net.UDPConn, data []byte, exp int) {
	r, err := conn.Write(data)
	assert.Nil(t, err)
	if exp > -1 && exp != r {
		assert.Equal(t, r, exp, "invalid conn write count")
	}
}

func TestPacketValidSimple(t *testing.T) {
	// setup test
	n, nConn := testSetupNode(t)
	addr, err := net.ResolveUDPAddr(p2p.Network, "localhost:")
	assert.Nil(t, err)
	pConn, err := net.ListenUDP(p2p.Network, addr)
	assert.Nil(t, err)
	buffer := make([]byte, p2p.MaxPayloadSize)

	// fake responder
	rPrivate, _ := ppk.NewPrivateKey()
	rPublic := rPrivate.PublicKey()

	// add peer to node
	p, err := n.PeerAdd(rPublic, pConn.LocalAddr().String(), 1)
	assert.Nil(t, err)

	// perform handshake initialization
	assert.Nil(t, n.Handshake(p))

	// receive initializer message
	r, err := pConn.Read(buffer)
	assert.Nil(t, err)

	// unpack initializer message
	iMsg := new(noise.MessageInitializer)
	assert.Nil(t, iMsg.UnmarshalBinary(buffer[:r]))

	// validate message
	assert.Equal(t, iMsg.PeerID, uint64(1), "unexpected peer ID")

	// perform response
	_, iPublic, rMsg, err := noise.Respond(iMsg, rPrivate)
	assert.Nil(t, err)

	// check matching initializer public key
	if !iPublic.Equals(n.PublicKey()) {
		t.Fatal("mismatch public key")
	}

	// assign ID
	rMsg.PeerID = 1

	// send message
	rData, _ := rMsg.MarshalBinary()
	assertConnWrite(t, nConn, rData, len(rData))

	time.Sleep(2 * time.Second)

	// check error count
	assertErrCount(t, n, 0)

	// wrap things up
	assert.Nil(t, nConn.Close())
	assert.Nil(t, n.Stop())

	// check error count
	assertErrCount(t, n, 0)
}

func TestPacketInvalidResponder(t *testing.T) {
	// setup test
	n, nConn := testSetupNode(t)
	addr, err := net.ResolveUDPAddr(p2p.Network, "localhost:")
	assert.Nil(t, err)
	pConn, err := net.ListenUDP(p2p.Network, addr)
	assert.Nil(t, err)
	buffer := make([]byte, p2p.MaxPayloadSize)

	// fake responder
	rPrivate, _ := ppk.NewPrivateKey()
	rPublic := rPrivate.PublicKey()

	// add peer to node
	p, err := n.PeerAdd(rPublic, pConn.LocalAddr().String(), 1)
	assert.Nil(t, err)

	// perform handshake initialization
	assert.Nil(t, n.Handshake(p))

	// receive initializer message
	r, err := pConn.Read(buffer)
	assert.Nil(t, err)

	// unpack initializer message
	iMsg := new(noise.MessageInitializer)
	assert.Nil(t, iMsg.UnmarshalBinary(buffer[:r]))

	// validate message
	assert.Equal(t, iMsg.PeerID, uint64(1), "unexpected peer ID")

	// perform bad response
	rMsg := new(noise.MessageResponder)
	rMsg.PeerID = 1

	// send message
	rData, _ := rMsg.MarshalBinary()
	assertConnWrite(t, nConn, rData, len(rData))

	time.Sleep(2 * time.Second)

	// invalid sizes messages
	invalid0 := make([]byte, noise.SizeMessageResponder-1)
	invalid0[0] = 1
	invalid1 := make([]byte, noise.SizeMessageResponder+1)
	invalid1[0] = 1

	// send messages
	assertConnWrite(t, nConn, invalid0, len(invalid0))
	assertConnWrite(t, nConn, invalid1, len(invalid1))

	time.Sleep(2 * time.Second)

	// check error count
	assertErrCount(t, n, 1)

	// wrap things up
	assert.Nil(t, nConn.Close())
	assert.Nil(t, n.Stop())

	// check error count
	assertErrCount(t, n, 1)
}

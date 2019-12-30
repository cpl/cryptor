package node_test

import (
	"testing"

	"cpl.li/go/cryptor/pkg/crypt/ppk"
	"cpl.li/go/cryptor/pkg/p2p/node"

	"github.com/stretchr/testify/assert"
)

var zeroKey ppk.PrivateKey

func assertErrCount(t *testing.T, n *node.Node, expected uint32) {
	// check error count
	assert.Equal(t, expected, n.ErrCount(), "unexpected error count")
}

func assertNodeAddr(t *testing.T, n *node.Node, expected string) {
	assert.Equal(t, n.Addr(), expected, "unexpected node address")
}

func TestNodeBasicRoutines(t *testing.T) {
	t.Parallel()

	n := node.NewNode("test", zeroKey)

	assert.Nil(t, n.Start())
	assert.Nil(t, n.Stop())

	assert.Nil(t, n.Start())
	assert.Nil(t, n.Connect())
	assert.Nil(t, n.Disconnect())
	assert.Nil(t, n.Stop())

	assert.Nil(t, n.Start())
	assert.Nil(t, n.Connect())
	assert.Nil(t, n.Stop())

	// check error count
	assertErrCount(t, n, 0)
}

func TestNodeBasicInvalidRoutines(t *testing.T) {
	t.Parallel()

	n := node.NewNode("test", zeroKey)

	assert.NotNil(t, n.Stop(), "stop: not running")       // 1 err count
	assert.NotNil(t, n.Disconnect(), "disc: not running") // 2
	assert.NotNil(t, n.Connect(), "conn: not running")    // 3

	assert.Nil(t, n.Start())
	assert.NotNil(t, n.Disconnect(), "disc: running") // 4
	assert.NotNil(t, n.Start(), "start: running")     // 5

	assert.Nil(t, n.Connect())
	assert.NotNil(t, n.Start(), "start: running connected")  // 6
	assert.NotNil(t, n.Connect(), "conn: running connected") // 7

	assert.Nil(t, n.Stop())

	// check error count
	assertErrCount(t, n, 7)
}

func TestNodeAddrSet(t *testing.T) {
	// create node and check default address
	n := node.NewNode("test", zeroKey)
	assertNodeAddr(t, n, "<nil>")

	// change address
	assert.Nil(t, n.SetAddr("127.0.0.1:8000"))

	// start node
	assert.Nil(t, n.Start())

	// change address while node is running
	assert.Nil(t, n.SetAddr(":"))

	// connect to random port on [::]/0.0.0.0
	assert.Nil(t, n.Connect())

	// change address while node is running and connected
	assert.NotNil(t, n.SetAddr("127.0.0.2:8001"),
		"changed addr while connected") // err 1

	// disconnect node
	assert.Nil(t, n.Disconnect())

	// attempt change to invalid addresses
	invalidAddresses := []string{ // + 5 err
		"8000",
		"127",
		"nil",
		"nosuchhostname",
		"127.0.0.1:noport",
	}
	for _, addr := range invalidAddresses {
		assert.NotNil(t, n.SetAddr(addr), "changed address to "+addr)
	}

	// change address
	assert.Nil(t, n.SetAddr("127.0.0.2:8001"))
	assert.NotNil(t, n.SetAddr("127.0.0.2:8001"),
		"change addr to same addr") // err 7

	// change address
	assert.Nil(t, n.SetAddr("127.0.0.1:40123"))

	// check address while not connected
	assertNodeAddr(t, n, "127.0.0.1:40123")

	// connect
	assert.Nil(t, n.Connect())

	// check address while connected
	assertNodeAddr(t, n, "127.0.0.1:40123")

	// stop node
	assert.Nil(t, n.Stop())

	// check error count
	assertErrCount(t, n, 7)

}

func TestConnectInvalidAddress(t *testing.T) {
	t.Parallel()

	n := node.NewNode("test", zeroKey)
	assert.Nil(t, n.SetAddr("example.com:80"))
	assert.Nil(t, n.Start())
	assert.NotNil(t, n.Connect(), "connect on invalid address")
	assert.Nil(t, n.Stop())
}

func TestNodeWait(t *testing.T) {
	t.Parallel()

	n := node.NewNode("test", zeroKey)
	go assert.Nil(t, n.Start())
	n.Wait()
	go assert.Nil(t, n.Stop())
	n.Wait()
	go assert.Nil(t, n.Start())
	n.Wait()
	go assert.NotNil(t, n.Start(), "did not fail to start again")
	n.Wait()
	go assert.Nil(t, n.Stop())
	n.Wait()
	assert.NotNil(t, n.Stop(), "did not fail to stop non-running node")
}

func assertNodeName(t *testing.T, name string) {
	n := node.NewNode(name, zeroKey)
	assert.Equal(t, n.Name(), name, "node name does not match")
}

func TestNodeName(t *testing.T) {
	t.Parallel()

	assertNodeName(t, "test")
	assertNodeName(t, " ")
	assertNodeName(t, "1234$$$$")
	assertNodeName(t, " multi   name  ")
	assertNodeName(t, "-==--")
	assertNodeName(t, "s")

	n := node.NewNode("", zeroKey)
	if n.Name() == "" {
		t.Fatal("expected mnemonic name, got nothing")
	}
	t.Log(n.Name())
}

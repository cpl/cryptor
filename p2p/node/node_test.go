package node_test

import (
	"testing"

	"cpl.li/go/cryptor/crypt/ppk"

	"cpl.li/go/cryptor/p2p/node"
	"cpl.li/go/cryptor/tests"
)

var zeroKey ppk.PrivateKey

func assertErrCount(t *testing.T, n *node.Node, expected int) {
	// check error count
	tests.AssertEqual(t, n.ErrCount(), expected, "unexpected error count")
}

func assertNodeAddr(t *testing.T, n *node.Node, expected string) {
	tests.AssertEqual(t, n.Addr(), expected, "unexpected node address")
}

func TestNodeBasicRoutines(t *testing.T) {
	t.Parallel()

	n := node.NewNode("test", zeroKey)

	tests.AssertNil(t, n.Start())
	tests.AssertNil(t, n.Stop())

	tests.AssertNil(t, n.Start())
	tests.AssertNil(t, n.Connect())
	tests.AssertNil(t, n.Disconnect())
	tests.AssertNil(t, n.Stop())

	tests.AssertNil(t, n.Start())
	tests.AssertNil(t, n.Connect())
	tests.AssertNil(t, n.Stop())

	// check error count
	assertErrCount(t, n, 0)
}

func TestNodeBasicInvalidRoutines(t *testing.T) {
	t.Parallel()

	n := node.NewNode("test", zeroKey)

	tests.AssertNotNil(t, n.Stop(), "stop: not running")       // 1 err count
	tests.AssertNotNil(t, n.Disconnect(), "disc: not running") // 2
	tests.AssertNotNil(t, n.Connect(), "conn: not running")    // 3

	tests.AssertNil(t, n.Start())
	tests.AssertNotNil(t, n.Disconnect(), "disc: running") // 4
	tests.AssertNotNil(t, n.Start(), "start: running")     // 5

	tests.AssertNil(t, n.Connect())
	tests.AssertNotNil(t, n.Start(), "start: running connected")  // 6
	tests.AssertNotNil(t, n.Connect(), "conn: running connected") // 7

	tests.AssertNil(t, n.Stop())

	// check error count
	assertErrCount(t, n, 7)
}

func TestNodeAddrSet(t *testing.T) {
	// create node and check default address
	n := node.NewNode("test", zeroKey)
	assertNodeAddr(t, n, "<nil>")

	// change address
	tests.AssertNil(t, n.SetAddr("127.0.0.1:8000"))

	// start node
	tests.AssertNil(t, n.Start())

	// change address while node is running
	tests.AssertNil(t, n.SetAddr(":"))

	// connect to random port on [::]/0.0.0.0
	tests.AssertNil(t, n.Connect())

	// change address while node is running and connected
	tests.AssertNotNil(t, n.SetAddr("127.0.0.2:8001"),
		"changed addr while connected") // err 1

	// disconnect node
	tests.AssertNil(t, n.Disconnect())

	// attempt change to invalid addresses
	invalidAddresses := []string{ // + 5 err
		"8000",
		"127",
		"nil",
		"nosuchhostname",
		"127.0.0.1:noport",
	}
	for _, addr := range invalidAddresses {
		tests.AssertNotNil(t, n.SetAddr(addr), "changed address to "+addr)
	}

	// change address
	tests.AssertNil(t, n.SetAddr("127.0.0.2:8001"))
	tests.AssertNotNil(t, n.SetAddr("127.0.0.2:8001"),
		"change addr to same addr") // err 7

	// change address
	tests.AssertNil(t, n.SetAddr("127.0.0.1:40123"))

	// check address while not connected
	assertNodeAddr(t, n, "127.0.0.1:40123")

	// connect
	tests.AssertNil(t, n.Connect())

	// check address while connected
	assertNodeAddr(t, n, "127.0.0.1:40123")

	// stop node
	tests.AssertNil(t, n.Stop())

	// check error count
	assertErrCount(t, n, 7)

}

func TestConnectInvalidAddress(t *testing.T) {
	t.Parallel()

	n := node.NewNode("test", zeroKey)
	tests.AssertNil(t, n.SetAddr("example.com:80"))
	tests.AssertNil(t, n.Start())
	tests.AssertNotNil(t, n.Connect(), "connect on invalid address")
	tests.AssertNil(t, n.Stop())
}

func TestNodeState(t *testing.T) {
	n := node.NewNode("test", zeroKey)
	tests.AssertEqual(t, n.State(), node.StateStopped,
		"unexpected node state")
	tests.AssertNil(t, n.Start())
	tests.AssertEqual(t, n.State(), node.StateRunning,
		"unexpected node state")
	tests.AssertNil(t, n.Stop())
	tests.AssertEqual(t, n.State(), node.StateStopped,
		"unexpected node state")
	tests.AssertNil(t, n.Start())
	tests.AssertNil(t, n.Connect())
	tests.AssertEqual(t, n.State(), node.StateConnected,
		"unexpected node state")
	tests.AssertNil(t, n.Stop())
	tests.AssertEqual(t, n.State(), node.StateStopped,
		"unexpected node state")
}

func TestNodeStateString(t *testing.T) {
	tests.AssertEqual(t, node.StateStopped.String(), "STOPPED",
		"unexpected state string")
	tests.AssertEqual(t, node.StateRunning.String(), "RUNNING",
		"unexpected state string")
	tests.AssertEqual(t, node.StateConnected.String(), "CONNECTED",
		"unexpected state string")
}

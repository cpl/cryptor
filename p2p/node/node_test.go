package node_test

import (
	"testing"

	"cpl.li/go/cryptor/crypt/ppk"

	"cpl.li/go/cryptor/p2p/node"
	"cpl.li/go/cryptor/tests"
)

var zeroKey ppk.PrivateKey

func TestNodeBasicRoutines(t *testing.T) {
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
	if count := n.ErrCount(); count != 0 {
		t.Fatalf("expected no errors, got %d\n", count)
	}
}

func TestNodeBasicInvalidRoutines(t *testing.T) {
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
	if count := n.ErrCount(); count != 7 {
		t.Fatalf("expected 7 errors, got %d\n", count)
	}
}

func TestNodeAddrSet(t *testing.T) {
	// create node and check default address
	n := node.NewNode("test", zeroKey)
	if addr := n.Addr(); addr != "<nil>" {
		t.Fatalf("unexpected node address, expected <nil>, got %s\n", addr)
	}

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
	if addr := n.Addr(); addr != "127.0.0.1:40123" {
		t.Fatalf(
			"unexpected node address, expected 127.0.0.1:40123, got %s\n", addr)
	}

	// connect
	tests.AssertNil(t, n.Connect())

	// check address while connected
	if addr := n.Addr(); addr != "127.0.0.1:40123" {
		t.Fatalf(
			"unexpected node address, expected 127.0.0.1:40123, got %s\n", addr)
	}

	// stop node
	tests.AssertNil(t, n.Stop())

	// check error count
	if count := n.ErrCount(); count != 7 {
		t.Fatalf("expected 7 errors, got %d\n", count)
	}
}

func TestConnectInvalidAddress(t *testing.T) {
	n := node.NewNode("test", zeroKey)
	tests.AssertNil(t, n.SetAddr("example.com:80"))
	tests.AssertNil(t, n.Start())
	tests.AssertNotNil(t, n.Connect(), "connect on invalid address")
	tests.AssertNil(t, n.Stop())
}

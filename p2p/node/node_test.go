package node_test

import (
	"testing"

	"cpl.li/go/cryptor/p2p/node"
	"cpl.li/go/cryptor/tests"
)

func TestNodeBasicRoutines(t *testing.T) {
	n := node.NewNode("test")

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
	n := node.NewNode("test")

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
		t.Fatalf("expected no errors, got %d\n", count)
	}
}

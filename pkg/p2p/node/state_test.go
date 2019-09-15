package node_test

import (
	"testing"

	"cpl.li/go/cryptor/p2p/node"
	"cpl.li/go/cryptor/tests"
)

func TestNodeState(t *testing.T) {
	t.Parallel()

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
	t.Parallel()

	tests.AssertEqual(t, node.StateStopped.String(), "STOPPED",
		"unexpected state string")
	tests.AssertEqual(t, node.StateRunning.String(), "RUNNING",
		"unexpected state string")
	tests.AssertEqual(t, node.StateConnected.String(), "CONNECTED",
		"unexpected state string")
}

package node

import (
	"testing"

	"cpl.li/go/cryptor/crypt/ppk"

	"cpl.li/go/cryptor/tests"
)

// TODO Add build tags to all test files for separating race cond risk tests

func TestKillNetwork(t *testing.T) {
	var zeroKey ppk.PrivateKey
	n := NewNode("test", zeroKey)

	tests.AssertNil(t, n.Start())
	tests.AssertNil(t, n.Connect())

	// external network close
	n.net.Lock()
	tests.AssertNil(t, n.net.conn.Close())
	n.net.Unlock()

	// wait for node to disconnect
	for n.state.isConnected {
	}

	// check node still running
	if !n.state.isRunning {
		t.Error("node is not in running state")
	}

	// attempt node disconnect
	tests.AssertNotNil(t, n.Disconnect(), "disc not connected")

	// attempt node start
	tests.AssertNotNil(t, n.Start(), "start already running")

	// check node address
	tests.AssertEqual(t, n.Addr(), "<nil>", "unexpected node address")

	// attempt re-connect and disconnect
	tests.AssertNil(t, n.Connect())
	tests.AssertNil(t, n.Disconnect())

	tests.AssertNil(t, n.Stop())
}

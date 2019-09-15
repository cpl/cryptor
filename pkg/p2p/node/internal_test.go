package node

import (
	"testing"

	"cpl.li/go/cryptor/pkg/crypt/ppk"

	"github.com/stretchr/testify/assert"
)

// TODO Add build tags to all test files for separating race cond risk tests

func TestKillNetwork(t *testing.T) {
	n := NewNode("test", ppk.PrivateKey{})

	assert.Nil(t, n.Start())
	assert.Nil(t, n.Connect())

	// external network close
	n.net.Lock()
	assert.Nil(t, n.net.conn.Close())
	n.net.Unlock()

	// wait for node to disconnect
	for n.state.isConnected {
	}

	// check node still running
	if !n.state.isRunning {
		t.Error("node is not in running state")
	}

	// attempt node disconnect
	assert.NotNil(t, n.Disconnect(), "disc not connected")

	// attempt node start
	assert.NotNil(t, n.Start(), "start already running")

	// check node address
	assert.Equal(t, n.Addr(), "<nil>", "unexpected node address")

	// attempt re-connect and disconnect
	assert.Nil(t, n.Connect())
	assert.Nil(t, n.Disconnect())

	assert.Nil(t, n.Stop())

	assert.Equal(t, n.ErrCount(), uint32(4), "unexpected error count")
}

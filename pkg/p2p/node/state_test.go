package node_test

import (
	"testing"

	"cpl.li/go/cryptor/pkg/p2p/node"

	"github.com/stretchr/testify/assert"
)

func TestNodeState(t *testing.T) {
	t.Parallel()

	n := node.NewNode("test", zeroKey)
	assert.Equal(t, n.State(), node.StateStopped,
		"unexpected node state")
	assert.Nil(t, n.Start())
	assert.Equal(t, n.State(), node.StateRunning,
		"unexpected node state")
	assert.Nil(t, n.Stop())
	assert.Equal(t, n.State(), node.StateStopped,
		"unexpected node state")
	assert.Nil(t, n.Start())
	assert.Nil(t, n.Connect())
	assert.Equal(t, n.State(), node.StateConnected,
		"unexpected node state")
	assert.Nil(t, n.Stop())
	assert.Equal(t, n.State(), node.StateStopped,
		"unexpected node state")
}

func TestNodeStateString(t *testing.T) {
	t.Parallel()

	assert.Equal(t, node.StateStopped.String(), "STOPPED",
		"unexpected state string")
	assert.Equal(t, node.StateRunning.String(), "RUNNING",
		"unexpected state string")
	assert.Equal(t, node.StateConnected.String(), "CONNECTED",
		"unexpected state string")
}

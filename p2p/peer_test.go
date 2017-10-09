package p2p_test

import (
	"testing"
	"time"

	"github.com/thee-engineer/cryptor/p2p"
)

func TestNodePeers(t *testing.T) {
	qc := make(chan struct{})

	// Create node
	node := p2p.NewNode("127.0.0.1", 2000, 9000, qc)
	go node.Start()

	time.Sleep(time.Second)

	// Add 20 peers
	for peerCount := 0; peerCount < 20; peerCount++ {
		node.AddPeer(p2p.NewPeer("127.0.0.1", 2100+peerCount, 9100+peerCount))
	}

	// Count peers
	if node.PeerCount() != 20 {
		t.Error("node: unexpected peer count")
	}

	if len(node.Peers()) != 20 {
		t.Error("node: unexpected peer count")
	}

	qc <- *new(struct{})
}

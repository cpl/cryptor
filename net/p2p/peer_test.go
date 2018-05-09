package p2p_test

import (
	"testing"
	"time"

	"github.com/thee-engineer/cryptor/net/p2p"
)

func assertPeerCount(real, expected int, t *testing.T) {
	if real != expected {
		t.Errorf("node: expceted %d got %d peer count", expected, real)
	}
}

func TestNodePeers(t *testing.T) {
	qc := make(chan struct{})

	// Create node
	node := p2p.NewNode("127.0.0.1", 2000, 9000, qc, nil)
	go node.Start()

	time.Sleep(time.Second)

	// Add 20 peers
	for peerCount := 0; peerCount < 20; peerCount++ {
		node.AddPeer(p2p.NewPeer("127.0.0.1", 2100+peerCount, 9100+peerCount))
	}

	// Count peers
	assertPeerCount(node.PeerCount(), 20, t)
	assertPeerCount(len(node.Peers()), 20, t)

	for _, peer := range node.Peers() {
		node.RemovePeer(peer)
	}

	assertPeerCount(node.PeerCount(), 0, t)
	assertPeerCount(len(node.Peers()), 0, t)
}

func TestNodeTrustedPeers(t *testing.T) {
	qc := make(chan struct{})

	// Create node config with 4 trusted peers
	config := &p2p.NodeConfig{
		TrustedPeers: []*p2p.Peer{
			p2p.NewPeer("127.0.0.2", 2061, 9061),
			p2p.NewPeer("127.0.0.3", 2062, 9062),
			p2p.NewPeer("127.0.0.4", 2063, 9063),
			p2p.NewPeer("127.0.0.5", 2064, 9064),
		},
	}

	// Create node with node config
	n0 := p2p.NewNode("127.0.0.1", 2060, 9060, qc, config)

	time.Sleep(time.Second)

	// Start node (peers should be assigned)
	go n0.Start()

	time.Sleep(time.Second)

	// Check peer count
	if n0.PeerCount() != 4 {
		t.Fail()
	}

	time.Sleep(time.Second)

	qc <- *new(struct{})
}

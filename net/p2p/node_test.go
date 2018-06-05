package p2p_test

import (
	"testing"
	"time"

	"github.com/thee-engineer/cryptor/net/p2p"
)

func TestNodeNotRunning(t *testing.T) {
	node := p2p.NewNode("localhost", "2000", nil)
	node.Disconnect()
	node.Stop()
	node.Connect()
	node.Disconnect()
	node.AddPeer(nil)
	node.DelPeer(nil)
	if node.CountPeers() != 0 {
		t.Fail()
	}
	if node.Peers() != nil {
		t.Fail()
	}
}

func TestNodeRunning(t *testing.T) {
	node := p2p.NewNode("localhost", "2000", nil)
	node.Start()
	node.Disconnect()
	node.Start()

	if lenPeers := len(node.Peers()); lenPeers != 0 {
		t.Errorf("lenPeers, expected 0, got %d", lenPeers)
	}

	if peerCount := node.CountPeers(); peerCount != 0 {
		t.Errorf("peerCount, expected 0, got %d", peerCount)
	}

	go node.Peers()
	go node.CountPeers()

	time.Sleep(time.Second)
	node.Stop()
}

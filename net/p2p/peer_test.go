package p2p_test

import (
	"testing"

	"github.com/thee-engineer/cryptor/crypt"
	"github.com/thee-engineer/cryptor/crypt/encode/b16"

	"github.com/thee-engineer/cryptor/net/p2p"
)

func newTestPeer() *p2p.Peer {
	return &p2p.Peer{
		PublicKey: b16.EncodeString(crypt.RandomData(32)),
		Address:   "testAddress",
	}
}

func TestPeerAdd(t *testing.T) {
	node := p2p.NewNode("localhost", "2000", nil)
	go node.Start()

	if lenPeers := len(node.Peers()); lenPeers != 0 {
		t.Errorf("lenPeers, expected 0, got %d", lenPeers)
	}

	if peerCount := node.CountPeers(); peerCount != 0 {
		t.Errorf("peerCount, expected 0, got %d", peerCount)
	}

	// Nomral peer add
	node.AddPeer(newTestPeer())
	node.AddPeer(newTestPeer())
	node.AddPeer(newTestPeer())
	node.AddPeer(newTestPeer())

	if lenPeers := len(node.Peers()); lenPeers != 4 {
		t.Errorf("lenPeers, expected 4, got %d", lenPeers)
	}

	if peerCount := node.CountPeers(); peerCount != 4 {
		t.Errorf("peerCount, expected 4, got %d", peerCount)
	}

	// Peer that will be deleted
	dp0 := newTestPeer()
	dp1 := newTestPeer()

	// Concurent peer add
	node.AddPeer(dp0)
	node.AddPeer(dp1)
	node.AddPeer(newTestPeer())
	node.AddPeer(newTestPeer())

	if lenPeers := len(node.Peers()); lenPeers != 8 {
		t.Errorf("lenPeers, expected 8, got %d", lenPeers)
	}

	if peerCount := node.CountPeers(); peerCount != 8 {
		t.Errorf("peerCount, expected 8, got %d", peerCount)
	}

	node.DelPeer(dp0)
	node.DelPeer(dp1)

	if lenPeers := len(node.Peers()); lenPeers != 6 {
		t.Errorf("lenPeers, expected 8, got %d", lenPeers)
	}

	if peerCount := node.CountPeers(); peerCount != 6 {
		t.Errorf("peerCount, expected 8, got %d", peerCount)
	}
}

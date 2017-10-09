package p2p_test

import (
	"testing"
	"time"

	"github.com/thee-engineer/cryptor/p2p"
)

func TestNodeDefault(t *testing.T) {
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

func TestNodeStop(t *testing.T) {
	qc := make(chan struct{})
	node := p2p.NewNode("127.0.0.1", 2002, 9002, qc)
	go node.Stop() // Stop before starting

	go node.Start()
	time.Sleep(time.Second)

	go node.Stop()
}

func TestNodeStart(t *testing.T) {
	qc := make(chan struct{})

	node := p2p.NewNode("127.0.0.1", 2001, 9001, qc)

	go node.Start()

	time.Sleep(time.Second)

	qc <- *new(struct{})
}

func TestNodeConnection(t *testing.T) {
	qc := make(chan struct{})

	n0 := p2p.NewNode("127.0.0.1", 2010, 9010, qc)
	n1 := p2p.NewNode("127.0.0.1", 2011, 9011, qc)

	n0.Send(p2p.NewUDPPacket([]byte("hello world"), n0.UDPAddr()))

	time.Sleep(time.Second)

	go n0.Start()
	go n1.Start()

	time.Sleep(time.Second)

	qc <- *new(struct{})
}

package p2p_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/thee-engineer/cryptor/network/p2p"
)

func TestNetwork(t *testing.T) {
	qc := make(chan struct{})

	n0 := p2p.NewNode("localhost", 2000, qc)

	go n0.Start()

	// time.Sleep(2000 * time.Millisecond)

	for port := 9990; port < 10000; port++ {
		go n0.AddPeer("localhost", port)
	}

	go n0.AddPeer("localhost", 9990)

	time.Sleep(2000 * time.Millisecond)

	fmt.Println(n0.PeerCount())
	fmt.Println(n0.Peers())

	// time.Sleep(2000 * time.Millisecond)

}

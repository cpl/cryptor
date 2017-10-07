package p2p_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/thee-engineer/cryptor/network/p2p"
)

func TestNetwork(t *testing.T) {
	qc0 := make(chan struct{})
	qc1 := make(chan struct{})

	fmt.Println("making nodes")
	n0 := p2p.NewNode("localhost", 2000, qc0)
	n1 := p2p.NewNode("localhost", 2001, qc1)

	fmt.Println("starting nodes")
	go n0.Start()
	go n1.Start()

	fmt.Println("adding peers")
	for port := 9000; port < 9020; port++ {
		go n0.AddPeer(p2p.NewPeer("localhost", port))
	}

	time.Sleep(2000 * time.Millisecond)

	fmt.Println("getting count")
	fmt.Println(n0.PeerCount())
	fmt.Println(n1.PeerCount())

	<-qc0
}

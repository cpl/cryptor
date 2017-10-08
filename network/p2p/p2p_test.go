package p2p_test

import (
	"fmt"
	"testing"

	"github.com/thee-engineer/cryptor/network/p2p"
)

func TestNetwork(t *testing.T) {
	qc0 := make(chan struct{})

	fmt.Println("making node")
	n0 := p2p.NewNode("localhost", 2000, qc0)

	fmt.Println("starting node")
	go n0.Start()

	<-qc0
}

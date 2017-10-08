package p2p_test

import (
	"testing"

	"github.com/thee-engineer/cryptor/network/p2p"
)

func TestNetwork(t *testing.T) {
	qc := make(chan struct{})

	n0 := p2p.NewNode("127.0.0.1", 2000, qc)
	n0.Start()

	<-qc
}

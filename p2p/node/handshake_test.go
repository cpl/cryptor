package node_test

import (
	"testing"
	"time"

	"cpl.li/go/cryptor/p2p/node"
)

func TestNodeHandshakeDEBUG(t *testing.T) {
	// create nodes
	sigma := node.NewNode("sigma")
	omega := node.NewNode("omega")

	// start nodes
	sigma.Start()
	defer sigma.Stop()
	omega.Start()
	defer omega.Stop()

	sigma.SetAddr("localhost:45000")
	omega.SetAddr("localhost:45001")

	// add external peer
	p, err := sigma.PeerAdd(omega.PublicKey(), omega.Addr())
	if err != nil {
		t.Fatal(err)
	}

	// list peers
	sigma.PeerList()

	// connect
	sigma.Connect()
	omega.Connect()

	// attempt handshake with peer
	sigma.Handshake(p)

	time.Sleep(5 * time.Second)
}

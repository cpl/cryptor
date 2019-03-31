package p2p

import (
	"sync"

	"cpl.li/go/cryptor/crypt/ppk"
)

// Node represents the local machine running and/or connected to the Cryptor
// network. Other nodes are represented as peers.
type Node struct {

	// static identity of a node
	static struct {
		sync.RWMutex

		privateKey ppk.PrivateKey
		publicKey  ppk.PublicKey
	}
}

// NewNode creates a node running on the local machine. The default starting
// state is NOT RUNNING and OFFLINE. Allowing the Node to be further configured
// before starting and connecting to the Cryptor Network.
func NewNode() *Node {
	return nil
}

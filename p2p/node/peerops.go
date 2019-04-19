package node

import (
	"errors"

	"cpl.li/go/cryptor/crypt/ppk"
	"cpl.li/go/cryptor/p2p/peer"
)

// TODO Remove running state check.

// PeerAdd takes the public key of a peer and creates a new entry in the node
// map. An optional string address can be passed.
func (n *Node) PeerAdd(pk ppk.PublicKey, addr string) (*peer.Peer, error) {
	n.lookup.Lock()
	defer n.lookup.Unlock()

	// check node is running
	if !n.state.isRunning {
		return nil, errors.New("can't add peer, node is not running")
	}

	// check if peer already exists
	if _, ok := n.lookup.peers[pk]; ok {
		return nil, errors.New("can't add peer, public key already in use")
	}

	// create new peer and map it
	n.lookup.peers[pk] = peer.NewPeer(pk, addr)

	return n.lookup.peers[pk], nil
}

// PeerDel removes the peer from the peers map.
func (n *Node) PeerDel(pk ppk.PublicKey) (err error) {
	n.lookup.Lock()
	defer n.lookup.Unlock()

	// check node is running
	if !n.state.isRunning {
		return errors.New("can't remove peer, node is not running")
	}

	// check if peer exists
	if _, ok := n.lookup.peers[pk]; !ok {
		return errors.New("can't remove peer, public key not matching")
	}
	delete(n.lookup.peers, pk)

	return nil
}

// PeerCount returns the number of counts present in the peers map.
func (n *Node) PeerCount() int {
	// check node is running
	if !n.state.isRunning {
		return -1
	}

	return len(n.lookup.peers)
}

// PeerList prints all the peers public keys and address to the node logger.
func (n *Node) PeerList() error {
	// check node is running
	if !n.state.isRunning {
		return errors.New("can't get peer list, node is not running")
	}

	// iterate peers
	count := 0
	for pk, p := range n.lookup.peers {
		// print public key as hex
		n.logger.Printf("%4d: %s %s\n", count, pk.ToHex(), p.Addr())

		count++
	}

	return nil
}

// PeerGet returns a peer from lookup if it exists. If the peer does not exist
// a nil pointer is returned.
func (n *Node) PeerGet(pk ppk.PublicKey) *peer.Peer {
	// lock lookup
	n.lookup.Lock()
	defer n.lookup.Unlock()

	// check if peer exists, return nil if not
	p, ok := n.lookup.peers[pk]
	if !ok {
		return nil
	}

	// return peer
	return p
}

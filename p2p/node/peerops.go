package node

import (
	"errors"

	"cpl.li/go/cryptor/crypt/ppk"
	"cpl.li/go/cryptor/p2p/peer"
)

// PeerAdd takes the public key of a peer and creates a new entry in the node
// map. It returns the Peer struct and an error.
func (n *Node) PeerAdd(pk ppk.PublicKey) (p *peer.Peer, err error) {
	n.lookup.Lock()
	defer n.lookup.Unlock()

	// check node is running
	if !n.state.isRunning {
		err = errors.New("can't add peer, node is not running")
		n.comm.err <- err
		return nil, err
	}

	// check if peer already exists
	if _, ok := n.lookup.peers[pk]; ok {
		err = errors.New("can't add peer, public key already in use")
		n.comm.err <- err
		return nil, err
	}

	// create new peer
	p = peer.New(pk)

	// map peer
	n.lookup.peers[pk] = p

	return p, nil
}

// PeerDel removes the peer from the peers map.
func (n *Node) PeerDel(pk ppk.PublicKey) (err error) {
	n.lookup.Lock()
	defer n.lookup.Unlock()

	// check node is running
	if !n.state.isRunning {
		err = errors.New("can't remove peer, node is not running")
		n.comm.err <- err
		return err
	}

	// check if peer exists
	if _, ok := n.lookup.peers[pk]; !ok {
		err = errors.New("can't remove peer, public key not matching")
		n.comm.err <- err
		return err
	}
	delete(n.lookup.peers, pk)

	return nil
}

// PeerCount returns the number of counts present in the peers map.
func (n *Node) PeerCount() int {

	// check node is running
	if !n.state.isRunning {
		err := errors.New("can't get peer count, node is not running")
		n.comm.err <- err
		return -1
	}

	return len(n.lookup.peers)
}

// PeerList prints all the peers public keys and address to the node logger.
func (n *Node) PeerList() {
	// check node is running
	if !n.state.isRunning {
		err := errors.New("can't get peer list, node is not running")
		n.comm.err <- err
		return
	}

	// iterate peers
	count := 0
	for pk, p := range n.lookup.peers {
		// print public key as hex
		n.logger.Printf("%4d: %s %s\n", count, pk.ToHex(), p.Addr())

		count++
	}
}

// PeerSetAddr sets an existing peer address to the given one.
func (n *Node) PeerSetAddr(pk ppk.PublicKey, addr string) (err error) {
	// check node is running
	if !n.state.isRunning {
		err = errors.New("can't set peer address, node is not running")
		n.comm.err <- err
		return err
	}

	// lookup peer
	p, ok := n.lookup.peers[pk]
	if !ok {
		err = errors.New("can't set peer address, peer not found")
		n.comm.err <- err
		return err
	}

	// set addr, return any errors if it fails
	err = p.SetAddr(addr)
	if err != nil {
		n.comm.err <- err
	}

	return err
}

package node

import (
	"errors"
	"net"

	"cpl.li/go/cryptor/crypt/ppk"
	"cpl.li/go/cryptor/p2p"
)

// PeerAdd takes the public key of a peer and creates a new entry in the node
// map. An optional string address can be passed.
func (n *Node) PeerAdd(pk ppk.PublicKey, addr string) (p *Peer, err error) {
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

	// create new peer
	p = NewPeer(pk, addr)

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
		n.logger.Printf("%4d: %s %s\n", count, pk.ToHex(), p.addr.String())

		count++
	}
}

// PeerSetAddr sets an existing peer address to the given one.
func (n *Node) PeerSetAddr(pk ppk.PublicKey, addr string) (err error) {
	// check node is running
	if !n.state.isRunning {
		return err
	}

	// lookup peer
	p, ok := n.lookup.peers[pk]
	if !ok {
		return errors.New("can't set peer address, peer not found")
	}

	// only set addr if no errors
	newaddr, err := net.ResolveUDPAddr(p2p.Network, addr)
	if err != nil {
		return err
	}
	p.addr = newaddr

	return nil
}

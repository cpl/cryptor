package node

import (
	"errors"

	"cpl.li/go/cryptor/crypt/ppk"
	"cpl.li/go/cryptor/p2p/peer"
)

// PeerAdd takes the public key of a peer and creates a new entry in the node
// map. An optional string address can be passed.
func (n *Node) PeerAdd(pk ppk.PublicKey, addr string, id uint64) (*peer.Peer, error) {
	n.lookup.Lock()
	defer n.lookup.Unlock()

	// check if peer already exists (public key)
	if _, ok := n.lookup.peers[pk]; ok {
		return nil, errors.New("can't add peer, public key already in use")
	}

	// check if id is 0
	if id == 0 {
		return nil, errors.New("can't add peer, id is 0")
	}

	// check if peer already exists (id)
	if _, ok := n.lookup.table[id]; ok {
		return nil, errors.New("can't add peer, id already in use")
	}

	// create new peer
	p := peer.NewPeer(pk, addr)
	p.ID = id

	// assign peer to public key lookup
	n.lookup.peers[pk] = p

	// assign peer to id lookup
	n.lookup.table[id] = p

	// increment count
	n.lookup.count++

	return p, nil
}

// ! delete both peer entries from the lookup maps without any validation
func (n *Node) unsafePeerDel(id uint64, pk ppk.PublicKey) {
	// delete peers
	delete(n.lookup.table, id)
	delete(n.lookup.peers, pk)

	// decrement count
	n.lookup.count--
}

// PeerDelID removes the peer from the lookup maps based on the id.
func (n *Node) PeerDelID(id uint64) error {
	n.lookup.Lock()
	defer n.lookup.Unlock()

	// check peer count
	if n.lookup.count == 0 {
		return errors.New("can't remove peer, lookup is empty")
	}

	// check if peer exists
	p0, ok := n.lookup.table[id]
	if !ok {
		return errors.New("can't remove peer, id not found")
	}
	p1, ok := n.lookup.peers[p0.PublicKey()]
	if !ok {
		return errors.New("can't remove peer, public key not found")
	}
	if p0 != p1 {
		return errors.New("peers found, but don't match")
	}

	// delete while keeping lookup locked
	n.unsafePeerDel(id, p0.PublicKey())
	return nil
}

// PeerDel removes the peer from the lookup maps based on the public key.
func (n *Node) PeerDel(pk ppk.PublicKey) error {
	n.lookup.Lock()
	defer n.lookup.Unlock()

	// check peer count
	if n.lookup.count == 0 {
		return errors.New("can't remove peer, lookup is empty")
	}

	// check if peer exists
	p0, ok := n.lookup.peers[pk]
	if !ok {
		return errors.New("can't remove peer, public key not found")
	}
	p1, ok := n.lookup.table[p0.ID]
	if !ok {
		return errors.New("can't remove peer, id not found")
	}
	if p0 != p1 {
		return errors.New("peers found, but don't match")
	}

	// delete while keeping lookup locked
	n.unsafePeerDel(p0.ID, pk)
	return nil
}

// PeerCount returns the number of counts present in the peers lookup.
func (n *Node) PeerCount() int {
	n.lookup.Lock()
	defer n.lookup.Unlock()
	return n.lookup.count
}

// PeerList prints all the peers public keys and address to the node logger.
func (n *Node) PeerList() error {
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
// a nil pointer is returned. If ID is 0, a public key lookup is done.
func (n *Node) PeerGet(id uint64, pk ppk.PublicKey) *peer.Peer {
	// lock lookup
	n.lookup.Lock()
	defer n.lookup.Unlock()

	// check if ID is non-zero, perform lookup and return nil or peer pointer
	if id != 0 {
		p, ok := n.lookup.table[id]
		if !ok {
			return nil
		}
		return p
	}

	// check if peer exists, return nil if not
	p, ok := n.lookup.peers[pk]
	if !ok {
		return nil
	}

	// return peer
	return p
}

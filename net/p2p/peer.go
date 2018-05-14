package p2p

import (
	"crypto/rsa"

	"github.com/thee-engineer/cryptor/crypt/encode/b16"
	"github.com/thee-engineer/cryptor/crypt/hashing"
	"github.com/thee-engineer/cryptor/crypt/ppk"
)

// Peer ...
type Peer struct {
	PublicKey *rsa.PublicKey
	Address   string
}

// NewPeer ...
func NewPeer(address string, key *rsa.PublicKey) *Peer {
	return &Peer{
		PublicKey: key,
		Address:   address,
	}
}

// Encrypt ...
func (p *Peer) Encrypt(data []byte) ([]byte, error) {
	return ppk.Encrypt(p.PublicKey, data)
}

// Verify ...
func (p *Peer) Verify(data, signature []byte) bool {
	return ppk.Verify(p.PublicKey, data, signature)
}

// Hash ...
func (p *Peer) Hash() string {
	return b16.EncodeString(hashing.Hash([]byte(p.Address)))
}

type peerMap map[string]*Peer
type peerFunc func(peerMap)

// CountPeers ...
func (n *Node) CountPeers() int {
	if !checkRunning(n) {
		return 0
	}

	var count int

	select {
	case n.peerOp <- func(peerList peerMap) {
		count = len(peerList)
	}:
		<-n.peerOpDone
	}

	return count
}

// Peers ...
func (n *Node) Peers() []*Peer {
	if !checkRunning(n) {
		return nil
	}

	var peerList []*Peer

	select {
	case n.peerOp <- func(peers peerMap) {
		for _, peer := range peers {
			peerList = append(peerList, peer)
		}
	}:
		<-n.peerOpDone
	}

	return peerList
}

// AddPeer ...
func (n *Node) AddPeer(peer *Peer) {
	if !checkRunning(n) {
		return
	}

	select {
	case n.peerOp <- func(peers peerMap) {
		peers[peer.Address] = peer
	}:
		<-n.peerOpDone
	}
	n.logChan <- "add peer: " + peer.Hash()
}

// DelPeer ...
func (n *Node) DelPeer(peer *Peer) {
	if !checkRunning(n) {
		return
	}

	select {
	case n.peerOp <- func(peers peerMap) {
		delete(peers, peer.Address)
	}:
		<-n.peerOpDone
	}
	n.logChan <- "del peer: " + peer.Hash()
}

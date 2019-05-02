package node

import (
	"errors"
	"sync/atomic"

	"cpl.li/go/cryptor/crypt"
	"cpl.li/go/cryptor/p2p/noise"
	"cpl.li/go/cryptor/p2p/packet"
	"cpl.li/go/cryptor/p2p/peer"
)

// Handshake will initialize the handshake protocol with the given peer.
func (n *Node) Handshake(p *peer.Peer) (err error) {
	// check if node is connected
	n.state.RLock()
	if !n.state.isConnected {
		n.state.RUnlock()
		atomic.AddUint32(&n.meta.errCount, 1)
		return errors.New("can't handshake peer, node not connected")
	}
	n.state.RUnlock()

	// check for nil peer
	if p == nil {
		atomic.AddUint32(&n.meta.errCount, 1)
		return errors.New("peer is nil")
	}

	// check if peer is missing address
	if p.AddrUDP() == nil {
		atomic.AddUint32(&n.meta.errCount, 1)
		return errors.New("peer address is nil")
	}

	// check if handshake is initialized
	// TODO Extend this case to allow retries
	p.RLock()
	if p.Handshake != nil {
		p.RUnlock()
		atomic.AddUint32(&n.meta.errCount, 1)
		return errors.New("peer handshake already exists")
	}
	p.RUnlock()

	// lock peer
	p.Lock()

	// create handshake
	var msg *noise.MessageInitializer
	p.Handshake, msg = noise.Initialize(
		n.identity.publicKey, p.PublicKey())

	// assign ID
	if p.ID == 0 {
		p.ID = crypt.RandomUint64()
	}
	msg.PeerID = p.ID

	// unlock peer
	p.Unlock()

	// create packet
	pack := new(packet.Packet)
	pack.Type = packet.TypeHandshakeInitializer
	pack.Address = p.AddrUDP()

	// prepare handshake message
	pack.Payload, err = msg.MarshalBinary()
	if err != nil {
		atomic.AddUint32(&n.meta.errCount, 1)
		return err
	}

	// send handshake
	go n.send(pack)

	return nil
}

package node

import (
	"errors"

	"cpl.li/go/cryptor/p2p/noise"
	"cpl.li/go/cryptor/p2p/packet"
	"cpl.li/go/cryptor/p2p/peer"
)

// Handshake will initialize the handshake protocol with the given peer.
func (n *Node) Handshake(p *peer.Peer) (err error) {
	// check if node is connected
	if !n.state.isConnected {
		n.meta.errCount++
		return errors.New("can't handshake peer, node not connected")
	}

	// check for nil peer
	if p == nil {
		n.meta.errCount++
		return errors.New("peer is nil")
	}

	// lock peer
	p.Lock()
	defer p.Unlock()

	// check if handshake is initialized
	if p.Handshake != nil {
		n.meta.errCount++
		return errors.New("peer handshake already exists")
	}

	// create handshake
	var msg *noise.MessageInitializer
	p.Handshake, msg = noise.Initialize(
		n.identity.publicKey, p.PublicKey())

	// create packet
	pack := new(packet.Packet)
	pack.Type = packet.TypeHandshakeInitializer
	pack.Address = p.AddrUDP()

	// check if peer is missing address
	if pack.Address == nil {
		n.meta.errCount++
		return errors.New("peer address is nil")
	}

	// prepare handshake message
	pack.Payload, err = msg.MarshalBinary()
	if err != nil {
		n.meta.errCount++
		return err
	}

	// send handshake
	go n.send(pack)

	return nil
}

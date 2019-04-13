package node

import (
	"cpl.li/go/cryptor/p2p/noise"
	"cpl.li/go/cryptor/p2p/packet"
	"cpl.li/go/cryptor/p2p/peer"
)

// Handshake will initialize the handshake protocol with the given peer.
func (n *Node) Handshake(p *peer.Peer) {
	p.Lock()
	defer p.Unlock()

	// check if handshake is initialized
	if !p.HasHandshake {
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
			return
		}

		// prepare handshake message
		pack.Payload, _ = msg.MarshalBinary()

		// send handshake
		go n.send(pack)
	}
}

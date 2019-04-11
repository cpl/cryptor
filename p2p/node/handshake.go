package node

import (
	"cpl.li/go/cryptor/p2p/noise"
	"cpl.li/go/cryptor/p2p/packet"
	"cpl.li/go/cryptor/p2p/peer"
)

// Handshake will initialize the handshake protocol with the given peer.
func (n *Node) Handshake(p *peer.Peer) {
	// check if handshake is initialized
	if !p.HasHandshake {
		// create handshake
		var encISPub [48]byte
		p.Handshake, encISPub = noise.Initialize(
			n.identity.publicKey, p.PublicKey())
		tempPubKey := p.Handshake.PublicKey()

		// create packet
		pack := new(packet.Packet)
		pack.Type = packet.TypeHandshakeInitializer
		pack.Address = p.AddrUDP()
		pack.Payload = append(tempPubKey[:], encISPub[:]...)

		// send handshake
		go n.send(pack)
	}

}

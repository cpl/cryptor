package node

import (
	"cpl.li/go/cryptor/pkg/p2p/noise"
	"cpl.li/go/cryptor/pkg/p2p/packet"
	"cpl.li/go/cryptor/pkg/p2p/peer"
)

func (n *Node) handleTransport(p *peer.Peer, pack *packet.Packet) {
	return
}

func (n *Node) handleInitialize(peerID uint64, pack *packet.Packet) {
	// check payload size to match handshake initializer size
	if len(pack.Payload) < noise.SizeMessageInitializer {
		// drop packet if not
		return
	}

	// extract message from payload
	message := new(noise.MessageInitializer)
	if err := message.UnmarshalBinary(pack.Payload[:noise.SizeMessageInitializer]); err != nil {
		n.comm.err <- err
		return
	}

	// perform protocol response
	handshake, iSPub, rmsg, err := noise.Respond(
		message, n.identity.privateKey)

	// check if auth failed
	if err != nil {
		n.comm.err <- err
		return
	}

	// add peer
	p, err := n.PeerAdd(iSPub, pack.Address.String(), peerID)
	if err != nil {
		n.comm.err <- err
		return
	}

	// update peer handshake and transport keys
	p.Handshake = handshake

	// finalize handshake to get transport keys
	send, recv, err := p.Handshake.Finalize()
	if err != nil {
		n.comm.err <- err
		return
	}

	// set keys
	p.SetTransportKeys(send, recv)

	// send response to initializer
	response := new(packet.Packet)
	response.Type = packet.TypeHandshakeResponder
	response.Address = pack.Address
	response.Payload, _ = rmsg.MarshalBinary()

	n.comm.send <- response
}

func (n *Node) handleResponse(p *peer.Peer, pack *packet.Packet) {
	// check payload size to match handshake responder size
	if len(pack.Payload) < noise.SizeMessageResponder {
		// drop packet if not
		return
	}

	// extract message from payload
	message := new(noise.MessageResponder)
	if err := message.UnmarshalBinary(pack.Payload[:noise.SizeMessageResponder]); err != nil {
		n.comm.err <- err
		return
	}

	// validate handshake
	if err := p.Handshake.Receive(message, n.identity.privateKey); err != nil {
		// failed to validate packet
		n.comm.err <- err
		return
	}

	// finalize handshake to get transport keys
	send, recv, err := p.Handshake.Finalize()
	if err != nil {
		n.comm.err <- err
		return
	}

	// set keys
	p.SetTransportKeys(send, recv)
	return
}

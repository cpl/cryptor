package node

import (
	"encoding/binary"

	"cpl.li/go/cryptor/p2p/noise"
	"cpl.li/go/cryptor/p2p/packet"
)

// TODO break down recv into multiple internal functions
// - handleTransport
// - handleInitialize
// - handleResponse

func (n *Node) recv(pack *packet.Packet) {
	// extract message ID
	msgID := binary.LittleEndian.Uint64(pack.Payload)

	// perform peer lookup based on ID
	n.lookup.Lock()
	p, ok := n.lookup.table[msgID]
	if ok {
		p.Lock()
		defer p.Unlock()
	}
	n.lookup.Unlock()

	// peer exists and has complete handshake
	// -> transport message
	if ok && p.Handshake != nil && p.HasHandshake {
		// TODO Transport message handling
		return
	}

	// peer exists and has incomplete handshake
	// -> responder message
	if ok && p.Handshake != nil {
		// check payload size to match handshake responder size
		if len(pack.Payload) != noise.SizeMessageResponder {
			// drop packet if not
			return
		}

		// extract message from payload
		message := new(noise.MessageResponder)
		if err := message.UnmarshalBinary(pack.Payload); err != nil {
			n.comm.err <- err
			return
		}

		// validate handshake
		if err := p.Handshake.Receive(message, n.identity.privateKey); err != nil {
			// failed to validate packet
			n.comm.err <- err
			return
		}

		// update peer handshake and transport keys
		p.HasHandshake = true

		// finalize handshake to get transport keys
		_, _, err := p.Handshake.Finalize()
		if err != nil {
			n.comm.err <- err
			return
		}

		// TODO Set peer keys for encryption/decryption
		// p.SetTransportKeys(send, recv [ppk.KeySize]byte)
		return
	}

	// peer may or may not exist, expect initializer message

	// check payload size to match handshake initializer size
	if len(pack.Payload) != noise.SizeMessageInitializer {
		// drop packet if not
		return
	}

	// extract message from payload
	message := new(noise.MessageInitializer)
	if err := message.UnmarshalBinary(pack.Payload); err != nil {
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

	// if peer does not exist, add peer to lookup table
	if !ok {
		p, err = n.PeerAdd(iSPub, pack.Address.String())
		if err != nil {
			n.comm.err <- err
			return
		}
		p.Lock()
		defer p.Unlock()
	}

	// update peer handshake and transport keys
	p.Handshake = handshake
	p.HasHandshake = true

	// finalize handshake to get transport keys
	_, _, err = p.Handshake.Finalize()
	if err != nil {
		n.comm.err <- err
		return
	}

	// TODO Set peer keys for encryption/decryption
	// p.SetTransportKeys(send, recv [ppk.KeySize]byte)

	// send response to initializer
	response := new(packet.Packet)
	response.Type = packet.TypeHandshakeResponder
	response.Address = pack.Address
	response.Payload, _ = rmsg.MarshalBinary()

	go n.send(response)

	return
}

func (n *Node) send(pack *packet.Packet) {
	// check node is connected
	n.state.Lock()
	defer n.state.Unlock()
	if !n.state.isConnected {
		return
	}

	// ! DEBUG
	n.logger.Printf("sent packet (%d) to (%s)\n",
		len(pack.Payload), pack.Address.String())

	// send packet payload to its address
	_, err := n.net.conn.WriteToUDP(pack.Payload, pack.Address)
	if err != nil {
		n.comm.err <- err
	}
}

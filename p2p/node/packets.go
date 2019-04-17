package node

import (
	"fmt"

	"cpl.li/go/cryptor/p2p/noise"
	"cpl.li/go/cryptor/p2p/packet"
)

func (n *Node) recv(pack *packet.Packet) {
	// perform peer lookup
	n.lookup.Lock()
	p, ok := n.lookup.address[pack.Address.String()]
	n.lookup.Unlock()

	// if peer is not found, accept only handshake requests
	if !ok {
		// check payload size to match handshake initializer size
		if len(pack.Payload) != noise.HandshakeSizeInitializer {
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

		// if handshake is OK, add peer to lookup
		p, err := n.PeerAdd(iSPub, pack.Address.String())
		if err != nil {
			n.comm.err <- err
			return
		}

		// update peer handshake
		p.Lock()
		p.Handshake = handshake
		p.HasHandshake = true
		p.Unlock()

		// send response to initializer
		response := new(packet.Packet)
		response.Type = packet.TypeHandshakeResponder
		response.Address = pack.Address
		response.Payload, _ = rmsg.MarshalBinary()

		go n.send(response)

		return
	}
	n.lookup.Unlock()

	// peer is known

	// ! DEBUG
	fmt.Println(p)
}

func (n *Node) send(pack *packet.Packet) {
	// check node is connected
	n.state.Lock()
	defer n.state.Unlock()
	if !n.state.isConnected {
		return
	}

	// ! DEBUG
	n.logger.Println("sent packet")

	// send packet payload to its address
	_, err := n.net.conn.WriteToUDP(pack.Payload, pack.Address)
	if err != nil {
		n.comm.err <- err
	}
}

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
	if !ok {
		// if peer is not found, check payload to be handshake initializer
		if len(pack.Payload) != noise.HandshakeSizeInitializer {
			n.lookup.Unlock()
			// drop packet if not
			return
		}

		// handle handshake initializer
	}
	n.lookup.Unlock()

	// peer is known

	// ! DEBUG
	fmt.Println(p)
}

func (n *Node) send(pack *packet.Packet) {
	// check node is connected
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

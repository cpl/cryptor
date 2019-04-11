package node

import (
	"net"

	"cpl.li/go/cryptor/crypt"
	"cpl.li/go/cryptor/p2p"
	"cpl.li/go/cryptor/p2p/packet"
)

func (n *Node) run() {
	for {
		select {
		// pick up and display errors
		case err := <-n.comm.err:
			if err != nil {
				n.logger.Println("err", err)
			}
		// listen for exit signal
		case <-n.comm.exi:
			return
		}
	}
}

func (n *Node) connect() (err error) {
	// lock network state
	n.net.Lock()
	defer n.net.Unlock()

	// network bind using Node address
	n.net.conn, err = net.ListenUDP(p2p.Network, n.net.addr)

	return err
}

func (n *Node) disconnect() error {
	// ignore if node is not running
	if !n.state.isRunning {
		return nil
	}

	// ignore if node is not connected
	if !n.state.isConnected {
		return nil
	}

	// set node as disconnected
	n.state.isConnected = false

	// disconnect network bind
	if err := n.net.conn.Close(); err != nil {
		// on error, set node back as connected
		n.state.isConnected = true

		return err
	}

	n.logger.Println("disconnected")

	return nil
}

// listen checks the network for incoming connections, extracts the data
// and passes on valid packets only
func (n *Node) listen() {
	// incoming data buffer
	buffer := make([]byte, p2p.MaxPayloadSize)

	// zero buffer on disconnect
	defer crypt.ZeroBytes(buffer)

	for {
		// check if still connected
		if !n.state.isConnected {
			return
		}

		// read from network
		r, addr, err := n.net.conn.ReadFromUDP(buffer)
		if err != nil {
			// if disconnected return without error
			if !n.state.isConnected {
				return
			}

			// send error to node and retry
			n.comm.err <- err
			continue
		}

		n.logger.Println("receive packet")

		// check min size, drop packets if too small
		if r < p2p.MinPayloadSize {
			continue
		}

		// parse payload into packet
		pack := new(packet.Packet)
		pack.Address = addr
		pack.Payload = buffer[:r]

		// forward packet to handler
		go n.recv(pack)
	}
}

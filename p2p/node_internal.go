package p2p

import (
	"net"

	"cpl.li/go/cryptor/crypt"
)

func (n *Node) run() {
	for {
		select {
		// pick up and display errors
		case err := <-n.comm.err:
			if err != nil {
				n.logger.Println("err", err)
			}
		// receive exit signal
		case <-n.comm.exi:
			return
		}
	}
}

func (n *Node) connect() (err error) {
	// lock network state
	n.net.Lock()
	defer n.net.Unlock()

	// network bind on entire local address space, using random port
	n.net.conn, err = net.ListenUDP(Network, n.net.addr)

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

	// send disconnect signal, stop listening
	n.comm.dis <- nil

	n.logger.Println("disconnected")

	return nil
}

// forward handles sending and receiving data to and from the network
func (n *Node) forward() {
	for {
		select {
		// disconnect
		case <-n.comm.dis:
			return
		// receive
		case pack := <-n.net.recv:
			go n.recv(pack)
		// send
		case pack := <-n.net.send:
			go n.send(pack)
		}
	}
}

// listen checks the network for incoming connections, extracts the data
// and passes on valid packets only
func (n *Node) listen() {
	// incoming data buffer
	buffer := make([]byte, MaxUDPSize)

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

		// check min size, drop packets if too small
		if r < MsgMinSize {
			continue
		}

		var pack *Packet

		// check connection type
		if p, ok := n.lookup.address[addr.String()]; ok {
			// known peer

			// TODO Write checks for handshake state and handle package as such

			// ! DEBUG
			if p.handshake.status == 0 {

			}

		} else {
			// unknown

			// packet must be handshake initialization otherwise drop it
			if r != MsgSizeHandshakeI {
				continue
			}

			// ! DEBUG
			// convert received data to pack
			pack = new(Packet)
			pack.Address = addr
			pack.Payload = buffer[:r]
			pack.Type = packetTypeHandshakeI
		}

		// send parsed packet
		n.net.recv <- pack
	}
}

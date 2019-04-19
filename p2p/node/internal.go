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
				n.meta.errCount++
				n.logger.Println("err", err)
			}
		// listen for exit signal
		case <-n.comm.exi:
			return
		}
	}
}

// ! unsafe function, must be used only when the node state is defined & locked
func (n *Node) connect() (err error) {
	// lock network state
	n.net.Lock()
	defer n.net.Unlock()

	// network bind using Node address
	n.net.conn, err = net.ListenUDP(p2p.Network, n.net.addr)

	return err
}

// ! unsafe function, must be used only when the node state is defined & locked
func (n *Node) disconnect() error {
	// lock network state
	n.net.Lock()
	defer n.net.Unlock()

	// set node as disconnected
	n.state.isConnected = false

	// disconnect
	err := n.net.conn.Close()
	if err != nil {
		return err
	}

	// signal disconnect
	n.comm.dis <- nil

	return nil
}

// listen checks the network for incoming connections, extracts the data
// and passes on valid packets only
func (n *Node) listen() {
	// incoming data buffer
	buffer := make([]byte, p2p.MaxPayloadSize+1)

	// zero buffer on disconnect
	defer crypt.ZeroBytes(buffer)

	for {
		select {
		// on disconnect
		case <-n.comm.dis:
			n.logger.Println("disconnected")
			return
		default:
			// read from network
			r, addr, err := n.net.conn.ReadFromUDP(buffer)
			if err != nil {
				// return if not connected anymore
				if !n.state.isConnected {
					continue
				} else {
					// attempt safe disconnect on failed connection
					n.Disconnect()
					continue
				}
			}

			// ! DEBUG
			n.logger.Printf("receive packet from %s\n", addr.String())

			// check min size, drop packets if too small
			if r < p2p.MinPayloadSize || r > p2p.MaxPayloadSize {
				// ! DEBUG
				n.logger.Printf("drop packet, size %d outside bounds\n", r)
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
}

package node

import (
	"encoding/binary"
	"net"
	"sync/atomic"

	"cpl.li/go/cryptor/pkg/p2p"
	"cpl.li/go/cryptor/pkg/p2p/noise"
	"cpl.li/go/cryptor/pkg/p2p/packet"
)

func (n *Node) run() {
	n.state.starting.Done()
	defer n.state.stopping.Done()
	for {
		select {
		// pick up and display errors
		case err := <-n.comm.err:
			if err != nil {
				atomic.AddUint32(&n.meta.errCount, 1)
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
	for {
		// incoming data buffer
		buffer := make([]byte, p2p.MaxPayloadSize+1)

		select {
		// on disconnect
		case <-n.comm.dis:
			n.logger.Println("disconnected")
			return
		default:
			// read from network
			r, addr, err := n.net.conn.ReadFromUDP(buffer)
			if err != nil {
				if !n.state.isConnected {
					// return if not connected anymore
					// as something might be waiting to send on n.comm.dis
					continue
				} else {
					// attempt safe disconnect on failed connection
					n.comm.err <- n.Disconnect()
					continue
				}
			}

			// check min size, drop packets if too small
			if r < p2p.MinPayloadSize || r > p2p.MaxPayloadSize {
				continue
			}

			// parse payload into packet
			pack := new(packet.Packet)
			pack.Address = addr
			pack.Payload = buffer[:r]

			// forward packet pointer
			// this will block the listener if the receiving channel is full,
			// this ensures we don't get overloaded and should be fine if it
			// happens to lose a few packets + some packets will get cached by
			// the OS/network
			n.comm.recv <- pack
		}
	}
}

func (n *Node) recv() {
	n.state.starting.Done()
	defer n.state.stopping.Done()
	for {
		select {
		case <-n.comm.exi:
			return
		case pack := <-n.comm.recv:
			// extract message ID
			peerID := binary.LittleEndian.Uint64(pack.Payload)

			// perform peer lookup based on ID
			n.lookup.RLock()
			p, ok := n.lookup.table[peerID]
			if ok {
				p.Lock()
			}
			n.lookup.RUnlock()

			// peer exists and has complete handshake
			// -> transport message
			if ok && p.Handshake != nil && p.Handshake.State() == noise.StateSuccessful {
				n.handleTransport(p, pack)
				p.Unlock()
				continue
			}

			// peer exists and has incomplete handshake (waiting for response)
			// -> responder message
			if ok && p.Handshake != nil && p.Handshake.State() == noise.StateInitialized {
				n.handleResponse(p, pack)
				p.Unlock()
				continue
			}

			// peer does not exist
			// -> initializer message
			if !ok {
				n.handleInitialize(peerID, pack)
				continue
			}

			// drop packet
		}
	}
}

func (n *Node) send() {
	n.state.starting.Done()
	defer n.state.stopping.Done()
	for {
		select {
		case <-n.comm.exi:
			return
		case pack := <-n.comm.send:
			// send packet payload to its address
			n.net.RLock()
			_, err := n.net.conn.WriteToUDP(pack.Payload, pack.Address)
			if err != nil {
				n.comm.err <- err
			}
			n.net.RUnlock()
		}
	}
}

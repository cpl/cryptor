package node

import (
	"cpl.li/go/cryptor/p2p/packet"
)

// TODO Finish implementing these

func (n *Node) recv(pack *packet.Packet) {

}

func (n *Node) send(pack *packet.Packet) {
	_, err := n.net.conn.WriteToUDP(pack.Payload, pack.Address)
	if err != nil {
		n.comm.err <- err
	}
}

package p2p

// TODO Finish implementing these

func (n *Node) recv(pack *Packet) {
	switch len(pack.Payload) {
	case MsgSizeHandshakeI:
	case MsgSizeHandshakeR:
	}
}

func (n *Node) send(pack *Packet) {
	_, err := n.net.conn.WriteToUDP(pack.Payload, pack.Address)
	if err != nil {
		n.comm.err <- err
	}
}

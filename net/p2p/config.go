package p2p

import "net"

// NodeConfig is a static structure defined at node creation.
type NodeConfig struct {
	IP           net.IP  // IP Address of the node
	TCP, UDP     int     // Ports used for node controll (TCP) and P2P (UDP)
	TrustedPeers []*Peer // Maps of trusted peers
}

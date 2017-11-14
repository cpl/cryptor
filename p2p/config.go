package p2p

import "net"

// NodeConfig is a static structure defined at node creation.
type NodeConfig struct {
	IP       net.IP // IP Address of the node
	TCP, UDP int    // Ports used for P2P and chunk sharing
}

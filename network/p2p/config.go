package p2p

import "net"

// NodeConfig ...
type NodeConfig struct {
	IP       net.IP // IP Adress of the node
	TCP, UDP int    // Ports used for P2P and file sharing
}

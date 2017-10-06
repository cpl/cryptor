package p2p

import "net"

// NodeConfig ...
type NodeConfig struct {
	IP       net.IP
	TCP, UDP int
}

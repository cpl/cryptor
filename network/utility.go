package network

import "net"

// IPPToUDP or IP:PORT TO UDP, takes an ip and port and converts them to an
// UDPAddr object.
func IPPToUDP(ip string, port int) *net.UDPAddr {
	return &net.UDPAddr{
		IP:   net.ParseIP(ip),
		Port: port,
	}
}

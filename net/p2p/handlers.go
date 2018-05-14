package p2p

import (
	"log"
	"net"
)

// !DEBUG
func handleConn(conn net.PacketConn) {
	buffer := make([]byte, 1024)
	for {
		n, addr, err := conn.ReadFrom(buffer)
		if err != nil {
			log.Fatalln(err)
		}
		log.Println("msg from", addr.String())
		log.Println(string(buffer[:n-1]))
	}
}

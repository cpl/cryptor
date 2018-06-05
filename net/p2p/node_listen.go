package p2p

import (
	"github.com/thee-engineer/cryptor/crypt"
)

func (n *Node) listen() {
	buffer := make([]byte, 1024)
	defer crypt.ZeroBytes(buffer)
	for {
		select {
		case <-n.disconnect:
			n.logChan <- "disconnecting ..."
			return
		default:
			r, addr, err := n.udpConn.ReadFrom(buffer)
			if err != nil {
				n.errChan <- err
				continue
			}
			n.logChan <- "new connection from " + addr.String()
			n.incoming <- buffer[:r-1]
		}
	}
}

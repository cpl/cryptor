package p2p

import "log"

func (n *Node) run() {
	for {
		select {
		case err := <-n.errChan:
			log.Println("node err:", err)
			n.status()
		case msg := <-n.logChan:
			log.Println("node log:", msg)
		case operation := <-n.peerOp:
			operation(n.peers)
			n.peerOpDone <- nil
		case <-n.quit:
			log.Println("node log: stopping ...")
			return
		case msg := <-n.incoming:
			smsg := string(msg)
			log.Println("msg:", smsg) // ! DEBUG
			if smsg == "stop" {
				go n.Stop()
			}
			if smsg == "dc" {
				go n.Disconnect()
			}
		}
	}
}

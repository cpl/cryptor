package p2p

import "log"

func checkRunning(n *Node) bool {
	if !n.isRunning {
		log.Println("node err: not running")
	}
	return n.isRunning
}

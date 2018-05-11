package p2p

import "log"

func checkRunning(n *Node) bool {
	n.lock.Lock()
	defer n.lock.Unlock()
	if !n.isRunning {
		log.Println("node err: not running")
	}
	return n.isRunning
}

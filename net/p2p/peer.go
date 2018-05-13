package p2p

// Peer ...
type Peer struct {
	PublicKey string
	Address   string
}

type peerMap map[string]*Peer
type peerFunc func(peerMap)

// CountPeers ...
func (n *Node) CountPeers() int {
	if !checkRunning(n) {
		return 0
	}

	var count int

	select {
	case n.peerOp <- func(peerList peerMap) {
		count = len(peerList)
	}:
		<-n.peerOpDone
	}

	return count
}

// Peers ...
func (n *Node) Peers() []*Peer {
	if !checkRunning(n) {
		return nil
	}

	var peerList []*Peer

	select {
	case n.peerOp <- func(peers peerMap) {
		for _, peer := range peers {
			peerList = append(peerList, peer)
		}
	}:
		<-n.peerOpDone
	}

	return peerList
}

// AddPeer ...
func (n *Node) AddPeer(peer *Peer) {
	if !checkRunning(n) {
		return
	}

	select {
	case n.peerOp <- func(peers peerMap) {
		peers[peer.PublicKey] = peer
	}:
		<-n.peerOpDone
	}
	n.logChan <- "add peer: " + peer.PublicKey
}

// DelPeer ...
func (n *Node) DelPeer(peer *Peer) {
	if !checkRunning(n) {
		return
	}

	select {
	case n.peerOp <- func(peers peerMap) {
		delete(peers, peer.PublicKey)
	}:
		<-n.peerOpDone
	}
	n.logChan <- "del peer: " + peer.PublicKey
}

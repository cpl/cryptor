package p2p_test

import (
	"testing"
	"time"

	"github.com/thee-engineer/cryptor/net/p2p"
)

func TestNodeStop(t *testing.T) {
	qc := make(chan struct{})
	node := p2p.NewNode("127.0.0.1", 2002, 9002, qc)
	go node.Stop() // Stop before starting

	go node.Start()
	time.Sleep(time.Second)

	go node.Stop()
}

func TestNodeStart(t *testing.T) {
	qc := make(chan struct{})

	node := p2p.NewNode("127.0.0.1", 2001, 9001, qc)

	go node.Start()

	time.Sleep(time.Second)

	qc <- *new(struct{})
}

func TestNodeConnection(t *testing.T) {
	qc := make(chan struct{})

	n0 := p2p.NewNode("127.0.0.1", 2010, 9010, qc)
	n1 := p2p.NewNode("127.0.0.1", 2011, 9011, qc)

	n0.Send(p2p.NewUDPPacket([]byte("hello world"), n0.UDPAddr()))
	n0.Send(p2p.NewUDPPacket([]byte("hello world node1"), n1.UDPAddr()))

	time.Sleep(time.Second)

	go n0.Start()
	go n1.Start()

	time.Sleep(time.Second)

	go n0.Listen()
	go n1.Listen()

	time.Sleep(time.Second)

	n0.Send(p2p.NewUDPPacket([]byte("hello world node1, 2"), n1.UDPAddr()))
	n1.Send(p2p.NewUDPPacket([]byte("hi node0, sup?"), n0.UDPAddr()))

	time.Sleep(time.Second)

	qc <- *new(struct{})
}

func TestSamePort(t *testing.T) {
	qc := make(chan struct{})

	n0 := p2p.NewNode("127.0.0.1", 2020, 9020, qc)
	n1 := p2p.NewNode("127.0.0.1", 2020, 9020, qc)

	time.Sleep(time.Second)

	go n0.Start()
	go n1.Start()

	time.Sleep(time.Second)

	go n0.Listen()

	time.Sleep(time.Second)

	go n1.Listen()

	time.Sleep(time.Second)

	qc <- *new(struct{})
}

func TestNodeMutex(t *testing.T) {
	qc := make(chan struct{})
	n0 := p2p.NewNode("127.0.0.1", 2030, 9030, qc)

	n0.Stop()
	go n0.Stop()
	n0.Stop()

	n0.Disconnect()
	go n0.Disconnect()
	n0.Disconnect()

	time.Sleep(time.Second)

	go n0.Start()
	go n0.Start()

	time.Sleep(time.Second)

	go n0.Listen()
	go n0.Listen()

	time.Sleep(time.Second)

	go n0.Listen()

	time.Sleep(time.Second)

	go n0.Disconnect()
	go n0.Stop()

	time.Sleep(time.Second)

	qc <- *new(struct{})
}

func TestNodeDisconnect(t *testing.T) {
	qc := make(chan struct{})
	n0 := p2p.NewNode("127.0.0.1", 2040, 9040, qc)

	go n0.Start()

	time.Sleep(time.Second)

	go n0.Listen()

	time.Sleep(time.Second)

	go n0.Disconnect()

	time.Sleep(time.Second)

	go n0.Stop()
}

func TestNodeDisconnectStrange(t *testing.T) {
	qc := make(chan struct{})
	n0 := p2p.NewNode("127.0.0.1", 2050, 9050, qc)

	go n0.Listen()
	go n0.Start()
	go n0.Listen()
	go n0.Stop()
	go n0.Disconnect()

	time.Sleep(2 * time.Second)
}

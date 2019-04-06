package proto

import (
	"bytes"
	"encoding/gob"
)

// Packet contains the necessary info for sending a message over the network.
// TODO Expand Packet structure to cover more needed information
type Packet struct {
	MsgType byte
	MsgData []byte
}

// Bytes returns the packet as a single byte array.
// TODO Improve binary marshaling of packets
func (p *Packet) Bytes() []byte {
	var buf bytes.Buffer
	defer buf.Reset()

	enc := gob.NewEncoder(&buf)
	enc.Encode(p)
	return buf.Bytes()
}

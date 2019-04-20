/*
Package p2p implements the networking aspect of the Cryptor Network. Allowing
nodes and peers to connect.
*/
package p2p // import "cpl.li/go/cryptor/p2p"

// Network tells golang net package which network type to use.
// In this case we use "udp", allowing "udp4" and "udp6".
const Network = "udp"

// ! DEBUG
const maxUDPSize = 1024

const (
	// MaxPayloadSize is the maximum size of a data payload.
	MaxPayloadSize = maxUDPSize

	// MinPayloadSize is the minimum size of a data payload.
	MinPayloadSize = 48
)

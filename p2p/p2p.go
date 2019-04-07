/*
Package p2p implements the networking aspect of the Cryptor Network. Allowing
nodes and peers to connect.
*/
package p2p // import "cpl.li/go/cryptor/p2p"

// Network tells golang net package which network type to use.
// In this case we use "udp", allowing "udp4" and "udp6".
const Network = "udp"

// MaxUDPSize is the maximum number of bytes a UDP packet may safely contain.
// TODO Readup on UDP packet safe size, 1024 is for testing only
const MaxUDPSize = 1024

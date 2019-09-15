package peer_test

import (
	"testing"

	"cpl.li/go/cryptor/pkg/crypt/ppk"
	"cpl.li/go/cryptor/pkg/p2p/peer"
	"github.com/stretchr/testify/assert"
)

var zeroKey ppk.PublicKey

func testSetAddr(t *testing.T, p *peer.Peer, addr, expected string) {
	assert.Nil(t, p.SetAddr(addr))
	assert.Equal(t, p.Addr(), expected, "unexpected address")
}

func TestPeerSetAddr(t *testing.T) {
	t.Parallel()

	p := peer.NewPeer(zeroKey, "")

	// valid tests
	testSetAddr(t, p, "", "<nil>")
	testSetAddr(t, p, "127.0.0.1:", "127.0.0.1:0")
	testSetAddr(t, p, "192.168.1.1:", "192.168.1.1:0")
	testSetAddr(t, p, ":", ":0")
	testSetAddr(t, p, "", "<nil>")
	testSetAddr(t, p, ":1234", ":1234")

	// invalid
	assert.NotNil(t, p.SetAddr("1.1.1.1"), "set invalid address, no port")
	assert.NotNil(t, p.SetAddr("nosuchhost:"), "set invalid address, host")
	assert.NotNil(t, p.SetAddr("1.1.1.1:-1"), "set invalid address, invalid port")

	// check unchanged valid address
	assert.Equal(t, p.Addr(), ":1234", "unexpected address")
}

func TestSetTransportKeys(t *testing.T) {
	t.Parallel()

	var zeroPk ppk.PublicKey
	var key1, key2 [ppk.KeySize]byte
	p := peer.NewPeer(zeroPk, "")
	p.SetTransportKeys(key1, key2)
}

func TestNewPeerNoAddr(t *testing.T) {
	t.Parallel()

	p := peer.NewPeer(zeroKey, "")

	// check for ID 0
	// peer ID is assigned during handshake and not creation
	assert.Equal(t, p.ID, uint64(0), "invalid peer ID")

	// validate key
	if !p.PublicKey().Equals(zeroKey) {
		t.Fatal("public key does not match")
	}

	// validate default address
	if addr := p.AddrUDP(); addr != nil {
		t.Fatal("got non-nil udp address", addr)
	}
	assert.Equal(t, p.Addr(), "<nil>", "invalid peer address")
}

func TestNewPeerWithAddr(t *testing.T) {
	t.Parallel()

	p := peer.NewPeer(zeroKey, "127.0.0.1:8080")

	// validate key
	if !p.PublicKey().Equals(zeroKey) {
		t.Fatal("public key does not match")
	}

	// validate address
	if p.AddrUDP() == nil {
		t.Fatal("got nil udp address")
	}

	assert.Equal(t, p.Addr(), "127.0.0.1:8080", "unexpected peer address")
}

func TestPeerAddressParsing(t *testing.T) {
	t.Parallel()

	addresses := []struct {
		addr string
		expe string
	}{
		{"127.0.0.1:8000", "127.0.0.1:8000"},
		{"127.0.0.1", "<nil>"},
		{"192.168.0.1:15000", "192.168.0.1:15000"},
		{"192.100.1.0:1", "192.100.1.0:1"},
		{"192.168.2.2", "<nil>"},
		{"[::1]:8080", "[::1]:8080"},
		{"[::1]", "<nil>"},
		{":8000", ":8000"},
		{"127.0.0.1:50000000", "<nil>"},
		{"[::1]:65535", "[::1]:65535"},
		{"nil", "<nil>"},
		{"no such host", "<nil>"},
		{"nil:1000", "<nil>"},
		{"0.0.0.0:1000", "0.0.0.0:1000"},
		{"[::]:22", "[::]:22"},
	}

	for _, addr := range addresses {
		p := peer.NewPeer(zeroKey, addr.addr)
		assert.Equal(t, p.Addr(), addr.expe, "unexpected peer address")
	}
}

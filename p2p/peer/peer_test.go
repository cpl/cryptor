package peer_test

import (
	"testing"

	"cpl.li/go/cryptor/crypt/ppk"
	"cpl.li/go/cryptor/p2p/peer"
	"cpl.li/go/cryptor/tests"
)

var zeroKey ppk.PublicKey

func testSetAddr(t *testing.T, p *peer.Peer, addr, expected string) {
	tests.AssertNil(t, p.SetAddr(addr))
	if p.Addr() != expected {
		t.Fatalf("expected %s, got %s\n", expected, p.Addr())
	}
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
	tests.AssertNotNil(t, p.SetAddr("1.1.1.1"), "set invalid address, no port")
	tests.AssertNotNil(t, p.SetAddr("nosuchhost:"), "set invalid address, host")
	tests.AssertNotNil(t, p.SetAddr("1.1.1.1:-1"), "set invalid address, invalid port")

	// check unchanged valid address
	if addr := p.Addr(); addr != ":1234" {
		t.Fatalf("expected :1234, got %s\n", addr)
	}
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
	if p.ID != 0 {
		t.Fatal("peer ID is not 0")
	}

	// validate key
	if !p.PublicKey().Equals(zeroKey) {
		t.Fatal("public key does not match")
	}

	// validate default address
	if p.AddrUDP() != nil {
		t.Fatal("got non-nil udp address")
	}
	if p.Addr() != "<nil>" {
		t.Fatalf("got non-nil address as string, %s\n", p.Addr())
	}
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
	if p.Addr() != "127.0.0.1:8080" {
		t.Fatalf("got an invalid address as string, %s\n", p.Addr())
	}
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
		if p.Addr() != addr.expe {
			t.Errorf("failed to match %s with %s\n", addr, p.Addr())
		}
	}
}

package peer_test

import (
	"testing"

	"cpl.li/go/cryptor/crypt/ppk"
	"cpl.li/go/cryptor/p2p/peer"
)

var zeroKey ppk.PublicKey

func TestNewPeerNoAddr(t *testing.T) {
	p := peer.NewPeer(zeroKey, "")

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

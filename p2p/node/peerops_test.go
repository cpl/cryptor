package node_test

import (
	"sync"
	"testing"

	"cpl.li/go/cryptor/crypt"
	"cpl.li/go/cryptor/crypt/ppk"
	"cpl.li/go/cryptor/p2p/node"
	"cpl.li/go/cryptor/tests"
)

// generates a random public key for testing only
func newRandomPublicKey() ppk.PublicKey {
	var public ppk.PublicKey
	copy(public[:], crypt.RandomBytes(ppk.KeySize))
	return public
}

// can be called as a goroutine
func parallelAdd(t *testing.T, n *node.Node, wg *sync.WaitGroup) {
	if _, err := n.PeerAdd(newRandomPublicKey(), ""); err != nil {
		t.Error(err)
	}
	wg.Done()
}

func TestPeerAddAndDel(t *testing.T) {
	// create node
	n := node.NewNode("test", zeroKey)

	// attempt to add peer while not running
	if _, err := n.PeerAdd(newRandomPublicKey(), ""); err == nil {
		t.Fatal("added peer while node is not running")
	}

	// attempt to del peer while not running
	if err := n.PeerDel(newRandomPublicKey()); err == nil {
		t.Fatal("deleted peer while node is not running")
	}

	// attempt peer count while not running
	if count := n.PeerCount(); count != -1 {
		t.Fatalf("unexpected value from PeerCount, wanted -1, got %d\n", count)
	}

	// start node, stop at the end
	n.Start()
	defer n.Stop()

	// add 100 peers
	for i := 0; i < 100; i++ {
		if _, err := n.PeerAdd(newRandomPublicKey(), ""); err != nil {
			t.Fatal(err)
		}
	}

	// count peers
	if count := n.PeerCount(); count != 100 {
		t.Fatalf("expected 100 peers, got %d\n", count)
	}

	// try to add another peer twice with the same key
	key := newRandomPublicKey()
	if _, err := n.PeerAdd(key, ""); err != nil {
		t.Fatal(err)
	}
	if _, err := n.PeerAdd(key, ""); err == nil {
		t.Fatal("added duplicate key peer")
	}

	// count peers
	if count := n.PeerCount(); count != 101 {
		t.Fatalf("expected 101 peers, got %d\n", count)
	}

	// add 100 new peers in parallel to test the mutex
	// this should not happen in parallel, as the lookup map
	// requires a lock
	wg := new(sync.WaitGroup)
	wg.Add(100)
	for i := 0; i < 100; i++ {
		go parallelAdd(t, n, wg)
	}
	wg.Wait()

	// count peers
	if count := n.PeerCount(); count != 201 {
		t.Fatalf("expected 201 peers, got %d\n", count)
	}

	// create 100 keys, add them as peers
	publicKeys := make([]ppk.PublicKey, 100)
	for idx := range publicKeys {
		publicKeys[idx] = newRandomPublicKey()
		if _, err := n.PeerAdd(publicKeys[idx], ""); err != nil {
			t.Fatal(err)
		}
	}

	// count peers
	if count := n.PeerCount(); count != 301 {
		t.Fatalf("expected 301 peers, got %d\n", count)
	}

	// remove the 100 peers
	for _, key := range publicKeys {
		if err := n.PeerDel(key); err != nil {
			t.Fatal(err)
		}
	}

	// attempt to remote random peer
	if err := n.PeerDel(newRandomPublicKey()); err == nil {
		t.Error("removed random peer")
	}

	// count peers
	if count := n.PeerCount(); count != 201 {
		t.Fatalf("expected 201 peers, got %d\n", count)
	}
}

func TestPeerList(t *testing.T) {
	// create node
	n := node.NewNode("test", zeroKey)

	// attempt to list peers while not running
	if err := n.PeerList(); err == nil {
		t.Fatal("listed peers while node not running")
	}

	// start node, stop at the end
	n.Start()
	defer n.Stop()

	// add 8 peers
	for i := 0; i < 8; i++ {
		if _, err := n.PeerAdd(newRandomPublicKey(), ""); err != nil {
			t.Fatal(err)
		}
	}

	// count peers
	if count := n.PeerCount(); count != 8 {
		t.Fatalf("expected 8 peers, got %d\n", count)
	}

	// peer list
	if err := n.PeerList(); err != nil {
		t.Fatal(err)
	}
}

func TestPeerGet(t *testing.T) {
	// create and start node
	n := node.NewNode("test", zeroKey)
	tests.AssertNil(t, n.Start())

	// generate keys and peers
	keys := make([]ppk.PublicKey, 8)
	for i := 0; i < 8; i++ {
		keys[i] = newRandomPublicKey()
		if _, err := n.PeerAdd(keys[i], ""); err != nil {
			t.Fatal(err)
		}
	}

	// check count
	if count := n.PeerCount(); count != 8 {
		t.Fatalf("expected 8 peers, got %d\n", count)
	}

	// invalid gets
	for i := 0; i < 8; i++ {
		if p := n.PeerGet(newRandomPublicKey()); p != nil {
			t.Errorf("got non-nil peer, expected non peer, %d\n", i)
		}
	}

	// valid gets
	for i := 0; i < 8; i++ {
		if p := n.PeerGet(keys[i]); p == nil {
			t.Errorf("got nil peer, expected non-nil peer, %d\n", i)
		}
	}

	// stop
	tests.AssertNil(t, n.Stop())
}

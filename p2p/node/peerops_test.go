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
	if _, err := n.PeerAdd(newRandomPublicKey(), "", crypt.RandomUint64()); err != nil {
		t.Error(err)
	}
	wg.Done()
}

func TestPeeropsInvalid(t *testing.T) {
	t.Parallel()

	// create node
	n := node.NewNode("test", zeroKey)

	// attempt add with ID 0
	if _, err := n.PeerAdd(newRandomPublicKey(), "", 0); err == nil {
		t.Fatal("added peer with ID 0")
	}

	// attempt to remove random peer (no peers)
	if err := n.PeerDel(newRandomPublicKey()); err == nil {
		t.Error("removed random peer")
	}

	// attempt to remove random peer (no peers)
	if err := n.PeerDelID(crypt.RandomUint64()); err == nil {
		t.Error("removed random peer")
	}

	// try to add another peer twice with the same id
	id := crypt.RandomUint64()
	if _, err := n.PeerAdd(newRandomPublicKey(), "", id); err != nil {
		t.Fatal(err)
	}
	if _, err := n.PeerAdd(newRandomPublicKey(), "", id); err == nil {
		t.Fatal("added duplicate id peer")
	}

	// try to add another peer twice with the same key
	key := newRandomPublicKey()
	if _, err := n.PeerAdd(key, "", crypt.RandomUint64()); err != nil {
		t.Fatal(err)
	}
	if _, err := n.PeerAdd(key, "", crypt.RandomUint64()); err == nil {
		t.Fatal("added duplicate id peer")
	}

	// attempt to remove random peer
	if err := n.PeerDel(newRandomPublicKey()); err == nil {
		t.Error("removed random peer")
	}

	// attempt to remove random peer
	if err := n.PeerDelID(crypt.RandomUint64()); err == nil {
		t.Error("removed random peer")
	}
}

func TestPeerAddAndDel(t *testing.T) {
	t.Parallel()

	// create node
	n := node.NewNode("test", zeroKey)

	// add 100 peers
	for i := 0; i < 100; i++ {
		if _, err := n.PeerAdd(newRandomPublicKey(), "", crypt.RandomUint64()); err != nil {
			t.Fatal(err)
		}
	}

	// count peers
	if count := n.PeerCount(); count != 100 {
		t.Fatalf("expected 100 peers, got %d\n", count)
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
	if count := n.PeerCount(); count != 200 {
		t.Fatalf("expected 200 peers, got %d\n", count)
	}

	// create 100 keys, add them as peers
	publicKeys := make([]ppk.PublicKey, 100)
	for idx := range publicKeys {
		publicKeys[idx] = newRandomPublicKey()
		if _, err := n.PeerAdd(publicKeys[idx], "", crypt.RandomUint64()); err != nil {
			t.Fatal(err)
		}
	}

	// count peers
	if count := n.PeerCount(); count != 300 {
		t.Fatalf("expected 300 peers, got %d\n", count)
	}

	// remove the 100 peers
	for _, key := range publicKeys {
		if err := n.PeerDel(key); err != nil {
			t.Fatal(err)
		}
	}

	// count peers
	if count := n.PeerCount(); count != 200 {
		t.Fatalf("expected 200 peers, got %d\n", count)
	}
}

func TestPeerList(t *testing.T) {
	t.Parallel()

	// create node
	n := node.NewNode("test", zeroKey)

	// add 8 peers
	for i := 0; i < 8; i++ {
		if _, err := n.PeerAdd(newRandomPublicKey(), "", crypt.RandomUint64()); err != nil {
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
	t.Parallel()

	// create and start node
	n := node.NewNode("test", zeroKey)

	// generate keys and peers
	keys := make([]ppk.PublicKey, 8)
	ids := make([]uint64, 8)
	for i := 0; i < 8; i++ {
		keys[i] = newRandomPublicKey()
		ids[i] = crypt.RandomUint64()
		if _, err := n.PeerAdd(keys[i], "", ids[i]); err != nil {
			t.Fatal(err)
		}
	}

	// check count
	if count := n.PeerCount(); count != 8 {
		t.Fatalf("expected 8 peers, got %d\n", count)
	}

	// invalid gets
	for i := 0; i < 8; i++ {
		if p := n.PeerGet(0, newRandomPublicKey()); p != nil {
			t.Errorf("got non-nil peer, expected non peer, %d\n", i)
		}
	}
	for i := 0; i < 8; i++ {
		if p := n.PeerGet(crypt.RandomUint64(), newRandomPublicKey()); p != nil {
			t.Errorf("got non-nil peer, expected non peer, %d\n", i)
		}
	}

	// valid gets
	for i := 0; i < 8; i++ {
		if p := n.PeerGet(0, keys[i]); p == nil {
			t.Errorf("got nil peer, expected non-nil peer, %d\n", i)
		}
	}
	for i := 0; i < 8; i++ {
		if p := n.PeerGet(ids[i], keys[i]); p == nil {
			t.Errorf("got nil peer, expected non-nil peer, %d\n", i)
		}
	}
}

func TestPeerDel(t *testing.T) {
	t.Parallel()

	// create and start node
	n := node.NewNode("test", zeroKey)

	// add peer
	key := newRandomPublicKey()
	id := crypt.RandomUint64()
	p, err := n.PeerAdd(key, "", id)
	tests.AssertNil(t, err)

	// check id and public key
	if p.ID != id {
		t.Fatalf("invalid id, expected %d, got %d\n", id, p.ID)
	}
	if !p.PublicKey().Equals(key) {
		t.Fatalf("mismatch public keys")
	}

	// delete with public key
	tests.AssertNil(t, n.PeerDel(key))

	// add again
	p, err = n.PeerAdd(key, "", id)
	tests.AssertNil(t, err)

	// delete with id
	tests.AssertNil(t, n.PeerDelID(id))
}

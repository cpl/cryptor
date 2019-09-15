package node_test

import (
	"sync"
	"testing"

	"cpl.li/go/cryptor/pkg/crypt"
	"cpl.li/go/cryptor/pkg/crypt/ppk"
	"cpl.li/go/cryptor/pkg/p2p/node"

	"github.com/stretchr/testify/assert"
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

func assertPeerCount(t *testing.T, n *node.Node, expected int) {
	assert.Equal(t, n.PeerCount(), expected, "unexpected peer count")
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
	assertPeerCount(t, n, 100)

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
	assertPeerCount(t, n, 200)

	// create 100 keys, add them as peers
	publicKeys := make([]ppk.PublicKey, 100)
	for idx := range publicKeys {
		publicKeys[idx] = newRandomPublicKey()
		if _, err := n.PeerAdd(publicKeys[idx], "", crypt.RandomUint64()); err != nil {
			t.Fatal(err)
		}
	}

	// count peers
	assertPeerCount(t, n, 300)

	// remove the 100 peers
	for _, key := range publicKeys {
		assert.Nil(t, n.PeerDel(key))
	}

	// count peers
	assertPeerCount(t, n, 200)
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
	assertPeerCount(t, n, 8)

	// peer list
	assert.Nil(t, n.PeerList())
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
	assertPeerCount(t, n, 8)

	// invalid gets
	for i := 0; i < 8; i++ {
		if p := n.PeerGet(0, newRandomPublicKey()); p != nil {
			t.Errorf("got non-nil peer, expected nil, %d\n", i)
		}
	}
	for i := 0; i < 8; i++ {
		if p := n.PeerGet(crypt.RandomUint64(), newRandomPublicKey()); p != nil {
			t.Errorf("got non-nil peer, expected nil, %d\n", i)
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
	assert.Nil(t, err)

	// check id and public key
	assert.Equal(t, p.ID, id, "invalid peer id")
	if !p.PublicKey().Equals(key) {
		t.Fatalf("mismatch public keys")
	}

	// delete with public key
	assert.Nil(t, n.PeerDel(key))

	// add again
	p, err = n.PeerAdd(key, "", id)
	assert.Nil(t, err)

	// delete with id
	assert.Nil(t, n.PeerDelID(id))
}

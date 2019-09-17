package hkdf_test

import (
	"crypto/hmac"
	"encoding/hex"
	"testing"

	"golang.org/x/crypto/blake2s"
	"github.com/stretchr/testify/assert"

	"cpl.li/go/cryptor/pkg/crypt"
	"cpl.li/go/cryptor/pkg/crypt/hkdf"
)

func TestHMAC(t *testing.T) {
	t.Parallel()

	// test data
	data := []byte("We attack at dawn.")
	invalidData := []byte("We attack at night.")
	invalidKey := crypt.RandomBytes(32)

	// test key
	hexkey := "93b5fb7506799aaeb93b91cff9fbbcf798d26df95c6d579fce45af881a0b471b"
	key, _ := hex.DecodeString(hexkey)

	// create sum arrays
	var sum0 [blake2s.Size]byte
	var sum1 [blake2s.Size]byte

	// compute sums
	hkdf.HMAC(&sum0, key, data)
	hkdf.HMAC(&sum1, key, data)
	assertSums(t, sum0, sum1, true)

	// check with invalid key but valid data
	hkdf.HMAC(&sum1, invalidKey, data)
	assertSums(t, sum0, sum1, false)

	// check with valid key but invalid data
	hkdf.HMAC(&sum1, key, invalidData)
	assertSums(t, sum0, sum1, false)

	// check with both invalid key and data
	hkdf.HMAC(&sum1, invalidKey, invalidData)
	assertSums(t, sum0, sum1, false)

	// recompute valid sum and check to be equal again
	hkdf.HMAC(&sum1, key, data)
	assertSums(t, sum0, sum1, true)
}

func assertSums(t *testing.T, sum0, sum1 [blake2s.Size]byte, expect bool) {
	if hmac.Equal(sum0[:], sum1[:]) != expect {
		t.Fatalf("invalid hmac sums")
	}
}

func TestHMACMultiple(t *testing.T) {
	// test data
	data0 := []byte("We attack at dawn.")
	data1 := []byte("You attack from the East.")

	// test key
	hexkey := "93b5fb7506799aaeb93b91cff9fbbcf798d26df95c6d579fce45af881a0b471b"
	key, _ := hex.DecodeString(hexkey)

	// create sum arrays
	var sum0 [blake2s.Size]byte
	var sum1 [blake2s.Size]byte

	// generate data in both orders
	hkdf.HMAC(&sum0, key, data0, data1)
	hkdf.HMAC(&sum1, key, data1, data0)

	// must not be equal
	assertSums(t, sum0, sum1, false)
}

type kdfTest struct {
	key  string
	data string

	// expected results
	t0 string
	t1 string
	t2 string
}

// test data taken from WireGuard - Noise protocol golang implementation
// https://git.zx2c4.com/wireguard-go/tree/device/kdf_test.go
var testData = []kdfTest{
	{
		key:  "746573742d6b6579",
		data: "746573742d696e707574",
		t0:   "6f0e5ad38daba1bea8a0d213688736f19763239305e0f58aba697f9ffc41c633",
		t1:   "df1194df20802a4fe594cde27e92991c8cae66c366e8106aaa937a55fa371e8a",
		t2:   "fac6e2745a325f5dc5d11a5b165aad08b0ada28e7b4e666b7c077934a4d76c24",
	},
	{
		key:  "776972656775617264",
		data: "776972656775617264",
		t0:   "491d43bbfdaa8750aaf535e334ecbfe5129967cd64635101c566d4caefda96e8",
		t1:   "1e71a379baefd8a79aa4662212fcafe19a23e2b609a3db7d6bcba8f560e3d25f",
		t2:   "31e1ae48bddfbe5de38f295e5452b1909a1b4e38e183926af3780b0c1e1f0160",
	},
	{
		key:  "",
		data: "",
		t0:   "8387b46bf43eccfcf349552a095d8315c4055beb90208fb1be23b894bc2ed5d0",
		t1:   "58a0e5f6faefccf4807bff1f05fa8a9217945762040bcec2f4b4a62bdfe0e86e",
		t2:   "0ce6ea98ec548f8e281e93e32db65621c45eb18dc6f0a7ad94178610a2f7338e",
	},
}

func runKDFTest(t *testing.T, count int) {
	for _, test := range testData {
		// decode test data
		key, _ := hex.DecodeString(test.key)
		data, _ := hex.DecodeString(test.data)

		// define output keys
		var key0, key1, key2 [blake2s.Size]byte

		// apply HKDF and output to key0, key1, key2
		hkdf.HKDF(key, data, &key0, &key1, &key2)

		// test HKDF with 1, 2, 3 keys generated
		switch count {
		case 3:
			assert.Equal(t, test.t2,
				hex.EncodeToString(key2[:]),
				"unexpected HKDF derivation")
			assert.Equal(t, test.t1,
				hex.EncodeToString(key1[:]),
				"unexpected HKDF derivation")
			assert.Equal(t, test.t0,
				hex.EncodeToString(key0[:]),
				"unexpected HKDF derivation")
		case 2:
			assert.Equal(t, test.t1,
				hex.EncodeToString(key1[:]),
				"unexpected HKDF derivation")
			assert.Equal(t, test.t0,
				hex.EncodeToString(key0[:]),
				"unexpected HKDF derivation")
		case 1:
			assert.Equal(t, test.t0,
				hex.EncodeToString(key0[:]),
				"unexpected HKDF derivation")
		}
	}
}

func TestHKDF(t *testing.T) {
	runKDFTest(t, 3)
	runKDFTest(t, 2)
	runKDFTest(t, 1)
}

package assembler

import (
	"os"
	"testing"

	"github.com/thee-engineer/cryptor/cachedb"
	"github.com/thee-engineer/cryptor/crypt"
)

func TestAssembler(t *testing.T) {
	t.Parallel()

	// Open cache
	cache, err := cachedb.NewLDBCache("data", 0, 0)
	if err != nil {
		t.Error(err)
	}
	defer cache.Close()

	// Create assembler
	hash, _ := crypt.DecodeString(
		"bbdb0ce1f3bae066a4bf8e232b216ae5692f9d428bd20e6d02cbc57f819423e5")
	asm := &Assembler{
		Tail:  hash,
		Cache: cache,
	}

	// Start assembling package
	err = asm.Assemble(crypt.NewKeyFromPassword("test"))
	if err != nil {
		t.Error(err)
	}

	os.Remove("untar")
}

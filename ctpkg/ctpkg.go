package ctpkg

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path"

	"github.com/thee-engineer/cryptor/archive"
	"github.com/thee-engineer/cryptor/chunker"
	"github.com/thee-engineer/cryptor/crypt"
	"github.com/thee-engineer/cryptor/utility"
)

const jsonExtension = "json"

// CTPKG ...
type CTPKG struct {
	Name      string
	Hash      string
	Tail      string
	Size      int
	ChunkSize uint32
	PKey      string
}

// NewCTPKG ...
func NewCTPKG(s, name string, chunkSize uint32, pKey crypt.AESKey) *CTPKG {
	contentBuffer := new(bytes.Buffer)

	// Create tar.gz of file/dir
	if err := archive.TarGz(s, contentBuffer); err != nil {
		panic(err)
	}

	// Hash tar.gz for integrity check
	contentHash := crypt.SHA256Data(contentBuffer.Bytes())

	// Get content lenght
	contentLen := contentBuffer.Len()

	// Generate a random primary key for the package
	if pKey == crypt.NullKey {
		pKey = crypt.NewKey()
	}

	// Create a chunker
	chunker := &chunker.Chunker{
		Size:   chunkSize,
		PKey:   pKey,
		Reader: contentBuffer,
	}

	// Start chunking the tar.gz
	tailHash, err := chunker.Chunk()
	if err != nil {
		panic(err)
	}

	// Create package info
	ctpkg := &CTPKG{
		Name:      name,
		Hash:      string(crypt.Encode(contentHash.Sum(nil))),
		Tail:      string(crypt.Encode(tailHash)),
		PKey:      pKey.String(),
		Size:      contentLen,
		ChunkSize: chunker.Size,
	}

	return ctpkg
}

// LoadCTPKG ...
func LoadCTPKG(ctpkgFile string) (ctpkg *CTPKG) {
	return ctpkg
}

// Assemble ...
func (ctpkg *CTPKG) Assemble(pKey crypt.AESKey) error {
	return nil
}

// ToJSON ...
func (ctpkg *CTPKG) ToJSON() ([]byte, error) {
	return json.MarshalIndent(ctpkg, "", "\t")
}

// Save ...
func (ctpkg *CTPKG) Save() error {
	// Create file name
	pkgFileName := fmt.Sprintf("%s.%s", ctpkg.Hash, jsonExtension)

	// Create file for CTPKG
	pkgFile, err := os.Create(path.Join(utility.GetPacksPath(), pkgFileName))
	if err != nil {
		return err
	}
	defer pkgFile.Close()

	// Convert CTPKG to JSON
	pkgJSON, err := ctpkg.ToJSON()
	if err != nil {
		return err
	}

	// Store CTPKG JSON to CTPKG File
	_, err = pkgFile.Write(pkgJSON)
	if err != nil {
		return err
	}

	return nil
}

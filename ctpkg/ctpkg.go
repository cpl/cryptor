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
	Size      int
	ChunkSize uint32
	PKey      string
}

// NewCTPKG ...
func NewCTPKG(source, name string, chunkSize uint32, pKey crypt.AESKey) *CTPKG {
	contentBuffer := new(bytes.Buffer)

	// Create tar.gz of file/dir
	if err := archive.TarGz(source, contentBuffer); err != nil {
		panic(err)
	}

	// Hash tar.gz for integrity check
	contentHash := crypt.SHA256Data(contentBuffer.Bytes())
	// Get content lenght
	contentLen := contentBuffer.Len()

	// Generate a random primary key for the package
	if pKey.Bytes() == nil {
		pKey = crypt.NewKey()
	}

	// Create a chunker
	chunker := &chunker.Chunker{
		Size:   chunkSize,
		Reader: contentBuffer,
	}

	// Start chunking the tar.gz
	if err := chunker.Chunk(); err != nil {
		panic(err)
	}

	// Create package info
	ctpkg := &CTPKG{
		Name:      name,
		Hash:      string(crypt.Encode(contentHash.Sum(nil))),
		PKey:      string(crypt.Encode(pKey[:])),
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
func Assemble(pKey *[32]byte) error {
	return nil
}

// ToJSON ...
func (ctpkg *CTPKG) ToJSON() ([]byte, error) {
	return json.MarshalIndent(ctpkg, "", "\t")
}

// Save ...
func (ctpkg *CTPKG) Save() error {
	pkgFileName := fmt.Sprintf("%s.%s", ctpkg.Hash, jsonExtension)
	pkgFile, err := os.Create(path.Join(utility.GetPacksPath(), pkgFileName))
	if err != nil {
		return err
	}
	defer pkgFile.Close()

	pkgJSON, err := ctpkg.ToJSON()
	if err != nil {
		return err
	}

	_, err = pkgFile.Write(pkgJSON)
	if err != nil {
		return err
	}

	return nil
}

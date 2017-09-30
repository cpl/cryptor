package ctpkg

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/thee-engineer/cryptor/archive"
	"github.com/thee-engineer/cryptor/cache"
	"github.com/thee-engineer/cryptor/chunker"
	"github.com/thee-engineer/cryptor/crypt"
)

const jsonExtension = "json"

// CTPKG ...
type CTPKG struct {
	Name      string `json:"name"`       // Name of the package
	Hash      string `json:"hash"`       // Hash of initial data
	Tail      string `json:"tail"`       // Tail chunk hash
	Size      int    `json:"size"`       // Size in bytes of initial data
	ChunkSize uint32 `json:"chunk_size"` // Chunk content size
	TKey      string `json:"tail_key"`   // Tail chunk key
}

// NewCTPKG ...
func NewCTPKG(s, n, o string, size uint32, tKey crypt.AESKey) (*CTPKG, error) {
	contentBuffer := new(bytes.Buffer)

	// Create tar.gz of file/dir
	if err := archive.TarGz(s, contentBuffer); err != nil {
		return nil, err
	}

	// Hash tar.gz for integrity check
	contentHash := crypt.SHA256Data(contentBuffer.Bytes())

	// Get content lenght
	contentLen := contentBuffer.Len()

	// Generate a random primary key for the package
	if tKey == crypt.NullKey {
		tKey = crypt.NewKey()
	}

	// Create a chunker
	chunker := &chunker.Chunker{
		Size:   size,
		Cache:  o,
		Reader: contentBuffer,
	}

	// Start chunking the tar.gz
	tailHash, err := chunker.Chunk(tKey)
	if err != nil {
		return nil, err
	}

	// Create package info
	ctpkg := &CTPKG{
		Name:      n,
		Hash:      string(crypt.Encode(contentHash.Sum(nil))),
		Tail:      string(crypt.Encode(tailHash)),
		TKey:      tKey.String(),
		Size:      contentLen,
		ChunkSize: chunker.Size,
	}

	return ctpkg, nil
}

// LoadCTPKG ...
func LoadCTPKG(ctpkgHash, source string) (*CTPKG, error) {
	// Default source directory
	if source == "" {
		source = cache.GetPacksPath()
	}

	// Obtain packs path
	packName := fmt.Sprintf("%s.%s", ctpkgHash, jsonExtension)
	packPath := path.Join(source, packName)

	// Check that pack exist
	_, err := os.Stat(packPath)
	if err != nil {
		return nil, err
	}

	// Obtain pack contents
	data, err := ioutil.ReadFile(packPath)
	if err != nil {
		return nil, err
	}

	// Convert JSON to CTPKG
	ctpkg := &CTPKG{}
	if err = json.Unmarshal(data, ctpkg); err != nil {
		return nil, err
	}

	return ctpkg, nil
}

func (ctpkg CTPKG) toJSON() ([]byte, error) {
	return json.MarshalIndent(ctpkg, "", "\t")
}

// Save ...
func (ctpkg CTPKG) Save(destination string) error {
	// Create file name
	pkgFileName := fmt.Sprintf("%s.%s", ctpkg.Tail, jsonExtension)

	// Resort to default location
	if destination == "" {
		destination = cache.GetPacksPath()
	}

	// Create file for CTPKG
	pkgFile, err := os.Create(path.Join(destination, pkgFileName))
	if err != nil {
		return err
	}
	defer pkgFile.Close()

	// Convert CTPKG to JSON
	pkgJSON, err := ctpkg.toJSON()
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

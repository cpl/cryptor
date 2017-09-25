package ctpkg

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/thee-engineer/cryptor/archive"
	"github.com/thee-engineer/cryptor/assembler"
	"github.com/thee-engineer/cryptor/cache"
	"github.com/thee-engineer/cryptor/chunker"
	"github.com/thee-engineer/cryptor/crypt"
)

const jsonExtension = "json"

// CTPKG ...
type CTPKG struct {
	Name      string `json:"name"`
	Hash      string `json:"hash"`
	Tail      string `json:"tail"`
	Size      int    `json:"size"`
	ChunkSize uint32 `json:"chunk_size"`
	TKey      string `json:"tail_key"`
}

// NewCTPKG ...
func NewCTPKG(s, name string, chunkSize uint32, tKey crypt.AESKey) *CTPKG {
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
	if tKey == crypt.NullKey {
		tKey = crypt.NewKey()
	}

	// Create a chunker
	chunker := &chunker.Chunker{
		Size:   chunkSize,
		Reader: contentBuffer,
	}

	// Start chunking the tar.gz
	tailHash, err := chunker.Chunk(tKey)
	if err != nil {
		panic(err)
	}

	// Create package info
	ctpkg := &CTPKG{
		Name:      name,
		Hash:      string(crypt.Encode(contentHash.Sum(nil))),
		Tail:      string(crypt.Encode(tailHash)),
		TKey:      tKey.String(),
		Size:      contentLen,
		ChunkSize: chunker.Size,
	}

	return ctpkg
}

// LoadCTPKG ...
func LoadCTPKG(ctpkgHash string) (ctpkg *CTPKG) {
	// Obtain packs path
	packsPath := cache.GetPacksPath()
	packName := fmt.Sprintf("%s.%s", ctpkgHash, jsonExtension)
	packPath := path.Join(packsPath, packName)

	// Check that pack exist
	_, err := os.Stat(packPath)
	if err != nil {
		panic(err)
	}

	// Obtain pack contents
	data, err := ioutil.ReadFile(packPath)
	if err != nil {
		panic(err)
	}

	// Convert JSON to CTPKG
	if err = json.Unmarshal(data, ctpkg); err != nil {
		panic(err)
	}

	return ctpkg
}

// Assemble ...
func (ctpkg *CTPKG) Assemble() error {
	tKey := crypt.NewKeyFromString(ctpkg.TKey)
	assembler.Assemble(ctpkg.Tail, tKey)
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
	pkgFile, err := os.Create(path.Join(cache.GetPacksPath(), pkgFileName))
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

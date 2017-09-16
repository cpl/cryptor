package chunker

import (
	"fmt"
	"io/ioutil"
)

const (
	ppFormatString = `
Hash: %s
Name: %s
Size: %d

ChunkSize: %d
ChunkCount: %d
`

	ppPathFormatString = "%s/%s.%s"
	ppExtension        = ".ctpp"
)

// PackageProfile ...
type PackageProfile struct {
	Hash       string
	Name       string
	Size       uint64
	ChunkSize  uint64
	ChunkCount uint64
}

// Generate ...
func (p *PackageProfile) Generate(filePath string) error {
	// Format the packageProfile string
	ppString := fmt.Sprintf(ppFormatString,
		p.Hash, p.Name, p.Size, p.ChunkSize, p.ChunkCount)

	// Write the formated string to a given path
	ppPath := fmt.Sprintf(ppPathFormatString, filePath, p.Hash, ppExtension)
	ioutil.WriteFile(ppPath, []byte(ppString), 0400)
	return nil
}

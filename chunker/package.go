package chunker

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

const (
	ppPathFormatString = "%s/%s.%s"
	ppExtension        = "json"
	ppFileName         = "profile"
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
	// Prepare package profile path
	ppPath := fmt.Sprintf(ppPathFormatString,
		filePath, ppFileName, ppExtension)

	// Generate JSON data for package profile
	jsonData, err := json.Marshal(p)
	if err != nil {
		return err
	}

	// Write JSON data to package profile file
	err = ioutil.WriteFile(ppPath, jsonData, 0400)
	if err != nil {
		return err
	}

	return nil
}

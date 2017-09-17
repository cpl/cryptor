package ctpkm

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

const (
	ppPathFormatString = "%s/%s.%s"

	// PPExtension ...
	PPExtension = "json"

	// PPFileName ...
	PPFileName = "profile"
)

// PackageProfile ...
type PackageProfile struct {
	Hash       string
	Name       string
	Size       uint
	ChunkSize  uint
	ChunkCount uint
}

// NewPackageProfile ...
func NewPackageProfile(packagePath, packageName string) *PackageProfile {

	return &PackageProfile{"", "", 0, 0, 0}
}

// Generate ...
func (p *PackageProfile) Generate(filePath string) error {
	// Prepare package profile path
	ppPath := fmt.Sprintf(ppPathFormatString,
		filePath, PPFileName, PPExtension)

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

package utility

import (
	"archive/tar"
	"compress/gzip"
)

// TarFile implement *tar.Writer
type TarFile struct {
	Writer     *tar.Writer
	Name       string
	GzWriter   *gzip.Writer
	Compressed bool
}

// Tar ...
func Tar() {

}

// UnTar ...
func UnTar() {

}

package archive

import (
	"compress/gzip"
	"io"
)

// TarGz takes a source (path to file/directory) as input, walks the filepath
// and creates a .tar.gz archive which is written on the given io.Writer.
func TarGz(source string, out io.Writer) error {
	gzipWriter := gzip.NewWriter(out)
	if err := tarArchive(source, gzipWriter); err != nil {
		return err
	}
	defer gzipWriter.Close()
	return nil
}

// UnTarGz takes io.Reader as input which should contain the .tar.gz archive
// data and outputs the original files/directories at the given destination.
func UnTarGz(destination string, in io.Reader) error {
	// Create gzip reader
	gzipReader, err := gzip.NewReader(in)
	if err != nil {
		return err
	}
	defer gzipReader.Close()
	return tarExtract(destination, gzipReader)
}

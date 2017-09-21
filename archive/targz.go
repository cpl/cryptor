package archive

import (
	"compress/gzip"
	"io"
)

// TarGz ...
func TarGz(source string, out io.Writer) error {
	gzipWriter := gzip.NewWriter(out)
	if err := tarArchive(source, gzipWriter); err != nil {
		return err
	}
	defer gzipWriter.Close()
	return nil
}

// UnTarGz ...
func UnTarGz(destination string, in io.Reader) error {
	// Create gzip reader
	gzipReader, err := gzip.NewReader(in)
	if err != nil {
		return err
	}
	defer gzipReader.Close()
	return tarExtract(destination, gzipReader)
}

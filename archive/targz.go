package archive

import (
	"compress/gzip"
	"io"
)

// TarGz ...
func TarGz(source string, out io.Writer) error {
	gzipWriter := gzip.NewWriter(out)
	defer gzipWriter.Close()
	if err := TarArchive(source, gzipWriter); err != nil {
		return err
	}
	return nil
}

// UnTarGz ...
func UnTarGz(destination string, in io.Reader) error {
	// Create gzip reader
	gzipReader, err := gzip.NewReader(in)
	defer gzipReader.Close()
	if err != nil {
		return err
	}
	return tarExtract(destination, gzipReader)
}

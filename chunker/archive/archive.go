//Package archive contains functions for reading files & directories and creating
// tar archives (https://golang.org/pkg/archive/tar/).
// This package also contains functions for compressing data using gzip
// (https://golang.org/pkg/compress/gzip/).
package archive

import (
	"archive/tar"
	"io"
	"os"
	"path/filepath"
)

func tarArchive(source string, out io.Writer) error {

	// Prepare tar writer
	tarWriter := tar.NewWriter(out)
	defer tarWriter.Close()

	// Walk source tree
	return filepath.Walk(source,
		func(file string, fileInfo os.FileInfo, err error) error {
			// Check for errors
			if err != nil {
				return err
			}

			// Check for symlinks and handle it
			var link string
			if fileInfo.Mode()&os.ModeSymlink == os.ModeSymlink {
				if link, err = os.Readlink(source); err != nil {
					return nil
				}
			}

			// Create tar header
			header, err := tar.FileInfoHeader(fileInfo, link)
			if err != nil {
				return err
			}

			// Update tar header to relative path
			if file == source {
				header.Name = ""
			} else {
				header.Name = file[len(source)+1:]
			}

			// Write file/dir header to tar
			if err := tarWriter.WriteHeader(header); err != nil {
				return err
			}

			// Directories have no content, return
			if fileInfo.Mode().IsDir() {
				return nil
			}

			// Open file for reading
			f, err := os.Open(file)
			if err != nil {
				return err
			}
			defer f.Close()

			// Copy file content to tar writer
			if _, err := io.Copy(tarWriter, f); err != nil {
				return err
			}

			return nil
		})
}

func tarExtract(destination string, in io.Reader) error {
	tarReader := tar.NewReader(in)

	for {
		// Read each header
		header, err := tarReader.Next()

		// Check for errors, EOF, or empty headers
		switch {
		case err == io.EOF:
			return nil
		case err != nil:
			return err
		case header == nil:
			return nil
		}

		// Location where dir/file should be created
		target := filepath.Join(destination, header.Name)

		switch header.Typeflag {
		// Check for directory
		case tar.TypeDir:
			if _, err := os.Stat(target); err != nil {
				if err := os.MkdirAll(target, 0755); err != nil {
					return err
				}
			}

		// Check for file
		case tar.TypeReg:
			// Create file, starting with header
			file, err := os.OpenFile(
				target, os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode))
			if err != nil {
				return err
			}
			defer file.Close()

			// Copy file contents to new file
			if _, err := io.Copy(file, tarReader); err != nil {
				return err
			}
		}
	}
}

package archive

import (
	"archive/tar"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// TarArchive ...
func tarArchive(source string, out io.Writer) error {
	// Get file/dir status
	stat, err := os.Stat(source)
	if err != nil {
		return err
	}

	// Prepare tar writer
	tarWriter := tar.NewWriter(out)
	defer tarWriter.Close()

	// Check if source is file or dir
	if stat.IsDir() {
		return tarDir(source, tarWriter)
	}

	return tarFile(stat, tarWriter)
}

func tarFile(fileInfo os.FileInfo, tarWriter *tar.Writer) error {
	// Create tar header
	header, err := tar.FileInfoHeader(fileInfo, fileInfo.Name())
	if err != nil {
		return err
	}

	// Write file header to tar
	if err := tarWriter.WriteHeader(header); err != nil {
		return err
	}

	// Open file for reading
	f, err := os.Open(fileInfo.Name())
	defer f.Close()
	if err != nil {
		return err
	}

	// Copy file content to tar writer
	if _, err := io.Copy(tarWriter, f); err != nil {
		return err
	}

	return nil
}

func tarDir(source string, tarWriter *tar.Writer) error {
	// Walk source tree
	return filepath.Walk(source,
		func(file string, fileInfo os.FileInfo, err error) error {
			// Check for errors
			if err != nil {
				return err
			}

			// Ignore root directory
			if source == file {
				return nil
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

			// Update header for directories
			header.Name = strings.TrimPrefix(
				strings.Replace(file, source, "", -1),
				string(filepath.Separator))

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
			defer f.Close()
			if err != nil {
				return err
			}

			// Copy file content to tar writer
			if _, err := io.Copy(tarWriter, f); err != nil {
				return err
			}

			return nil
		})
}

func tarExtract(destination string, in io.Reader) error {
	tarReader := tar.NewReader(in)

	// Make sure destination exists
	if _, err := os.Stat(destination); err != nil {
		if err := os.MkdirAll(destination, 0755); err != nil {
			return err
		}
	}

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

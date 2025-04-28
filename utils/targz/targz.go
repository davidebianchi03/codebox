package targz

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type TarGZManager struct {
	Filepath string
}

// create an empty tar.gz archive
func (tgm *TarGZManager) CreateArchive() error {
	file, err := os.Create(tgm.Filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	gw := gzip.NewWriter(file)
	defer gw.Close()

	tw := tar.NewWriter(gw)
	defer tw.Close()

	return nil
}

// write file inside a tar.gz archive,
// if file already exists it will be overwritten
func (tgm *TarGZManager) WriteFile(filename string, data []byte) error {
	entries, err := tgm.readAll()
	if err != nil {
		return err
	}

	entries[filename] = data

	return tgm.writeAll(entries)
}

// ReadFile, read file from archive
func (tgm *TarGZManager) ReadFile(filename string) ([]byte, error) {
	file, err := os.Open(tgm.Filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	gr, err := gzip.NewReader(file)
	if err != nil {
		return nil, err
	}
	defer gr.Close()

	tr := tar.NewReader(gr)

	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break // End of archive
		}
		if err != nil {
			return nil, err
		}

		if hdr.Name == filename {
			var buf bytes.Buffer
			if _, err := io.Copy(&buf, tr); err != nil {
				return nil, err
			}
			return buf.Bytes(), nil
		}
	}

	return nil, fmt.Errorf("file %s not found in archive", filename)
}

// DeleteFile removes a file from the archive
func (tgm *TarGZManager) RemoveFileFromArchive(filename string) error {
	entries, err := tgm.readAll()
	if err != nil {
		return err
	}

	if _, exists := entries[filename]; !exists {
		return fmt.Errorf("file %s not found in archive", filename)
	}

	delete(entries, filename)

	return tgm.writeAll(entries)
}

// compress folder into a tar.gz archive
func (tgm *TarGZManager) CompressFolder(srcDir string) error {
	// Create the output file
	outFile, err := os.Create(tgm.Filepath)
	if err != nil {
		return fmt.Errorf("failed to create archive file: %w", err)
	}
	defer outFile.Close()

	// Create a gzip writer
	gzipWriter := gzip.NewWriter(outFile)
	defer gzipWriter.Close()

	// Create a tar writer
	tarWriter := tar.NewWriter(gzipWriter)
	defer tarWriter.Close()

	// Walk the directory and add files to the archive
	err = filepath.Walk(srcDir, func(file string, fi os.FileInfo, err error) error {
		if file == srcDir {
			return nil
		}

		if err != nil {
			return err
		}

		// Create header for the file
		header, err := tar.FileInfoHeader(fi, file)
		if err != nil {
			return fmt.Errorf("failed to get file info header: %w", err)
		}

		// Replace absolute paths with relative ones
		header.Name = strings.TrimPrefix(strings.Replace(file, srcDir, "", -1), string(filepath.Separator))

		// Write the header to the tarball
		if err := tarWriter.WriteHeader(header); err != nil {
			return fmt.Errorf("failed to write header: %w", err)
		}

		// If it's a directory, we don't need to write content
		if fi.IsDir() {
			return nil
		}

		// Open the file for reading
		fileToArchive, err := os.Open(file)
		if err != nil {
			return fmt.Errorf("failed to open file: %w", err)
		}
		defer fileToArchive.Close()

		// Copy the file content to the tar writer
		if _, err := io.Copy(tarWriter, fileToArchive); err != nil {
			return fmt.Errorf("failed to copy file content: %w", err)
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("error walking the directory: %w", err)
	}

	return nil
}

// extra tar.gz archive into a folder
func (tgm *TarGZManager) ExtractTarGz(destination string) error {
	tarGzFile, err := os.Open(tgm.Filepath)
	if err != nil {
		return fmt.Errorf("failed to open tar.gz file: %v", err)
	}
	defer tarGzFile.Close()

	uncompressedStream, err := gzip.NewReader(tarGzFile)
	if err != nil {
		return fmt.Errorf("failed to read gzip: %v", err)
	}
	defer uncompressedStream.Close()

	tarReader := tar.NewReader(uncompressedStream)

	for {
		header, err := tarReader.Next()

		if err == io.EOF {
			break
		}

		if err != nil {
			return fmt.Errorf("failed to open tar file: %v", err)
		}

		target := filepath.Join(destination, header.Name)

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(target, os.FileMode(header.Mode)); err != nil {
				return fmt.Errorf("failed to create output directory: %v", err)
			}

		case tar.TypeReg:
			if err := os.MkdirAll(filepath.Dir(target), 0755); err != nil {
				return fmt.Errorf("failed to create parent directory: %v", err)
			}

			outFile, err := os.Create(target)
			if err != nil {
				return fmt.Errorf("failed to create file: %v", err)
			}
			defer outFile.Close()

			if _, err := io.Copy(outFile, tarReader); err != nil {
				return fmt.Errorf("failed to write file: %v", err)
			}

		default:
			fmt.Printf("Tipo di file non gestito: %c nel file %s\n", header.Typeflag, header.Name)
		}
	}

	return nil
}

// readAll reads all files inside the tar.gz archive into a map
func (tgm *TarGZManager) readAll() (map[string][]byte, error) {
	entries := make(map[string][]byte)

	file, err := os.Open(tgm.Filepath)
	if err != nil {
		if os.IsNotExist(err) {
			return entries, nil // Return empty if file doesn't exist yet
		}
		return nil, err
	}
	defer file.Close()

	gr, err := gzip.NewReader(file)
	if err != nil {
		return nil, err
	}
	defer gr.Close()

	tr := tar.NewReader(gr)

	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		var buf bytes.Buffer
		if _, err := io.Copy(&buf, tr); err != nil {
			return nil, err
		}

		entries[hdr.Name] = buf.Bytes()
	}

	return entries, nil
}

// writeAll writes all files from the map into the tar.gz archive
func (tgm *TarGZManager) writeAll(entries map[string][]byte) error {
	file, err := os.Create(tgm.Filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	gw := gzip.NewWriter(file)
	defer gw.Close()

	tw := tar.NewWriter(gw)
	defer tw.Close()

	for name, data := range entries {
		hdr := &tar.Header{
			Name:    name,
			Mode:    0600,
			Size:    int64(len(data)),
			ModTime: time.Now(),
		}

		if err := tw.WriteHeader(hdr); err != nil {
			return err
		}

		if _, err := tw.Write(data); err != nil {
			return err
		}
	}

	return nil
}

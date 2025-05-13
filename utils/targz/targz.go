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

type TarEntry struct {
	Path    string `json:"name"`
	Type    string `json:"type"`
	Content []byte `json:"content"`
}

type TarTreeItem struct {
	Name     string         `json:"name"`
	FullPath string         `json:"full_path"`
	Type     string         `json:"type"`
	Children []*TarTreeItem `json:"children"`
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
func (tgm *TarGZManager) WriteFile(path string, data []byte) error {
	entries, err := tgm.ListEntries()
	if err != nil {
		return err
	}

	found := false
	for i, e := range entries {
		if e.Path == path {
			entries[i].Content = data
			found = true
		}
	}

	if !found {
		entries = append(entries, TarEntry{
			Path:    path,
			Type:    "file",
			Content: data,
		})
	}

	return tgm.writeAll(entries)
}

// create a folder
func (tgm *TarGZManager) Mkdir(path string) error {
	entries, err := tgm.ListEntries()
	if err != nil {
		return err
	}

	entries = append(entries, TarEntry{
		Path:    path,
		Type:    "dir",
		Content: []byte{},
	})

	return tgm.writeAll(entries)
}

// Delete removes a path and its children elements
func (tgm *TarGZManager) Delete(path string) error {
	entries, err := tgm.ListEntries()
	if err != nil {
		return err
	}

	newEntries := make([]TarEntry, 0)

	for _, entry := range entries {
		if !strings.HasPrefix(entry.Path, path) && strings.TrimSuffix(entry.Path, "/") != strings.TrimSuffix(path, "/") {
			newEntries = append(newEntries, entry)
		}
	}

	return tgm.writeAll(newEntries)
}

func (tgm *TarGZManager) Move(oldPath string, newPath string) error {
	entries, err := tgm.ListEntries()
	if err != nil {
		return err
	}

	newEntries := make([]TarEntry, 0)

	for _, entry := range entries {
		if strings.HasPrefix(entry.Path, oldPath) || strings.TrimSuffix(entry.Path, "/") == strings.TrimSuffix(oldPath, "/") {
			entry.Path = strings.Replace(entry.Path, oldPath, newPath, 1)
		}
		newEntries = append(newEntries, entry)
	}

	return tgm.writeAll(newEntries)
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

func (tgm *TarGZManager) ListEntries() ([]TarEntry, error) {
	entries := make([]TarEntry, 0)
	file, err := os.Open(tgm.Filepath)
	if err != nil {
		if os.IsNotExist(err) {
			return []TarEntry{}, nil // Return empty if file doesn't exist yet
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

		if hdr.Typeflag == tar.TypeDir {
			entries = append(entries, TarEntry{
				Path:    hdr.Name,
				Type:    "dir",
				Content: []byte{},
			})
		} else {
			// read content
			var buf bytes.Buffer
			if _, err := io.Copy(&buf, tr); err != nil {
				return nil, err
			}

			entries = append(entries, TarEntry{
				Path:    hdr.Name,
				Type:    "file",
				Content: buf.Bytes(),
			})
		}
	}

	return entries, nil
}

func (tgm *TarGZManager) EntriesTree() ([]*TarTreeItem, error) {
	entries, err := tgm.ListEntries()
	if err != nil {
		return nil, err
	}

	var tree *[]*TarTreeItem
	t := make([]*TarTreeItem, 0)
	tree = &t

	for _, entry := range entries {
		// remove trailing slash
		entryName := strings.TrimSuffix(entry.Path, "/")

		parts := strings.Split(entryName, "/")

		if len(parts) > 0 {
			// check if exists
			var parentTreeItem *TarTreeItem

			for _, part := range parts {
				var source *[]*TarTreeItem
				if parentTreeItem == nil {
					source = tree
				} else {
					source = &parentTreeItem.Children
				}

				found := false
				for _, ti := range *source {
					if ti.Name == part {
						found = true
						parentTreeItem = ti
					}
				}

				if !found {
					item := TarTreeItem{
						Name:     part,
						FullPath: strings.TrimPrefix(entry.Path, "./"),
						Children: make([]*TarTreeItem, 0),
						Type:     entry.Type,
					}

					*source = append(*source, &item)
					parentTreeItem = &item
				}
			}
		}
	}

	return *tree, nil
}

// retrieve entry by path, return nil if not found
func (tgm *TarGZManager) RetrieveEntry(path string) (*TarEntry, error) {
	entries, err := tgm.ListEntries()
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if strings.TrimSuffix(entry.Path, "/") == strings.TrimSuffix(path, "/") {
			return &entry, nil
		}
	}

	return nil, nil
}

// writeAll writes all files from the map into the tar.gz archive
func (tgm *TarGZManager) writeAll(entries []TarEntry) error {
	file, err := os.Create(tgm.Filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	gw := gzip.NewWriter(file)
	defer gw.Close()

	tw := tar.NewWriter(gw)
	defer tw.Close()

	for _, entry := range entries {
		tarTypeFlag := tar.TypeReg // file
		if entry.Type == "dir" {
			tarTypeFlag = tar.TypeDir
		}

		hdr := &tar.Header{
			Name:     entry.Path,
			Typeflag: byte(tarTypeFlag),
			Mode:     0600,
			Size:     int64(len(entry.Content)),
			ModTime:  time.Now(),
		}

		if err := tw.WriteHeader(hdr); err != nil {
			return err
		}

		if _, err := tw.Write(entry.Content); err != nil {
			return err
		}
	}

	return nil
}

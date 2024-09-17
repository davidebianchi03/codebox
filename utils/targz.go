package utils

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func CreateNewTarGzArchive(srcDir, outputFilePath string) error {
	// Create the output file
	outFile, err := os.Create(outputFilePath)
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

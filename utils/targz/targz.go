package targz

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func CreateArchive(srcDir string, outputFilePath string) error {
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

func ExtractTarGz(tarGzPath, destination string) error {
	// Aprire il file tar.gz
	tarGzFile, err := os.Open(tarGzPath)
	if err != nil {
		return fmt.Errorf("errore nell'apertura del file tar.gz: %v", err)
	}
	defer tarGzFile.Close()

	// Decomprimere il gzip
	uncompressedStream, err := gzip.NewReader(tarGzFile)
	if err != nil {
		return fmt.Errorf("errore nella lettura del gzip: %v", err)
	}
	defer uncompressedStream.Close()

	// Leggere l'archivio tar
	tarReader := tar.NewReader(uncompressedStream)

	for {
		// Leggere la prossima intestazione di file
		header, err := tarReader.Next()

		// Se raggiungiamo la fine dell'archivio, esci
		if err == io.EOF {
			break
		}

		// Gestire altri errori
		if err != nil {
			return fmt.Errorf("errore nella lettura del tar: %v", err)
		}

		// Creare il percorso completo
		target := filepath.Join(destination, header.Name)

		// Controllare il tipo di file
		switch header.Typeflag {
		case tar.TypeDir:
			// Creare la directory se non esiste già
			if err := os.MkdirAll(target, os.FileMode(header.Mode)); err != nil {
				return fmt.Errorf("errore nella creazione della directory: %v", err)
			}

		case tar.TypeReg:
			// Creare la directory padre, se necessario
			if err := os.MkdirAll(filepath.Dir(target), 0755); err != nil {
				return fmt.Errorf("errore nella creazione della directory padre: %v", err)
			}

			// Creare il file
			outFile, err := os.Create(target)
			if err != nil {
				return fmt.Errorf("errore nella creazione del file: %v", err)
			}
			defer outFile.Close()

			// Scrivere il contenuto del file
			if _, err := io.Copy(outFile, tarReader); err != nil {
				return fmt.Errorf("errore nella scrittura del file: %v", err)
			}

		default:
			// Se ci sono altri tipi di file, ignorarli (puoi gestirli se necessario)
			fmt.Printf("Tipo di file non gestito: %c nel file %s\n", header.Typeflag, header.Name)
		}
	}

	return nil
}

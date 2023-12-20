package model

import (
	"os"
	"path/filepath"
	"strings"
)

const hiperDataPath = "/var/hiper"

// saveFile saves a file to the hiper data directory.
func saveFile(relativePath string, content []byte) error {
	fullPath := filepath.Join(hiperDataPath, relativePath)
	// Create the directory path
	dirPath := filepath.Dir(fullPath)
	if err := os.MkdirAll(dirPath, 0755); err != nil {
		return err
	}

	// Write to a temporary file first, then rename it to the target file.
	// Because renaming is atomic in Linux.
	tempFile, err := os.CreateTemp(dirPath, "temp")
	if err != nil {
		return err
	}
	defer tempFile.Close()
	if _, err = tempFile.Write(content); err != nil {
		return err
	}
	// Rename the temporary file to the target file
	if err = os.Rename(tempFile.Name(), fullPath); err != nil {
		return err
	}

	return nil
}

// getFileWithAutoExt returns a file from the hiper data directory.
// Automatically add the extension to the file name if needed.
func getFileWithAutoExt(relativePath string) ([]byte, error) {
	dirPath := filepath.Join(hiperDataPath, filepath.Dir(relativePath))
	baseName := filepath.Base(relativePath)

	// Open the directory
	dir, err := os.Open(dirPath)
	if err != nil {
		return nil, err
	}
	defer dir.Close()

	// List all files in the directory
	entries, err := dir.Readdir(-1)
	if err != nil {
		return nil, err
	}

	// Find the file with the same base name
	for _, entry := range entries {
		if strings.HasPrefix(entry.Name(), baseName) {
			// Read the file
			data, err := os.ReadFile(filepath.Join(dirPath, entry.Name()))
			if err != nil {
				return nil, err
			}
			return data, nil
		}
	}

	return nil, os.ErrNotExist
}

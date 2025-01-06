package main

import (
	"errors"
	"fmt"
	"log"
	"path/filepath"

	"github.com/spf13/afero"
)

var (
	ErrFileNotFound   = errors.New("file not found")
	ErrDirNotFound    = errors.New("directory not found")
	ErrInvalidFileExt = errors.New("invalid file extension")
	ErrInvalidPath    = errors.New("invalid path")
)

func main() {
	// Create a new virtual file system
	fs := afero.NewMemMapFs()

	// Setup some mock files and directories
	err := setupMockFilesystem(fs)
	if err != nil {
		log.Fatalf("Error setting up mock filesystem: %v", err)
	}

	// Perform some invalid operations to test error handling
	performInvalidOperations(fs)
}

func setupMockFilesystem(fs afero.Fs) error {
	// Create some directories and files
	// ... (same code as before)

	return nil
}

func performInvalidOperations(fs afero.Fs) {
	// Try to read a non-existent file
	content, err := afero.ReadFile(fs, "/non-existent-file.txt")
	if err != nil {
		fmt.Println("Error reading non-existent file:", err)
	} else {
		fmt.Println("Content of non-existent file:", string(content))
	}

	// Try to write a file with an invalid extension
	err = afero.WriteFile(fs, "/invalid.file", []byte("Invalid file content"), 0644)
	if err != nil {
		fmt.Println("Error writing file with invalid extension:", err)
	}

	// Try to delete a non-existent directory
	err = fs.RemoveAll("/non-existent-dir")
	if err != nil {
		fmt.Println("Error deleting non-existent directory:", err)
	}

	// Try to create a directory with an invalid path
	err = fs.MkdirAll("invalid/path/to/dir", 0755)
	if err != nil {
		fmt.Println("Error creating directory with invalid path:", err)
	}
}

func validatePath(path string) error {
	if path == "" {
		return ErrInvalidPath
	}
	if !filepath.IsAbs(path) {
		return ErrInvalidPath
	}
	return nil
}

func validateFileExt(filePath string) error {
	ext := filepath.Ext(filePath)
	if ext != ".txt" && ext != ".bin" {
		return ErrInvalidFileExt
	}
	return nil
}

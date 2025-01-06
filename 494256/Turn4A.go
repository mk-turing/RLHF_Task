package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/afero"
)

func main() {
	// Create a new virtual file system
	fs := afero.NewMemMapFs()

	// Setup some mock files and directories
	err := setupMockFilesystem(fs)
	if err != nil {
		log.Fatalf("Error setting up mock filesystem: %v", err)
	}

	// Traverse and list the virtual filesystem
	rootDir := "/"
	traverseAndList(fs, rootDir)

	// Create invalid operations to test error handling
	testInvalidOperations(fs)
}

func setupMockFilesystem(fs afero.Fs) error {
	// Create some directories
	err := fs.MkdirAll("/docs/subdir1", 0755)
	if err != nil {
		return err
	}
	err = fs.MkdirAll("/docs/subdir2", 0755)
	if err != nil {
		return err
	}

	// Create some files
	err = afero.WriteFile(fs, "/docs/myfile.txt", []byte("Hello from docs/myfile.txt\n"), 0644)
	if err != nil {
		return err
	}
	err = afero.WriteFile(fs, "/docs/subdir1/subfile1.txt", []byte("Hello from docs/subdir1/subfile1.txt\n"), 0644)
	if err != nil {
		return err
	}
	err = afero.WriteFile(fs, "/data/datafile.bin", []byte("This is a binary file"), 0600)
	if err != nil {
		return err
	}

	return nil
}

func traverseAndList(fs afero.Fs, rootDir string) {
	err := afero.Walk(fs, rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Remove the root directory from the path to print relative paths
		path = strings.TrimPrefix(path, rootDir)

		var fileSize int64
		if !info.IsDir() {
			fileSize = info.Size()
		}

		fmt.Printf("%s (%s) - %d bytes\n", path, filepath.Ext(path), fileSize)
		return nil
	})
	if err != nil {
		log.Fatalf("Error walking through filesystem: %v", err)
	}
}

func testInvalidOperations(fs afero.Fs) {
	// Test invalid directory creation
	invalidDir := "/invalid/path/to/create/"
	err := fs.MkdirAll(invalidDir, 0755)
	if err != nil {
		fmt.Printf("Error creating directory %s: %v\n", invalidDir, err)
	}

	// Test file creation over existing directory
	existingDir := "/docs"
	err = afero.WriteFile(fs, existingDir, []byte("This should not work"), 0644)
	if err != nil {
		fmt.Printf("Error creating file %s: %v\n", existingDir, err)
	}

	// Test file deletion that doesn't exist
	nonExistentFile := "/non/existent/file.txt"
	err = fs.Remove(nonExistentFile)
	if err != nil {
		fmt.Printf("Error deleting file %s: %v\n", nonExistentFile, err)
	}

	// Test directory deletion that doesn't exist
	nonExistentDir := "/non/existent/dir"
	err = fs.RemoveAll(nonExistentDir)
	if err != nil {
		fmt.Printf("Error deleting directory %s: %v\n", nonExistentDir, err)
	}
}

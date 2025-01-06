package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/afero"
)

func main() {
	// Create a new virtual file system
	fs := afero.NewMemMapFs()

	// Create a directory structure
	dirPath := "/docs"
	err := fs.MkdirAll(dirPath, 0755)
	if err != nil {
		log.Fatalf("Error creating directory: %v", err)
	}

	filePaths := []string{
		dirPath + "/myfile.txt",
		dirPath + "/subdir/subfile.txt",
		dirPath + "/subdir/subsubdir/deepfile.txt",
	}

	for _, filePath := range filePaths {
		err = afero.WriteFile(fs, filePath, []byte("Sample content"), 0644)
		if err != nil {
			log.Fatalf("Error creating file: %v", err)
		}
	}

	// Traverse and list files and directories
	err = filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Print file or directory information
		fmt.Printf("%s - %s, %d bytes\n", info.Mode(), filepath.Base(path), info.Size())

		return nil
	})
	if err != nil {
		log.Fatalf("Error walking directory: %v", err)
	}
}

package main

import (
	"fmt"
	"log"

	"github.com/spf13/afero"
)

func main() {
	// Create a new virtual file system
	fs := afero.NewMemMapFs()

	// Create a directory
	dirPath := "/docs"
	err := fs.MkdirAll(dirPath, 0755)
	if err != nil {
		log.Fatalf("Error creating directory: %v", err)
	}
	fmt.Printf("Created directory: %s\n", dirPath)

	// Create a file inside the directory
	filePath := dirPath + "/myfile.txt"
	err = afero.WriteFile(fs, filePath, []byte("Hello, this is a mock file!\n"), 0644)
	if err != nil {
		log.Fatalf("Error creating file: %v", err)
	}
	fmt.Printf("Created file: %s\n", filePath)

	// Read the file
	content, err := afero.ReadFile(fs, filePath)
	if err != nil {
		log.Fatalf("Error reading file: %v", err)
	}
	fmt.Printf("File contents: %s\n", string(content))

	// Append data to the file
	newContent := []byte("\nAdded some more content.")
	err = afero.WriteFile(fs, filePath, append(content, newContent...), 0644)
	if err != nil {
		log.Fatalf("Error writing to file: %v", err)
	}

	// Read the updated file
	content, err = afero.ReadFile(fs, filePath)
	if err != nil {
		log.Fatalf("Error reading file: %v", err)
	}
	fmt.Printf("Updated file contents: %s\n", string(content))

	// Delete the file
	err = fs.Remove(filePath)
	if err != nil {
		log.Fatalf("Error deleting file: %v", err)
	}
	fmt.Printf("Deleted file: %s\n", filePath)

	// Check if the file exists
	exists, err := afero.Exists(fs, filePath)
	if err != nil {
		log.Fatalf("Error checking file existence: %v", err)
	}
	fmt.Printf("File exists: %t\n", exists)

	// Delete the directory and its contents
	err = fs.RemoveAll(dirPath)
	if err != nil {
		log.Fatalf("Error deleting directory: %v", err)
	}
	fmt.Printf("Deleted directory: %s\n", dirPath)
}

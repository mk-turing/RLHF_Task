package main

import (
	"fmt"
	"os"
)

type Reader interface {
	Read(p []byte) (n int, err error)
}

// Existing implementation of the Reader interface
type fileReader struct {
	file *os.File
}

func (f *fileReader) Read(p []byte) (n int, err error) {
	return f.file.Read(p)
}

// Refactored implementation of the Reader interface with defer statements
type refactoredFileReader struct {
	file *os.File
}

func (f *refactoredFileReader) Read(p []byte) (n int, err error) {
	// Defer the file closure to ensure it's closed after reading
	defer f.file.Close()

	return f.file.Read(p)
}

func main() {
	// Open a file for reading
	file, err := os.Open("example.txt")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}

	// Create a reader instance using the existing implementation
	existingReader := &fileReader{file: file}

	// Read data from the file using the existing reader
	data := make([]byte, 1024)
	n, err := existingReader.Read(data)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	fmt.Println("Read", n, "bytes from the file using the existing reader:")
	fmt.Println(string(data[:n]))

	// Create a reader instance using the refactored implementation
	refactoredReader := &refactoredFileReader{file: file}

	// Read data from the file using the refactored reader
	data = make([]byte, 1024)
	n, err = refactoredReader.Read(data)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	fmt.Println("Read", n, "bytes from the file using the refactored reader:")
	fmt.Println(string(data[:n]))

	// The file will be automatically closed by the defer statement in the refactored reader
}

package main

import (
	"fmt"
	"os"
)

type Reader interface {
	Read(p []byte) (n int, err error)
	Close() error
}

// Existing implementation of the Reader interface
type fileReader struct {
	file *os.File
}

func (f *fileReader) Read(p []byte) (n int, err error) {
	return f.file.Read(p)
}

func (f *fileReader) Close() error {
	return f.file.Close()
}

// Refactored implementation of the Reader interface with defer statements
type refactoredFileReader struct {
	file     *os.File
	tempFile *os.File
}

func (f *refactoredFileReader) Read(p []byte) (n int, err error) {
	// Defer the closing of both files
	defer f.close()
	return f.file.Read(p)
}

func (f *refactoredFileReader) Close() error {
	return f.file.Close()
}

func (f *refactoredFileReader) close() {
	if err := f.file.Close(); err != nil {
		fmt.Println("Error closing file:", err)
	}
	if f.tempFile != nil {
		if err := f.tempFile.Close(); err != nil {
			fmt.Println("Error closing temp file:", err)
		}
	}
}

func main() {
	// Open a file for reading
	file, err := os.Open("example.txt")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}

	// Create a reader instance using the refactored implementation
	refactoredReader := &refactoredFileReader{file: file}

	// Read data from the file using the refactored reader
	data := make([]byte, 1024)
	n, err := refactoredReader.Read(data)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	fmt.Println("Read", n, "bytes from the file using the refactored reader:")
	fmt.Println(string(data[:n]))

	// Demonstrating resource error handling
	// For example, we might open another temp file and encounter an error
	tempFile, err := os.CreateTemp("", "example-tmp")
	if err != nil {
		fmt.Println("Error creating temp file:", err)
		return
	}
	defer tempFile.Close() // Ensure temp file is closed

	refactoredReader.tempFile = tempFile

	// Simulating an error in the close method
	if _, err := refactoredReader.Read(data); err != nil {
		fmt.Println("Error reading file:", err)
	}

	// The close method is called automatically, handling the error
}

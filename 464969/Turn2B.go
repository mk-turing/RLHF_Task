package main

import (
	"fmt"
	"os"
	"strings"
)

type Reader interface {
	Read(p []byte) (n int, err error)
	Close() error
}

// A custom struct that represents a file reader with a buffer
type BufferedFileReader struct {
	file   *os.File
	buffer *strings.Builder
}

// NewBufferedFileReader creates a new BufferedFileReader instance
func NewBufferedFileReader(filename string) (*BufferedFileReader, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	b := strings.Builder{}
	return &BufferedFileReader{file: f, buffer: &b}, nil
}

func (f *BufferedFileReader) Read(p []byte) (n int, err error) {
	// Implements the io.Reader interface
	// For simplicity, we'll just read data from the file into the buffer.
	n, err = f.file.Read(p)
	if err != nil {
		return
	}
	_, err = f.buffer.Write(p[:n])
	return n, err
}

func (f *BufferedFileReader) Close() error {
	// Close the file first
	err := f.file.Close()
	if err != nil {
		return fmt.Errorf("error closing file: %w", err)
	}

	// Deferring the close of the buffer would be wrong
	// as we need to ensure any potential buffer error is
	// handled after the file is closed.

	// Check for buffer errors
	if err := f.Close(); err != nil {
		return fmt.Errorf("error closing buffer: %w", err)
	}
	return nil
}

func main() {
	filename := "example.txt"

	// Create a reader using the BufferedFileReader
	reader, err := NewBufferedFileReader(filename)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}

	// Properly close the reader at the end
	defer reader.Close()

	data := make([]byte, 1024)
	n, err := reader.Read(data)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	fmt.Println("Read", n, "bytes from the file:")
	fmt.Println(reader.buffer.String())
}

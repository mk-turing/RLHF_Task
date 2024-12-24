package main

import (
	"fmt"
	"io"
	"os"
)

// FileManager interface for file operations
type FileManager interface {
	Open(name string) (io.ReadWriter, error)
	Close(io.ReadWriter) error
	Write(io.ReadWriter, []byte) error
}

// Basic implementation of FileManager
type BasicFileManager struct{}

func (f *BasicFileManager) Open(name string) (io.ReadWriter, error) {
	// Open the file with read and write permissions; create it if it doesn't exist
	file, err := os.OpenFile(name, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func (f *BasicFileManager) Close(rw io.ReadWriter) error {
	file, ok := rw.(*os.File)
	if !ok {
		return fmt.Errorf("cannot close non-os.File")
	}
	return file.Close()
}

func (f *BasicFileManager) Write(rw io.ReadWriter, data []byte) error {
	_, err := rw.Write(data)
	return err
}

func main() {
	fm := &BasicFileManager{}

	// Open the file
	file, err := fm.Open("example.txt")
	if err != nil {
		panic(err)
	}

	// Ensure the file is closed
	defer func() {
		if err := fm.Close(file); err != nil {
			fmt.Printf("Error closing file: %v\n", err)
		}
	}()

	// Write data to the file
	data := []byte("Hello, World!")
	if err := fm.Write(file, data); err != nil {
		panic(err)
	}

	fmt.Println("File written successfully!")
}

package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

// File represents a file in the virtual file system
type File struct {
	data []byte
}

// Dir represents a directory in the virtual file system
type Dir struct {
	files map[string]*File
	dirs  map[string]*Dir
}

// VirtualFileSystem represents the entire virtual file system
type VirtualFileSystem struct {
	root *Dir
}

// NewVirtualFileSystem creates a new virtual file system with an empty root directory
func NewVirtualFileSystem() *VirtualFileSystem {
	return &VirtualFileSystem{root: &Dir{files: map[string]*File{}, dirs: map[string]*Dir{}}}
}

// OpenFile opens a file for reading, writing, or appending
func (vfs *VirtualFileSystem) OpenFile(path string, mode int, perm os.FileMode) (io.ReadWriteCloser, error) {
	parts := splitPath(path)
	dir := vfs.root
	for _, part := range parts[:-1] {
		if d, ok := dir.dirs[part]; !ok {
			return nil, os.ErrNotExist
		}
		dir = d
	}
	fileName := parts[-1]

	file, ok := dir.files[fileName]
	if !ok {
		if mode != os.O_CREATE {
			return nil, os.ErrNotExist
		}
		file = &File{data: []byte{}}
		dir.files[fileName] = file
	}

	return file, nil
}

// CreateFile creates a new file
func (vfs *VirtualFileSystem) CreateFile(path string, perm os.FileMode) (io.WriteCloser, error) {
	file, err := vfs.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, perm)
	if err != nil {
		return nil, err
	}
	return file, nil
}

// ReadFile reads the content of a file
func (vfs *VirtualFileSystem) ReadFile(path string) ([]byte, error) {
	file, err := vfs.OpenFile(path, os.O_RDONLY, 0644)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// WriteFile writes data to a file
func (vfs *VirtualFileSystem) WriteFile(path string, data []byte, perm os.FileMode) error {
	file, err := vfs.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, perm)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	_, err = writer.Write(data)
	if err != nil {
		return err
	}
	if err := writer.Flush(); err != nil {
		return err
	}
	return nil
}

// RemoveFile removes a file
func (vfs *VirtualFileSystem) RemoveFile(path string) error {
	parts := splitPath(path)
	dir := vfs.root
	for _, part := range parts[:-1] {
		if d, ok := dir.dirs[part]; !ok {
			return os.ErrNotExist
		}
		dir = d
	}
	fileName := parts[-1]

	if _, ok := dir.files[fileName]; !ok {
		return os.ErrNotExist
	}
	delete(dir.files, fileName)
	return nil
}

func splitPath(path string) []string {
	return strings.Split(path, "/")
}

func main() {
	vfs := NewVirtualFileSystem()

	// Create a directory
	dirPath := "docs"
	if err := vfs.MkdirAll(dirPath, 0755); err != nil {
		panic(err)
	}
	fmt.Printf("Created directory at: %s\n", dirPath)

	// Create a file inside the directory
	filePath := dirPath + "/myfile.txt"
	data := "Hello, this is a mock file!"
	if err := vfs.WriteFile(filePath, []byte(data), 0644); err != nil {
		panic(err)
	}
	fmt.Printf("Created file at: %s\n", filePath)

	// Read the file
	contents, err := vfs.ReadFile(filePath)
	if err != nil {
		panic(err)
	}
	fmt.Printf("File contents: %s\n", string(contents))

	// Append to the file
	appendedData := " This is an additional line."
	if err := vfs.WriteFile(filePath, []byte(appendedData), 0644); err != nil {
		panic(err)
	}
	fmt.Printf("Appended to file at: %s\n", filePath)

	// Read the file again
	contents, err = vfs.ReadFile(filePath)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Updated file contents: %s\n", string(contents))

	// Remove the file
	if err := vfs.RemoveFile(filePath); err != nil {
		panic(err)
	}
	fmt.Printf("Deleted file at: %s\n", filePath)
}

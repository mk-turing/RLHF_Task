package main
import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)
// MockDir represents a single directory in the mock file system.
type MockDir struct {
	name     string
	children map[string]*MockDir
}
// NewMockDir creates a new MockDir with the given name.
func NewMockDir(name string) *MockDir {
	return &MockDir{
		name:     name,
		children: make(map[string]*MockDir),
	}
}
// CreateDir creates a new subdirectory in the current directory.
func (d *MockDir) CreateDir(name string) error {
	if _, ok := d.children[name]; ok {
		return fmt.Errorf("directory '%s' already exists", name)
	}
	d.children[name] = NewMockDir(name)
	return nil
}
// ReadDir reads the names of the subdirectories in the current directory.
func (d *MockDir) ReadDir() ([]string, error) {
	names := make([]string, 0, len(d.children))
	for name := range d.children {
		names = append(names, name)
	}
	return names, nil
}
// DeleteDir deletes the specified subdirectory from the current directory.
func (d *MockDir) DeleteDir(name string) error {
	if _, ok := d.children[name]; !ok {
		return fmt.Errorf("directory '%s' not found", name)
	}
	delete(d.children, name)
	return nil
}
// MockFS represents the entire mock file system.
type MockFS struct {
	root *MockDir
}
// NewMockFS creates a new MockFS with a root directory named "/".
func NewMockFS() *MockFS {
	return &MockFS{
		root: NewMockDir("/"),
	}
}
// Chdir changes the current working directory in the mock file system.
func (fs *MockFS) Chdir(dir string) error {
	current := fs.root
	for _, part := range strings.Split(dir, "/") {
		if part == "" {
			continue
		}
		if child, ok := current.children[part]; ok {
			current = child
		} else {
			return fmt.Errorf("directory '%s' not found", part)
		}
	}
	return nil
}
// MkdirAll creates a directory structure if it does not exist.
func (fs *MockFS) MkdirAll(dir string) error {
	current := fs.root
	for _, part := range strings.Split(dir, "/") {
		if part == "" {
			continue
		}
		if child, ok := current.children[part]; ok {
			current = child
		} else {
			if err := current.CreateDir(part); err != nil {
				return err
			}
			current = current.children[part]
		}
	}
	return nil
}
// RmdirAll removes a directory structure and its contents.
func (fs *MockFS) RmdirAll(dir string) error {
	parent := fs.root
	for _, part := range strings.Split(dir, "/") {
		if part == "" {
			continue
		}
		if child, ok := parent.children[part]; ok {
			parent = child
		} else {
			return fmt.Errorf("directory '%s' not found", part)
		}
	}
	// Remove the directory and its children
	delete(parent.children, parent.name)
	return nil
}

func main() {
	fs := NewMockFS()

	// Create a directory structure
	fs.MkdirAll("/usr/local/bin")
	fs.MkdirAll("/var/log")

	// Change the current working directory
	fs.Chdir("/var/log")
package main

import (
	"io/ioutil"
	"os"
	"sync"
	"syscall"

	"golang.org/x/sys/unix"
)

// fileCache represents a cache using file-based storage with advanced techniques
type fileCache struct {
	cacheDir string
	mutex     sync.Mutex
}

// newFileCache creates a new fileCache instance and initialises the cache directory
func newFileCache(cacheDir string) (*fileCache, error) {
	// Create the cache directory if it doesn't exist
	err := os.MkdirAll(cacheDir, 0755)
	if err != nil {
		return nil, err
	}

	return &fileCache{cacheDir: cacheDir}, nil
}

// get retrieves the value associated with the specified key from the cache
func (c *fileCache) get(key string) ([]byte, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	// Open the file corresponding to the key
	file, err := os.OpenFile(c.cacheDir+"/"+key, os.O_RDONLY, 0644)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil // Key not found in the cache
		}
		return nil, err
	}
	defer file.Close()

	// Mmap the file into memory
	data, err := syscall.Mmap(int(file.Fd()), 0, 0, syscall.PROT_READ, syscall.MAP_SHARED)
	if err != nil {
		return nil, err
	}
	defer syscall.Munmap(data)

	return data, nil
}

// set stores the specified key-value pair in the cache
func (c *fileCache) set(key string, value []byte) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	// Create a temporary file for writing
	tempFile, err := ioutil.TempFile(c.cacheDir, key+".tmp.")
	if err != nil {
		return err
	}
	defer os.Remove(tempFile.Name()) // Delete the temporary file after use

	// Write the value to the temporary file
	_, err = tempFile.Write(value)
	if err != nil {
		return err
	}

	// Flush and sync the temporary file to ensure data is written to disk
	if err := tempFile.Sync(); err != nil {
		return err
	}

	// Close the temporary file
	if err := tempFile.Close(); err != nil {
		return err
	}

	// Rename the temporary file to the desired cache file name
	cacheFileName := c.cacheDir + "/" + key
	if err := os.Rename(tempFile.Name(), cacheFileName); err != nil {
		return err
	}

	// Lock the cache file to prevent other processes from accessing it while it's being updated
	if err := unix.Flock(int(tempFile.Fd()), unix.LOCK_EX|unix.LOCK_NB); err != nil {
		return err
	}
	defer unix.Flock(int(tempFile.Fd()), unix.LOCK_UN)

	// Mmap the cache file into memory
	data, err := syscall.Mmap(int(tempFile.Fd()), 0, 0, syscall.PROT_READ, syscall.MAP_SHARED)
	if err != nil {
		return err
	}
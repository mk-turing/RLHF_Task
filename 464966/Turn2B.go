package main

import (
	"fmt"
	"os"
	"sync"
	"time"
)

type Config struct {
	mu     sync.RWMutex
	data   map[string]string
	file   *os.File
	reload chan struct{}
}

func NewConfig(filePath string) *Config {
	c := &Config{
		reload: make(chan struct{}),
	}
	go c.load(filePath)
	return c
}

func (c *Config) load(filePath string) {
	for {
		select {
		case <-c.reload:
			// Close the existing file handle if present.
			if c.file != nil {
				c.file.Close()
			}

			// Open the configuration file.
			file, err := os.Open(filePath)
			if err != nil {
				fmt.Println("Error opening configuration file:", err)
				continue
			}

			// Close the file when the function returns.
			defer file.Close()

			// Read the configuration file and populate the config map.
			// (Implementation details omitted for brevity)
			newConfig := make(map[string]string)

			// Update the config data with the new configuration.
			c.mu.Lock()
			c.data = newConfig
			c.file = file
			c.mu.Unlock()

			fmt.Println("Configuration reloaded successfully.")
		}
	}
}

func (c *Config) Get(key string) (string, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	val, ok := c.data[key]
	return val, ok
}

func (c *Config) Reload() {
	select {
	case c.reload <- struct{}{}:
		fmt.Println("Configuration reload requested.")
	default:
		fmt.Println("Configuration reload already in progress.")
	}
}

func (c *Config) Close() {
	close(c.reload)
	c.file.Close()
}

func main() {
	config := NewConfig("config.ini")
	defer config.Close()

	// Use the loaded configuration.
	for {
		time.Sleep(time.Second)
		key := "some_key"
		if val, ok := config.Get(key); ok {
			fmt.Printf("%s = %s\n", key, val)
		} else {
			fmt.Printf("%s not found in configuration.\n", key)
		}

		// Request a configuration reload.
		config.Reload()
	}
}

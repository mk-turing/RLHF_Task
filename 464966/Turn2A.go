package main

import (
	"fmt"
	"os"
	"time"
)

// LoadConfig loads configuration from a file.
func LoadConfig(filePath string, timeout time.Duration) (map[string]string, error) {
	errChan := make(chan error)

	// Start a goroutine to load the configuration with a timeout.
	go func() {
		file, err := os.Open(filePath)
		if err != nil {
			errChan <- err
			return
		}

		defer func() {
			err := file.Close()
			if err != nil {
				select {
				case errChan <- err:
				default:
					// Do nothing if an error has already been reported.
				}
			}
		}()

		// Simulate the work of loading configuration.
		// (Implementation details omitted for brevity)

		// Populate the config map with loaded values.
		// (Implementation details omitted for brevity)

		errChan <- nil
	}()

	select {
	case <-errChan:
		return nil, <-errChan
	case <-time.After(timeout):
		return nil, fmt.Errorf("timeout reached after %s", timeout)
	}
}

func main() {
	filePath := "config.ini"
	timeout := time.Second * 5 // 5-second timeout

	// Example: Initial load of configuration.
	config, err := LoadConfig(filePath, timeout)
	if err != nil {
		fmt.Println("Error loading configuration:", err)
		return
	}

	// Use the loaded configuration.
	fmt.Println("Configuration loaded successfully:")
	for key, value := range config {
		fmt.Printf("%s = %s\n", key, value)
	}

	// Example: Simulating a timeout scenario when reloading.
	config, err = LoadConfig(filePath, time.Millisecond*10) // Simulate a short timeout
	if err != nil {
		fmt.Println("Error reloading configuration:", err)
		return
	}
}

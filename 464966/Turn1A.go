package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

func loadConfig(filename string) (map[string]int, error) {
	// Open the configuration file
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %v", err)
	}

	// Defer the closing of the file to ensure it's closed regardless of the function's return
	defer file.Close()

	// Read the entire content of the file
	data, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %v", err)
	}

	// Parse the file content into a map of string to int
	config := make(map[string]int)
	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		if line == "" {
			continue // skip empty lines
		}
		parts := strings.Split(line, "=")
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid config line: %s", line)
		}
		key, value := parts[0], parts[1]
		v, err := strconv.Atoi(value)
		if err != nil {
			return nil, fmt.Errorf("failed to parse value %q to int: %v", value, err)
		}
		config[key] = v
	}

	return config, nil
}

func main() {
	configFile := "config.ini"
	config, err := loadConfig(configFile)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Loaded config:", config)
}

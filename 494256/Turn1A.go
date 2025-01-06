package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

func main() {
	// Create a directory
	dirPath := "docs"
	err := os.MkdirAll(dirPath, 0755)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Created directory at: %s\n", dirPath)

	// Create a file inside the directory
	filePath := dirPath + "/myfile.txt"
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	if _, err := writer.WriteString("Hello, this is a mock file!\n"); err != nil {
		panic(err)
	}
	if err := writer.Flush(); err != nil {
		panic(err)
	}

	fmt.Printf("Created file at: %s\n", filePath)

	// Read the file
	file, err = os.OpenFile(filePath, os.O_RDONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	contents, err := io.ReadAll(reader)
	if err != nil {
		panic(err)
	}
	fmt.Printf("File contents: %s\n", string(contents))

	// Delete the directory and its contents
	err = os.RemoveAll(dirPath)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Deleted directory at: %s\n", dirPath)
}

package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func main() {
	if len(os.Args) != 2 {
		usage()
		return
	}

	apiURL := os.Args[1]
	outputDir := "output_data"

	err := os.MkdirAll(outputDir, 0755) // Ensure output directory exists
	if err != nil {
		log.Fatalf("Failed to create output directory: %v", err)
	}

	defer cleanUp() // Clean up all resources

	data, err := fetchData(apiURL)
	if err != nil {
		log.Fatalf("Failed to fetch data: %v", err)
	}

	err = processDataAndSave(outputDir, data)
	if err != nil {
		log.Fatalf("Failed to process and save data: %v", err)
	}

	fmt.Println("Data successfully processed and saved.")
}

func cleanUp() {
	log.Println("Exiting and cleaning up resources...")
	// Close any open files, network connections, etc.
}

func fetchData(url string) ([]byte, error) {
	client := http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status: %d", resp.StatusCode)
	}

	return ioutil.ReadAll(resp.Body)
}

func processDataAndSave(outputDir string, data []byte) error {
	// Process data (e.g., JSON decoding, transformation) here
	// For simplicity, let's just write the raw data to a file.

	outputFile := filepath.Join(outputDir, "data.json")
	f, err := os.Create(outputFile)
	if err != nil {
		return err
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	defer w.Flush()

	_, err = w.Write(data)
	if err != nil {
		return err
	}

	return nil
}

func usage() {
	fmt.Println("Usage:", os.Args[0], "<API URL>")
}

package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		usage()
		return
	}

	url := os.Args[1]
	outputFile := "output.json"

	defer cleanUp() // Ensure cleanup happens regardless of the outcome

	resp, err := fetchData(url)
	if err != nil {
		log.Fatalf("Failed to fetch data: %v", err)
	}
	defer resp.Body.Close() // Ensure response body is closed

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Failed to read response: %v", err)
	}

	output, err := writeData(outputFile, data)
	if err != nil {
		log.Fatalf("Failed to write data to file: %v", err)
	}

	fmt.Println("Data written to", output)
}

func cleanUp() {
	log.Println("Exiting and cleaning up resources")
	// Close files, connections, etc., as necessary
}

func fetchData(url string) (*http.Response, error) {
	client := http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	return client.Do(req)
}

func writeData(file string, data []byte) (string, error) {
	f, err := os.Create(file)
	if err != nil {
		return "", err
	}
	defer f.Close() // Ensure file handle is closed

	_, err = f.Write(data)
	if err != nil {
		return "", err
	}

	return file, nil
}

func usage() {
	fmt.Println("Usage:", os.Args[0], "<API URL>")
}

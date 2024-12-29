package main

import (
	"context"
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

	ctx := context.Background()
	outputFile := "output.json"

	defer cleanUp(ctx) // Ensure cleanup happens regardless of the outcome

	resp, err := fetchData(ctx, url)
	if err != nil {
		log.Fatalf("Failed to fetch data: %v", err)
	}
	defer resp.Body.Close() // Ensure response body is closed

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Failed to read response: %v", err)
	}

	err = writeData(ctx, outputFile, data)
	if err != nil {
		log.Fatalf("Failed to write data to file: %v", err)
	}

	fmt.Println("Data written to", outputFile)
}

func cleanUp(ctx context.Context) {
	log.Println("Exiting and cleaning up resources")
	// Add any additional cleanup logic here as necessary
}

func fetchData(ctx context.Context, url string) (*http.Response, error) {
	client := http.Client{}
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	return client.Do(req)
}

func writeData(ctx context.Context, file string, data []byte) error {
	f, err := os.Create(file)
	if err != nil {
		return err
	}
	defer f.Close() // Ensure file handle is closed

	_, err = f.Write(data)
	if err != nil {
		return err
	}

	return nil
}

func usage() {
	fmt.Println("Usage:", os.Args[0], "<API URL>")
}

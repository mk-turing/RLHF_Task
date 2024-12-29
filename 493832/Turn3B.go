package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sync"
)

const (
	numWorkers = 10 // Number of concurrent workers to fetch data
)

func main() {
	if len(os.Args) < 2 {
		usage()
		return
	}

	urls := os.Args[1:]
	outputFile := "output.json"

	ctx := context.Background()
	var wg sync.WaitGroup
	wg.Add(len(urls))

	// Channel to receive fetched data
	dataChan := make(chan []byte, len(urls))

	defer cleanUp(ctx)

	// Start worker goroutines
	for i := 0; i < numWorkers; i++ {
		go worker(ctx, &wg, dataChan)
	}

	// Send URLs to worker goroutines
	for _, url := range urls {
		go fetchData(ctx, url, dataChan)
	}

	wg.Wait() // Wait for all workers to complete
	close(dataChan)

	err := writeData(ctx, outputFile, <-dataChan) // Assume only one result is needed
	if err != nil {
		log.Fatalf("Failed to write data to file: %v", err)
	}

	fmt.Println("Data written to", outputFile)
}

func cleanUp(ctx context.Context) {
	log.Println("Exiting and cleaning up resources")
}

func worker(ctx context.Context, wg *sync.WaitGroup, dataChan chan []byte) {
	defer wg.Done()
	i := 0
	for data := range dataChan {
		defer wg.Done()
		i += 1
		err := writeData(ctx, "output_worker_"+fmt.Sprintf("%d", i)+".json", data)
		if err != nil {
			log.Printf("Failed to write data: %v", err)
		}
	}
}

func fetchData(ctx context.Context, url string, dataChan chan []byte) {
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("Failed to fetch data from %s: %v", url, err)
		return
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Failed to read response from %s: %v", url, err)
		return
	}

	select {
	case <-ctx.Done():
		return // Cancel if context is canceled
	case dataChan <- data:
	}
}

func writeData(ctx context.Context, file string, data []byte) error {
	f, err := os.Create(file)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write(data)
	if err != nil {
		return err
	}

	return nil
}

func usage() {
	fmt.Println("Usage:", os.Args[0], "<API URL1> <API URL2> ...")
}

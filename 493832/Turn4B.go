package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

const (
	numGoroutines = 10 // Number of concurrent goroutines
	outputDir     = "outputs"
	// Tuning these values may help in handling contention
	maxBodySize = 1 << 20 // 1 MB limit for response body
	maxParallel = 50      // Maximum number of concurrent API requests
)

// requestMetrics stores metrics related to an individual HTTP request
type requestMetrics struct {
	time     time.Duration
	bodySize int64
	err      error
}

func main() {
	if len(os.Args) != 2 {
		usage()
		return
	}

	urls := os.Args[1:]
	ctx := context.Background()

	if err := os.MkdirAll(outputDir, 0755); err != nil {
		log.Fatalf("Failed to create output directory: %v", err)
	}

	// Setup logging and metrics collection
	logFile, err := os.Create("tool.log")
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}
	defer logFile.Close()
	log.SetOutput(logFile)
	log.Println("Tool started...")

	var workQueue = make(chan string, len(urls)) // Bounded work queue

	// Set maximum number of concurrent API requests
	apiSem := make(chan struct{}, maxParallel)

	var wg sync.WaitGroup

	// Initialize metrics
	var numRequests uint64
	var totalDuration time.Duration
	var largestResponseSize int64
	var errors = make(map[string]int)

	wg.Add(numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go worker(ctx, workQueue, apiSem, &wg, &numRequests, &totalDuration, &largestResponseSize, &errors)
	}

	defer cleanUp(ctx, &wg, errors) // Ensure cleanup happens

	// Start timer
	start := time.Now()

	// Add all URLs to the workQueue
	for _, url := range urls {
		workQueue <- url
	}
	close(workQueue) // Close the channel once all URLs are added

	wg.Wait() // Wait for all workers to finish
	end := time.Now()

	// Print overall metrics
	log.Println("-- Overall Metrics --")
	log.Printf("Total requests processed: %d\n", numRequests)
	log.Printf("Total execution time: %s\n", end.Sub(start))
	log.Printf("Average request time: %s\n", time.Duration(totalDuration)/time.Duration(numRequests))
	log.Printf("Largest response size: %d bytes\n", largestResponseSize)

	// Print detailed error metrics
	log.Println("-- Error Metrics --")
	for err, count := range errors {
		log.Printf("Error: %s, Count: %d\n", err, count)
	}
}

func cleanUp(ctx context.Context, wg *sync.WaitGroup) {
	log.Println("Exiting and cleaning up resources")
	wg.Wait()
}

func worker(ctx context.Context, workQueue <-chan string, wg *sync.WaitGroup) {
	defer wg.Done()

	for url := range workQueue {
		if ctx.Err() != nil {
			return
		}

		log.Printf("Fetching and writing data from: %s\n", url)

		resp, err := fetchData(ctx, url)
		if err != nil {
			log.Printf("Failed to fetch data from %s: %v\n", url, err)
			continue
		}
		defer resp.Body.Close() // Ensure response body is closed

		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Printf("Failed to read response from %s: %v\n", url, err)
			continue
		}

		err = writeData(ctx, outputDir, data)
		if err != nil {
			log.Printf("Failed to write data from %s: %v\n", url, err)
			continue
		}
	}
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

func writeData(ctx context.Context, outputDir string, data []byte) error {
	const fileName = "output.json"
	filePath := fmt.Sprintf("%s/%s", outputDir, fileName)

	f, err := os.Create(filePath)
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
	fmt.Println("Usage:", os.Args[0], "<API URL> [<API URL> ...]")
}

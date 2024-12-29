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

	"github.com/prometheus/client_golang/prometheus"
)

const (
	numGoroutines = 5 // Number of concurrent goroutines
	outputDir     = "outputs"
)

var (
	workerRequestsCounter = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "worker_requests_total",
		Help: "Total number of requests processed by workers.",
	})

	workerErrorsCounter = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "worker_errors_total",
		Help: "Total number of errors encountered by workers.",
	})

	responseTimeHistogram = prometheus.NewHistogram(prometheus.HistogramOpts{
		Name:    "response_time_seconds",
		Help:    "Response time histogram in seconds.",
		Buckets: prometheus.DefBuckets,
	})
)

func main() {
	prometheus.MustRegister(workerRequestsCounter)
	prometheus.MustRegister(workerErrorsCounter)
	prometheus.MustRegister(responseTimeHistogram)

	if len(os.Args) != 2 {
		usage()
		return
	}

	urls := os.Args[1:]
	ctx := context.Background()

	if err := os.MkdirAll(outputDir, 0755); err != nil {
		log.Fatalf("Failed to create output directory: %v", err)
	}

	workQueue := make(chan string, len(urls)) // Bounded work queue

	// First, add all URLs to the workQueue
	for _, url := range urls {
		workQueue <- url
	}
	close(workQueue) // Close the channel once all URLs are added

	var wg sync.WaitGroup

	// Start the workers after adding URLs to the channel
	wg.Add(numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go worker(ctx, workQueue, &wg)
	}

	defer cleanUp(ctx, &wg) // Ensure cleanup happens regardless of the outcome

	wg.Wait() // Wait for all workers to finish
	fmt.Println("All operations completed")
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

		startTime := time.Now()
		resp, err := fetchData(ctx, url)
		if err != nil {
			log.Printf("Failed to fetch data from %s: %v\n", url, err)
			workerErrorsCounter.Inc()
			continue
		}
		responseTimeHistogram.Observe(time.Since(startTime).Seconds())

		defer resp.Body.Close() // Ensure response body is closed

		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Printf("Failed to read response from %s: %v\n", url, err)
			workerErrorsCounter.Inc()
			continue
		}

		err = writeData(ctx, outputDir, data)
		if err != nil {
			log.Printf("Failed to write data from %s: %v\n", url, err)
			workerErrorsCounter.Inc()
			continue
		}

		workerRequestsCounter.Inc()
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

package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

var url = "https://jsonplaceholder.typicode.com/todos/1"

// Shared data structure
var results []string
var mu sync.Mutex // Mutex to guard access to the shared results slice

// Function to make a request and store response
func makeRequest(wg *sync.WaitGroup, ch chan<- time.Duration) {
	defer wg.Done()
	start := time.Now()

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error making request:", err)
		ch <- 0 // Send zero duration in case of an error
		return
	}
	defer resp.Body.Close()

	// Lock the mutex to modify the shared data structure safely
	mu.Lock()
	results = append(results, fmt.Sprintf("Response status: %s", resp.Status))
	mu.Unlock()

	// Measure the duration of the request
	duration := time.Since(start)
	ch <- duration
}

// Benchmark function to simulate a high-concurrency scenario
func benchmarkGoroutines(numGoroutines int) {
	var wg sync.WaitGroup
	ch := make(chan time.Duration, numGoroutines)

	startTime := time.Now()

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go makeRequest(&wg, ch)
	}

	wg.Wait() // Wait for all goroutines to finish

	close(ch)

	var totalDuration time.Duration
	var totalRequests int

	for duration := range ch {
		if duration > 0 {
			totalDuration += duration
			totalRequests++
		}
	}

	averageDuration := totalDuration / time.Duration(totalRequests)
	totalTime := time.Since(startTime)

	fmt.Printf("Goroutines: %d, Total Time: %v, Total Requests: %d, Average Latency: %v, Throughput: %.2f RPS\n",
		numGoroutines, totalTime, totalRequests, averageDuration, float64(totalRequests)/totalTime.Seconds())
}

func main() {
	results = make([]string, 0) // Initialize the shared results slice
	fmt.Println("Running benchmarks with varying number of goroutines:")

	for numGoroutines := 10; numGoroutines <= 100; numGoroutines += 10 {
		benchmarkGoroutines(numGoroutines)
	}
}

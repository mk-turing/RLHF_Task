package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

const url = "https://jsonplaceholder.typicode.com/todos/1"

// makeRequest simulates a single HTTP request and returns its duration
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

	// Measure the duration of the request
	duration := time.Since(start)
	ch <- duration
}

func main() {
	// Number of concurrent requests (can be set to 100, 200, etc., for testing)
	var maxGoroutines int
	fmt.Print("Enter the number of concurrent requests: ")
	fmt.Scanf("%d", &maxGoroutines)

	// Simulating a higher number of total requests
	const totalRequests = 500

	var wg sync.WaitGroup
	ch := make(chan time.Duration, totalRequests) // Buffer sizes to total requests

	start := time.Now()

	// Generate totalRequests with a limit on concurrent goroutines
	for i := 0; i < totalRequests; i++ {
		wg.Add(1)

		// Use a semaphore-like pattern to limit the number of concurrent goroutines
		if i%maxGoroutines == 0 && i != 0 {
			wg.Wait() // Wait for all current goroutines to finish before launching more
		}

		go makeRequest(&wg, ch)
	}

	wg.Wait() // Wait for all goroutines to finish
	close(ch) // Close the channel

	var totalDuration time.Duration
	for duration := range ch {
		totalDuration += duration
	}

	averageDuration := totalDuration / totalRequests
	totalTimeTaken := time.Since(start)

	fmt.Printf("Total Time Taken: %v\n", totalTimeTaken)
	fmt.Printf("Average Latency per Request: %v\n", averageDuration)
	fmt.Printf("Throughput: %.2f requests per second\n", float64(totalRequests)/totalTimeTaken.Seconds())
}

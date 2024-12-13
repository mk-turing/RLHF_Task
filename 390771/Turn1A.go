package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

const numRequests = 5
const url = "https://jsonplaceholder.typicode.com/todos/1"

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
	var wg sync.WaitGroup
	ch := make(chan time.Duration, numRequests) // Buffered channel to hold latencies

	start := time.Now()

	for i := 0; i < numRequests; i++ {
		wg.Add(1)
		go makeRequest(&wg, ch)
	}

	wg.Wait() // Wait for all goroutines to finish
	close(ch) // Close the channel

	var totalDuration time.Duration
	for duration := range ch {
		totalDuration += duration
	}

	averageDuration := totalDuration / numRequests
	fmt.Printf("Total Time Taken: %v\n", time.Since(start))
	fmt.Printf("Average Latency: %v\n", averageDuration)
}

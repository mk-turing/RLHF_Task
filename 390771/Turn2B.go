package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

var url = "https://jsonplaceholder.typicode.com/todos/1"

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

func benchmarkGoroutines(numGoroutines int) {
	var wg sync.WaitGroup
	ch := make(chan time.Duration, numGoroutines)

	startTime := time.Now()

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go makeRequest(&wg, ch)
	}

	wg.Wait()
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
	fmt.Println("Running benchmarks with varying number of goroutines:")

	for numGoroutines := 1; numGoroutines <= 20; numGoroutines += 5 {
		benchmarkGoroutines(numGoroutines)
	}
}

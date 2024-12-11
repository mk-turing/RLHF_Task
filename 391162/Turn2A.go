package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"
)

const (
	url                = "http://localhost:8080"
	totalRequests      = 1000 // Total number of requests
	concurrentRequests = 100  // Number of concurrent requests
)

type Stats struct {
	mu            sync.Mutex
	responseTimes []time.Duration
	errorCount    int
}

func (s *Stats) recordResponseTime(duration time.Duration) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.responseTimes = append(s.responseTimes, duration)
}

func (s *Stats) recordError() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.errorCount++
}

func benchmark() {
	var wg sync.WaitGroup
	stats := &Stats{}
	sem := make(chan struct{}, concurrentRequests)

	start := time.Now()

	for i := 0; i < totalRequests; i++ {
		wg.Add(1)
		sem <- struct{}{} // Acquire a semaphore
		go func() {
			defer wg.Done()
			defer func() { <-sem }() // Release the semaphore

			startTime := time.Now()
			resp, err := http.Get(url)
			if err != nil {
				fmt.Println("Error:", err)
				stats.recordError()
				return
			}
			_, _ = ioutil.ReadAll(resp.Body) // Read and discard the response body
			resp.Body.Close()
			stats.recordResponseTime(time.Since(startTime))
		}()
	}

	wg.Wait() // Wait for all requests to finish
	elapsed := time.Since(start)

	// Analyze results
	totalTime := float64(elapsed.Milliseconds())
	avgResponseTime := totalTime / float64(totalRequests)
	throughput := float64(totalRequests) / totalTime * 1000 // Requests per second

	fmt.Printf("Total Time: %.2f ms\n", totalTime)
	fmt.Printf("Average Response Time: %.2f ms\n", avgResponseTime)
	fmt.Printf("Throughput: %.2f requests/sec\n", throughput)
	fmt.Printf("Error Rate: %.2f%%\n", float64(stats.errorCount)/float64(totalRequests)*100)

	// Calculate response time statistics
	if len(stats.responseTimes) > 0 {
		min, max, total := stats.responseTimes[0], stats.responseTimes[0], 0
		for _, rt := range stats.responseTimes {
			if rt < min {
				min = rt
			}
			if rt > max {
				max = rt
			}
			total += int(rt.Milliseconds())
		}
		mean := total / len(stats.responseTimes)
		fmt.Printf("Min Response Time: %v\n", min)
		fmt.Printf("Max Response Time: %v\n", max)
		fmt.Printf("Mean Response Time: %v ms\n", mean)
	} else {
		fmt.Println("No successful responses to analyze.")
	}
}

func main() {
	benchmark()
}

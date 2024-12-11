package main

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"sync"
	"sync/atomic"
	"time"
)

const (
	url                   = "http://localhost:8080"
	totalRequests         = 1000 // Total number of requests
	minConcurrentRequests = 10   // Min concurrent requests for each load phase
	maxConcurrentRequests = 100  // Max concurrent requests for each load phase
)

type LoadTestStats struct {
	mu            sync.Mutex
	responseTimes []time.Duration
	errorCount    int32 // Using atomic integer for concurrency safety
}

func (s *LoadTestStats) recordResponseTime(duration time.Duration) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.responseTimes = append(s.responseTimes, duration)
}

func (s *LoadTestStats) recordError() {
	atomic.AddInt32(&s.errorCount, 1)
}

func simulateDynamicLoad(stats *LoadTestStats) {
	var wg sync.WaitGroup
	for i := 0; i < totalRequests; i++ {
		wg.Add(1)
		currentConcurrency := rand.Intn(maxConcurrentRequests-minConcurrentRequests+1) + minConcurrentRequests
		go func(concurrency int) {
			defer wg.Done()
			sem := make(chan struct{}, concurrency) // control concurrency using a semaphore

			for j := 0; j < concurrency; j++ {
				sem <- struct{}{} // Acquire semaphore
				go func() {
					defer func() { <-sem }() // Release semaphore
					startTime := time.Now()
					resp, err := http.Get(url)
					if err != nil {
						fmt.Println("Error:", err)
						stats.recordError()
						return
					}
					_, _ = ioutil.ReadAll(resp.Body) // Read response
					resp.Body.Close()
					stats.recordResponseTime(time.Since(startTime))
				}()
			}
			// wait for goroutines to finish
			for j := 0; j < concurrency; j++ {
				sem <- struct{}{} // Enforce wait if in large concurrency
			}
		}(currentConcurrency)
	}
	wg.Wait() // Wait for all requests to finish
}

func analyzeResults(stats *LoadTestStats) {
	totalResponse := len(stats.responseTimes)
	if totalResponse == 0 {
		fmt.Println("No successful responses to analyze.")
		return
	}

	minRT := stats.responseTimes[0]
	maxRT := stats.responseTimes[0]
	totalRT := int64(0)

	for _, rt := range stats.responseTimes {
		if rt < minRT {
			minRT = rt
		}
		if rt > maxRT {
			maxRT = rt
		}
		totalRT += int64(rt.Milliseconds())
	}

	//_ := totalRT / int64(totalResponse)
	errorCount := atomic.LoadInt32(&stats.errorCount)

	// Print final statistics
	fmt.Printf("Min Response Time: %v\n", minRT)
	fmt.Printf("Max Response Time: %v\n", maxRT)
	fmt.Printf("Mean Response Time: %.2f ms\n", float64(totalRT)/float64(totalResponse))
	fmt.Printf("Error Rate: %.2f%%\n", float64(errorCount)*100/float64(totalRequests))
}

func main() {
	stats := &LoadTestStats{}
	simulateDynamicLoad(stats)
	analyzeResults(stats)
}

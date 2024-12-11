package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
	"sync"
	"time"
)

const (
	url                = "http://localhost:8080"
	maxConcurrentUsers = 1000             // Maximum concurrent users
	totalRequests      = 5000             // Total number of requests
	rampUpTime         = 30 * time.Second // Ramp up time for requests
)

type Stats struct {
	mu            sync.Mutex
	responseTimes []time.Duration
	errorCount    int
	queueTimes    []time.Duration
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

func (s *Stats) recordQueueTime(duration time.Duration) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.queueTimes = append(s.queueTimes, duration)
}

func makeRequest(stats *Stats) {
	// Start queuing time measurement
	startQueueTime := time.Now()

	resp, err := http.Get(url)
	queueTime := time.Since(startQueueTime)

	stats.recordQueueTime(queueTime)

	if err != nil {
		stats.recordError()
		return
	}
	defer resp.Body.Close()

	_, _ = ioutil.ReadAll(resp.Body)
	stats.recordResponseTime(time.Since(startQueueTime))
}

func main() {
	stats := &Stats{}
	var wg sync.WaitGroup

	start := time.Now()
	currentUsers := 0

	for i := 0; i < totalRequests; i++ {
		if currentUsers < maxConcurrentUsers {
			wg.Add(1)
			currentUsers++
			go func() {
				defer wg.Done()
				makeRequest(stats)
			}()
		}

		// Wait for ramp up time
		if time.Since(start) < rampUpTime {
			time.Sleep(1 * time.Millisecond) // Adjust to control ramp-up speed
		}

		if i > 0 && i%1000 == 0 {
			fmt.Printf("Completed %d requests, Active Users: %d\n", i, currentUsers)
		}
	}

	wg.Wait() // Wait for all goroutines to finish

	elapsed := time.Since(start)
	totalRequests := len(stats.responseTimes)
	avgResponseTime := float64(elapsed.Milliseconds()) / float64(totalRequests)

	fmt.Printf("Total Time: %.2f ms\n", float64(elapsed.Milliseconds()))
	fmt.Printf("Average Response Time: %.2f ms\n", avgResponseTime)
	fmt.Printf("Error Count: %d\n", stats.errorCount)
	fmt.Printf("Queue Times: %v\n", stats.queueTimes)

	// Calculate 95th and 99th percentile response times
	calculatePercentiles(stats.responseTimes, []float64{95, 99})
}

func calculatePercentiles(times []time.Duration, percentiles []float64) {
	sortedTimes := make([]time.Duration, len(times))
	copy(sortedTimes, times)
	sort.Slice(sortedTimes, func(i, j int) bool { return sortedTimes[i] < sortedTimes[j] })

	for _, p := range percentiles {
		rank := int(p / 100.0 * float64(len(sortedTimes)))
		fmt.Printf("%dth Percentile: %.2f ms\n", int(p), float64(sortedTimes[rank].Milliseconds()))
	}
}

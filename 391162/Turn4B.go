package main

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"sort"
	"sync"
	"time"
)

const (
	url                = "http://localhost:8080"
	rampUpTime         = 30 * time.Second // Time to ramp up to max concurrency
	rampDownTime       = 30 * time.Second // Time to ramp down to zero concurrency
	maxConcurrentUsers = 100
	requestDuration    = 1 * time.Second // Duration for which each user is active
)

type Stats struct {
	mu            sync.Mutex
	responseTimes []time.Duration
	errorCount    int
	queueTimes    []time.Duration // New: Queue times before request is processed
}

func (s *Stats) recordResponseTime(startTime time.Time) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.responseTimes = append(s.responseTimes, time.Since(startTime))
}

func (s *Stats) recordError() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.errorCount++
}

func (s *Stats) recordQueueTime(startTime time.Time) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.queueTimes = append(s.queueTimes, time.Since(startTime))
}

func makeRequest() {
	//startQueueTime := time.Now()
	//defer stats.recordQueueTime(startQueueTime)

	//startTime := time.Now()
	resp, err := http.Get(url)
	if err != nil {
		//stats.recordError()
		return
	}
	defer resp.Body.Close()
	_, _ = ioutil.ReadAll(resp.Body)
	//stats.recordResponseTime(startTime)
}

func main() {
	stats := &Stats{}
	rand.Seed(time.Now().UnixNano())

	start := time.Now()
	var wg sync.WaitGroup

	currentUsers := 0
	stepUsers := 10 // Increment/decrement users by this number

	for {
		time.Sleep(1 * time.Second)

		if currentUsers == 0 && time.Since(start) < rampUpTime {
			currentUsers += stepUsers
		} else if currentUsers == maxConcurrentUsers && time.Since(start) >= rampUpTime+requestDuration {
			currentUsers -= stepUsers
			if currentUsers <= 0 {
				break
			}
		}

		fmt.Printf("Current Users: %d, Total Time: %v\n", currentUsers, time.Since(start))

		//newUsers := currentUsers - wg.Wait() // Wait for active users to complete
		//for i := 0; i < newUsers; i++ {
		//	wg.Add(1)
		//	go makeRequest()
		//}
	}

	wg.Wait()
	elapsed := time.Since(start)

	// Analyze results
	totalTime := float64(elapsed.Milliseconds())
	avgResponseTime := totalTime / float64(len(stats.responseTimes))
	throughput := float64(len(stats.responseTimes)) / totalTime * 1000 // Requests per second

	fmt.Printf("Total Time: %.2f ms\n", totalTime)
	fmt.Printf("Average Response Time: %.2f ms\n", avgResponseTime)
	fmt.Printf("Throughput: %.2f requests/sec\n", throughput)
	fmt.Printf("Error Rate: %.2f%%\n", float64(stats.errorCount)/float64(len(stats.responseTimes))*100)

	// Calculate response time and queue time percentiles
	percentiles := []float64{50, 75, 90, 95, 99}
	calculatePercentiles(stats.responseTimes, "Response Times", percentiles)
	calculatePercentiles(stats.queueTimes, "Queue Times", percentiles)
}

func calculatePercentiles(times []time.Duration, name string, percentiles []float64) {
	sortedTimes := make([]time.Duration, len(times))
	copy(sortedTimes, times)
	sort.Slice(sortedTimes, func(i, j int) bool { return sortedTimes[i] < sortedTimes[j] })

	fmt.Printf("\n%s Percentiles:\n", name)
	for _, p := range percentiles {
		rank := int(p / 100.0 * float64(len(sortedTimes)))
		fmt.Printf("%dth Percentile: %.2f ms\n", p, float64(sortedTimes[rank].Milliseconds()))
	}
}

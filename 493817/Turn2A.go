package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func main() {
	// Set the number of concurrent system calls to simulate
	numConcurrency := 10
	// Set the total number of system calls to simulate during startup
	numCalls := 1000
	// Set the maximum latency for each system call
	maxLatency := 10 * time.Millisecond

	// Seed the random number generator for latency simulation
	rand.Seed(time.Now().UnixNano())

	// Measure the startup time with concurrent system calls
	var startupTimes []time.Duration
	for i := 0; i < 5; i++ {
		var wg sync.WaitGroup
		start := time.Now()

		for j := 0; j < numConcurrency; j++ {
			wg.Add(1)
			go simulateSystemCalls(numCalls, maxLatency, &wg)
		}

		wg.Wait()
		startupTimes = append(startupTimes, time.Since(start))
	}

	// Analyze the startup times and calculate statistics
	totalStartupTime := time.Duration(0)
	for _, startupTime := range startupTimes {
		totalStartupTime += startupTime
	}
	averageStartupTime := totalStartupTime / time.Duration(len(startupTimes))
	minStartupTime := startupTimes[0]
	maxStartupTime := startupTimes[0]
	for _, startupTime := range startupTimes[1:] {
		if startupTime < minStartupTime {
			minStartupTime = startupTime
		} else if startupTime > maxStartupTime {
			maxStartupTime = startupTime
		}
	}

	// Print the analysis results
	fmt.Println("Concurrent Startup Time Analysis:")
	fmt.Printf("Average Startup Time: %v\n", averageStartupTime)
	fmt.Printf("Minimum Startup Time: %v\n", minStartupTime)
	fmt.Printf("Maximum Startup Time: %v\n", maxStartupTime)
}

func simulateSystemCalls(numCalls int, maxLatency time.Duration, wg *sync.WaitGroup) {
	defer wg.Done()

	for j := 0; j < numCalls; j++ {
		// Simulate system call latency with random duration
		latency := time.Duration(rand.Intn(int(maxLatency)))
		time.Sleep(latency)
	}
}

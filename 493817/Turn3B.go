package main

import (
	"fmt"
	"math/rand"
	"os"
	"runtime/pprof"
	"sync"
	"time"
)

func main() {
	// Set the number of concurrent system calls to simulate
	numConcurrentCalls := 10
	// Set the number of system calls per concurrent group
	numCallsPerGroup := 100
	// Set the maximum latency for each system call
	maxLatency := 10 * time.Millisecond

	// Enable CPU profiling
	f, err := os.Create("cpu.prof")
	if err != nil {
		panic(err)
	}
	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()

	// Seed the random number generator for latency simulation
	rand.Seed(time.Now().UnixNano())

	// Measure the startup time with concurrent system calls and varying latencies
	var startupTimes []time.Duration
	for i := 0; i < 5; i++ {
		start := time.Now()
		wg := sync.WaitGroup{}
		wg.Add(numConcurrentCalls)
		for j := 0; j < numConcurrentCalls; j++ {
			go func() {
				defer wg.Done()
				for k := 0; k < numCallsPerGroup; k++ {
					// Simulate system call latency with random duration
					latency := time.Duration(rand.Intn(int(maxLatency)))
					time.Sleep(latency)
				}
			}()
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

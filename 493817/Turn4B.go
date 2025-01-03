package main

import (
	"fmt"
	"math/rand"
	"os"
	"runtime/pprof"
	"sync"
	"time"
)

// Global variables to simulate external factors
var (
	cpuLoad     float64
	memoryUsage float64
)

func simulateSystemCallLatency() time.Duration {
	// Base latency
	latency := 5 * time.Millisecond

	// Random fluctuation based on CPU load
	fluctuation := time.Duration(float64(latency) * 0.2 * (1 - cpuLoad))
	latency += time.Duration(rand.Int63n(int64(fluctuation)))

	// Random fluctuation based on memory usage
	fluctuation = time.Duration(float64(latency) * 0.1 * (1 - memoryUsage))
	latency += time.Duration(rand.Int63n(int64(fluctuation)))

	// Add random noise to simulate unpredictable variations
	noise := time.Duration(rand.Int63n(int64(10 * time.Millisecond / 2)))
	latency += noise

	return latency
}

func generateRandomCPUAndMemoryUsage() {
	for {
		// Simulate random CPU load between 0.2 and 0.8
		cpuLoad = 0.2 + 0.6*rand.Float64()
		// Simulate random memory usage between 0.3 and 0.9
		memoryUsage = 0.3 + 0.6*rand.Float64()
		time.Sleep(1 * time.Second)
	}
}

func main() {
	// Initialize external factors
	cpuLoad = 0.5     // Average CPU load
	memoryUsage = 0.6 // Average memory usage

	// Set the number of concurrent system calls to simulate
	numConcurrentCalls := 10
	// Set the number of system calls per concurrent group
	numCallsPerGroup := 100

	// Seed the random number generator for latency simulation
	rand.Seed(time.Now().UnixNano())

	// Start CPU profiling
	f, err := os.Create("cpu.prof")
	if err != nil {
		fmt.Println("error creating cpu profile:", err)
		return
	}
	defer f.Close()
	if err := pprof.StartCPUProfile(f); err != nil {
		fmt.Println("error starting cpu profile:", err)
		return
	}
	defer pprof.StopCPUProfile()

	// Measure the startup time with concurrent system calls and varying latencies influenced by external factors
	var startupTimes []time.Duration
	for i := 0; i < 5; i++ {
		start := time.Now()
		wg := sync.WaitGroup{}
		wg.Add(numConcurrentCalls)
		for j := 0; j < numConcurrentCalls; j++ {
			go func() {
				defer wg.Done()
				for k := 0; k < numCallsPerGroup; k++ {
					// Simulate system call latency with random duration influenced by external factors
					latency := simulateSystemCallLatency()
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
	fmt.Println("Dynamic Startup Time Analysis:")
	fmt.Printf("Average Startup Time: %v\n", averageStartupTime)
	fmt.Printf("Minimum Startup Time: %v\n", minStartupTime)
	fmt.Printf("Maximum Startup Time: %v\n", maxStartupTime)

	// Initialize external factors
	cpuLoad = 0.5     // Average CPU load
	memoryUsage = 0.6 // Average memory usage

	// Generate random CPU and memory usage in the background
	go generateRandomCPUAndMemoryUsage()

	// Rest of the benchmarking code remains the same...
}

package main

import (
	"fmt"
	"github.com/xlab/randx"
	"math/rand"
	"os/exec"
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

	// Seed the random number generator for latency simulation
	rand.Seed(time.Now().UnixNano())

	// Set up dynamic latency distribution using randx
	dist := randx.NewNormal(float64(maxLatency/2), float64(maxLatency/4))

	// Benchmarking framework
	for cpuLoad := 0; cpuLoad <= 100; cpuLoad += 10 {
		for memUsage := 0; memUsage <= 90; memUsage += 10 {
			// Run system stress tests
			runStressTests(cpuLoad, memUsage)

			// Measure startup time with concurrent system calls and varying latencies
			var startupTimes []time.Duration
			for i := 0; i < 5; i++ {
				start := time.Now()
				wg := sync.WaitGroup{}
				wg.Add(numConcurrentCalls)
				for j := 0; j < numConcurrentCalls; j++ {
					go func() {
						defer wg.Done()
						for k := 0; k < numCallsPerGroup; k++ {
							// Simulate system call latency with dynamic distribution
							latency := time.Duration(dist.Float64() * 1000)
							time.Sleep(latency * time.Millisecond)
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
			fmt.Printf("CPU Load: %d%% Mem Usage: %d%% - Average Startup Time: %v\n", cpuLoad, memUsage, averageStartupTime)
		}
	}
}

func runStressTests(cpuLoad, memUsage int) {
	// Stress CPU
	cmd := exec.Command("stress", "-c", fmt.Sprintf("%d", cpuLoad/10))
	if err := cmd.Run(); err != nil {
		fmt.Println("Error running CPU stress test:", err)
	}

	// Stress memory
	cmd = exec.Command("stress-ng", "-m", fmt.Sprintf("%d", memUsage/10))
	if err := cmd.Run(); err != nil {
		fmt.Println("Error running memory stress test:", err)
	}
}

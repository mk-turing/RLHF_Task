package main

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"runtime/pprof"
	"sync"
	"testing"
	"time"
)

func task(ctx context.Context, wg *sync.WaitGroup, counter *int) {
	defer wg.Done()

	for i := 0; i < 1000; i++ {
		// Simulate work
		for j := 0; j < 1000; j++ {
			*counter++
		}
		// Randomize task duration
		time.Sleep(time.Duration(rand.Intn(10)) * time.Millisecond)
	}
}

// BenchmarkWorkload is a benchmark function to measure throughput for a given workload.
func BenchmarkWorkload(b *testing.B) {
	rand.Seed(time.Now().UnixNano())

	workloads := []struct {
		numWorkers   int
		totalTasks   int
		totalTimeout time.Duration
	}{
		{1, 10000, 5 * time.Second},
		{5, 10000, 5 * time.Second},
		{10, 10000, 5 * time.Second},
		{20, 20000, 10 * time.Second},
		{50, 20000, 15 * time.Second},
	}

	for _, workload := range workloads {
		b.Run(fmt.Sprintf("Workload-%d-%d-%s", workload.numWorkers, workload.totalTasks, workload.totalTimeout), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				RunWorkload(workload)
			}
		})
	}
}

// RunWorkload is a function to run a single workload
func RunWorkload(workload struct {
	numWorkers   int
	totalTasks   int
	totalTimeout time.Duration
}) {
	counter := 0
	wg := sync.WaitGroup{}
	ctx, cancel := context.WithTimeout(context.Background(), workload.totalTimeout)
	defer cancel()

	for i := 0; i < workload.numWorkers; i++ {
		wg.Add(1)
		go task(ctx, &wg, &counter)
	}
	wg.Wait()
	_ = counter // Avoid unused variable warning
	// Calculate and output metrics
}

func main() {
	// Profiling setup
	f, err := os.Create("cpu.prof")
	if err != nil {
		fmt.Printf("Error creating profile file: %v\n", err)
		os.Exit(1)
	}
	defer f.Close()
	if err := pprof.StartCPUProfile(f); err != nil {
		fmt.Printf("Error starting CPU profile: %v\n", err)
		os.Exit(1)
	}
	defer pprof.StopCPUProfile()

	// Run the benchmarking tests
	testing.Benchmark(BenchmarkWorkload)
}

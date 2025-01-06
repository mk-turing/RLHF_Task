package main

import (
	"fmt"
	"sync"
	"time"
)

type Benchmark interface {
	Run() (time.Duration, error)
}

type Test1 struct{}

func (t *Test1) Run() (time.Duration, error) {
	start := time.Now()
	// Benchmark logic here
	end := time.Now()
	return end.Sub(start), nil
}

func calculateMean(durations []time.Duration) float64 {
	sum := float64(0)
	for _, duration := range durations {
		sum += float64(duration)
	}
	return sum / float64(len(durations))
}

func aggregateResults(results map[string][]time.Duration) map[string]float64 {
	aggregated := map[string]float64{}
	for machine, durations := range results {
		mean := calculateMean(durations)
		aggregated[machine] = mean
	}
	return aggregated
}

func main() {
	workers := []string{"worker1", "worker2", "worker3"} // List of worker machine addresses
	numTests := 10
	test := &Test1{}

	workerResults := map[string][]time.Duration{}
	var mu sync.Mutex // Mutex to synchronize map access

	// Send test to workers
	for _, worker := range workers {
		fmt.Println("Sending test to", worker)
	}

	// Receive results from workers
	var wg sync.WaitGroup
	for _, worker := range workers {
		wg.Add(1)
		go func(worker string) {
			defer wg.Done()
			results := []time.Duration{}
			for i := 0; i < numTests; i++ {
				duration, _ := test.Run()
				results = append(results, duration)
			}

			mu.Lock() // Lock before modifying the map
			workerResults[worker] = results
			mu.Unlock() // Unlock after modifying the map
		}(worker)
	}

	wg.Wait()

	// Aggregate results
	aggregated := aggregateResults(workerResults)
	for machine, mean := range aggregated {
		fmt.Printf("Machine: %s, Mean duration: %f\n", machine, mean)
	}
}

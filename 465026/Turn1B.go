package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func worker(wg *sync.WaitGroup, numTasks int, latencyChan chan time.Duration, throughputChan chan int) {
	defer wg.Done()

	for i := 0; i < numTasks; i++ {
		// Simulate task execution with a random latency
		latency := time.Duration(rand.Intn(100)) * time.Millisecond
		time.Sleep(latency)

		// Send latency and throughput data to the respective channels
		latencyChan <- latency
		throughputChan <- 1
	}
}

func measureThroughputAndLatency(numWorkers int, numTasks int) (throughput float64, latency time.Duration) {
	var wg sync.WaitGroup
	latencyChan := make(chan time.Duration)
	throughputChan := make(chan int)

	// Start the workers
	wg.Add(numWorkers)
	for i := 0; i < numWorkers; i++ {
		go worker(&wg, numTasks/numWorkers, latencyChan, throughputChan)
	}

	// Wait for all workers to complete
	close(latencyChan)
	close(throughputChan)

	// Calculate throughput
	totalThroughput := 0
	for throughput := range throughputChan {
		totalThroughput += throughput
	}
	throughput = float64(totalThroughput) / float64(numTasks) * float64(time.Second)

	// Calculate latency
	totalLatency := time.Duration(0)
	for latency := range latencyChan {
		totalLatency += latency
	}
	latency = totalLatency / time.Duration(numTasks)

	return throughput, latency
}

func main() {
	numWorkers := 4
	numTasks := 10000

	throughput, latency := measureThroughputAndLatency(numWorkers, numTasks)

	fmt.Printf("Throughput: %.2f tasks/second\n", throughput)
	fmt.Printf("Latency: %s\n", latency)
}

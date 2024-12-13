package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

const (
	numRequests = 100000
	//numGoroutines = []int{1, 5, 10, 20, 50, 100, 200}
)

var url = "https://jsonplaceholder.typicode.com/todos/1"
var sharedCounter int64
var mutex sync.Mutex
var atomicCounter int64

func incrementWithMutex(wg *sync.WaitGroup) {
	defer wg.Done()
	for i := 0; i < numRequests; i++ {
		mutex.Lock()
		sharedCounter++
		mutex.Unlock()
	}
}

func incrementAtomic(wg *sync.WaitGroup) {
	defer wg.Done()
	for i := 0; i < numRequests; i++ {
		atomic.AddInt64(&atomicCounter, 1)
	}
}

func benchmarkConcurrency(numGoroutines int) {
	var wg sync.WaitGroup

	startTime := time.Now()

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go incrementWithMutex(&wg) // Uncomment the next line to use atomic increments
		// go incrementAtomic(&wg)
	}

	wg.Wait()

	totalTime := time.Since(startTime)

	fmt.Printf("Goroutines: %d, Total Time: %v, Shared Counter: %d, Atomic Counter: %d, Throughput: %.2f RPS\n",
		numGoroutines, totalTime, sharedCounter, atomicCounter, float64(numRequests*numGoroutines)/totalTime.Seconds())

	sharedCounter = 0
	atomicCounter = 0
}

func main() {
	fmt.Println("Running benchmarks with varying number of goroutines using mutex:")

	for _, numGoroutines := range []int{1, 5, 10, 20, 50, 100, 200} {
		benchmarkConcurrency(numGoroutines)
	}

	// Uncomment the following lines to run benchmarks with atomic operations
	/*
		fmt.Println("\nRunning benchmarks with varying number of goroutines using atomic operations:")

		for _, numGoroutines := range numGoroutines {
			benchmarkConcurrency(numGoroutines)
		}
	*/
}

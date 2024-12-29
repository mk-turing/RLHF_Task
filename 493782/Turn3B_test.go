package main

import (
	"errors"
	"math/rand"
	"sync"
	"testing"
	"time"
)

// Define an empty error type for simulation purposes.
var EmptyDataError = errors.New("data is empty")

// processData processes data and calls a callback on success or error.
func processData(data []byte, callback func(error)) {
	// Simulated error handling
	if data == nil {
		callback(errors.New("data is nil"))
		return
	}
	// Simulated work that might fail.
	if len(data) == 0 {
		callback(EmptyDataError)
		return
	}
	// Simulated successful processing.
	callback(nil)
}

func processDataWithConcurrency(data [][]byte, callbacks []func(error)) {
	wg := sync.WaitGroup{}
	wg.Add(len(data))

	for _, d := range data {
		go func(data []byte, callback func(error)) {
			defer wg.Done()
			// Simulated processing
			if data == nil {
				callback(errors.New("data is nil"))
				return
			}
			// Simulate error after a random delay to simulate concurrent execution
			time.Sleep(time.Duration(rand.Intn(50)) * time.Millisecond)
			if len(data) == 0 {
				callback(EmptyDataError)
				return
			}
			callback(nil)
		}(d, callbacks[0]) // Call the same callback for simplicity.
	}
	wg.Wait()
}

// BenchmarkProcessData runs a benchmark for processData.
func BenchmarkProcessData(b *testing.B) {
	data := make([]byte, 100) // 100-byte data
	for i := 0; i < b.N; i++ {
		processData(data, func(error) {}) // Ignore the error for benchmarking.
	}
}

// BenchmarkProcessDataWithConcurrency runs a benchmark for processDataWithConcurrency.
func BenchmarkProcessDataWithConcurrency(b *testing.B) {
	data := make([]byte, 100) // 100-byte data
	for i := 0; i < b.N; i++ {
		// Simulate parallel processing by calling processDataWithConcurrency with a single goroutine (single callback).
		processDataWithConcurrency([][]byte{data}, []func(error){func(error) {}})
	}
}

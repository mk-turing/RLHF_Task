package main

import (
	"fmt"
	"os"
	"sync"
	"testing"
)

var counter int
var mu sync.Mutex

// A function with a race condition
func increment() {
	for i := 0; i < 100000; i++ {
		mu.Lock()
		counter++
		mu.Unlock()
	}
}

// Test for race detection (manual inspection with -race flag)
func TestIncrementRace(t *testing.T) {
	if testing.Short() {
		t.Skip("race detector test is skipped in short mode")
	}

	counter = 0

	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(2)
		go func() {
			defer wg.Done()
			increment()
		}()
		go func() {
			defer wg.Done()
			increment()
		}()
	}

	wg.Wait()
	fmt.Println("Counter:", counter)
}

// Thread stress test to simulate high concurrency
func StressTestIncrement(t *testing.T) {
	if testing.Short() {
		t.Skip("stress test is skipped in short mode")
	}

	counter = 0

	const numGoroutines = 100
	const iterations = 100000

	var wg sync.WaitGroup
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < iterations; j++ {
				mu.Lock()
				counter++
				mu.Unlock()
			}
		}()
	}

	wg.Wait()
	fmt.Println("Counter:", counter)
}

func TestMain(m *testing.M) {
	// Running tests with race detection requires the -race flag
	fmt.Println("Running tests with potential race conditions...")
	os.Exit(m.Run())
}

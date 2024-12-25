package _65040

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func task(wg *sync.WaitGroup) {
	defer wg.Done()
	// Simulate some work being done by the goroutine
	time.Sleep(time.Millisecond * 50)
}

func TestSyncWaitGroup(t *testing.T) {
	wg := &sync.WaitGroup{}
	numGoroutines := 1000

	// Start timer before running goroutines
	start := time.Now()

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go task(wg)
	}

	// Wait for all goroutines to finish
	wg.Wait()

	// Calculate execution time
	end := time.Now()
	executionTime := end.Sub(start)

	fmt.Printf("Execution time with %d goroutines: %s\n", numGoroutines, executionTime)
}

package _65040

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestWaitGroupPerformance(t *testing.T) {
	const numGoroutines = 1000
	const iterations = 10000

	wg := &sync.WaitGroup{}
	start := time.Now()

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < iterations; j++ {
				// Simulate work
				time.Sleep(1 * time.Nanosecond)
			}
		}()
	}

	wg.Wait()
	end := time.Now()

	fmt.Printf("Time taken: %v\n", end.Sub(start))
}

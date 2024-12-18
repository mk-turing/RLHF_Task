package main

import (
	"sync"
	"testing"
)

var counter int // Shared mutable state

func incrementCounter() {
	counter++
}

func TestIncrementCounter(t *testing.T) {
	// Run the test multiple times to increase the likelihood of reproducing the race
	for i := 0; i < 1000; i++ {
		counter = 0
		var wg sync.WaitGroup
		wg.Add(2)

		go func() {
			incrementCounter()
			wg.Done()
		}()

		go func() {
			incrementCounter()
			wg.Done()
		}()

		wg.Wait()

		if counter != 2 {
			t.Errorf("Expected counter to be 2 after incrementing twice, but got: %d", counter)
		}
	}
}

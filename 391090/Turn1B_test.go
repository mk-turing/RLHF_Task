package main

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func IncrementCounter(counter *int, wg *sync.WaitGroup) {
	defer wg.Done()
	for i := 0; i < 1000; i++ {
		*counter++
	}
}

func TestConcurrentIncrement(t *testing.T) {
	counter := 0
	var wg sync.WaitGroup
	wg.Add(2)

	go IncrementCounter(&counter, &wg)
	go IncrementCounter(&counter, &wg)

	wg.Wait()

	// Expecting the counter to be 2000 after both goroutines complete
	expectedCounter := 2000
	if counter != expectedCounter {
		t.Errorf("Expected counter to be %d, got %d", expectedCounter, counter)
	}
}

func ExampleSynchronizedIncrement() {
	counter := 0
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		for i := 0; i < 1000; i++ {
			// Synchronize access to the counter using a mutex
			var mutex sync.Mutex
			mutex.Lock()
			counter++
			mutex.Unlock()
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < 1000; i++ {
			// Synchronize access to the counter using a mutex
			var mutex sync.Mutex
			mutex.Lock()
			counter++
			mutex.Unlock()
		}
	}()

	wg.Wait()

	fmt.Println("Final counter value:", counter) // Output: Final counter value: 2000
}

func TestSynchronizedIncrement(t *testing.T) {
	counter := 0
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		for i := 0; i < 1000; i++ {
			// Synchronize access to the counter using a mutex
			var mutex sync.Mutex
			mutex.Lock()
			counter++
			mutex.Unlock()
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < 1000; i++ {
package _65040
import (
	"fmt"
	"sync"
	"testing"
	"time"
)

// Global variable shared among goroutines
var counter int

// --------------------- sync.Mutex approach ----------------------
func taskMutex(wg *sync.WaitGroup, mutex *sync.Mutex) {
	defer wg.Done()
	for i := 0; i < 10000; i++ {
		mutex.Lock()
		counter++
		mutex.Unlock()
	}
}

// --------------------- sync.WaitGroup approach ----------------------
func taskWaitGroup(wg *sync.WaitGroup) {
	defer wg.Done()
	// Simulate some work being done by the goroutine
	time.Sleep(time.Millisecond * 50)
}

// --------------------- Channel approach ----------------------
func taskChannel(wg *sync.WaitGroup, resultChannel chan int) {
	defer wg.Done()
	total := 0
	for i := 0; i < 10000; i++ {
		total++
	}
	resultChannel <- total
}

func TestSyncMutex(t *testing.T) {
	var wg sync.WaitGroup
	var mutex sync.Mutex
	numGoroutines := 1000
	// Reset the counter
	counter = 0

	// Start timer before running goroutines
	start := time.Now()

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go taskMutex(&wg, &mutex)
	}

	// Wait for all goroutines to finish
	wg.Wait()

	// Calculate execution time
	end := time.Now()
	executionTime := end.Sub(start)

	fmt.Printf("Execution time with %d goroutines (Mutex): %s\n", numGoroutines, executionTime)
	fmt.Printf("Shared Counter Value: %d\n", counter)
}

func TestSyncWaitGroup(t *testing.T) {
	var wg sync.WaitGroup
	numGoroutines := 1000

	// Start timer before running goroutines
	start := time.Now()

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go taskWaitGroup(&wg)
	}

	// Wait for all goroutines to finish
	wg.Wait()

	// Calculate execution time
	end := time.Now()
	executionTime := end.Sub(start)

	fmt.Printf("Execution time with %d goroutines (WaitGroup): %s\n", numGoroutines, executionTime)
}

func TestChannelApproach(t *testing.T) {
	var wg sync.WaitGroup
	numGoroutines := 1000
	resultChannel := make(chan int, numGoroutines)
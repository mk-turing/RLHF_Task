package main

import (
	"fmt"
	"log"
	"sync"
	"sync/atomic"
	"time"
)

// Callback type for our callbacks
type Callback func()

// CallbackManager manages a list of callbacks
type CallbackManager struct {
	callbacks []Callback
	mu        sync.Mutex
	wg        sync.WaitGroup
}

// Add adds a callback to the manager
func (cm *CallbackManager) Add(cb Callback) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	cm.callbacks = append(cm.callbacks, cb)
}

// Execute executes all callbacks concurrently
func (cm *CallbackManager) Execute(numWorkers int) {
	cm.wg.Add(numWorkers)
	for i := 0; i < numWorkers; i++ {
		go func() {
			for {
				cm.mu.Lock()
				if len(cm.callbacks) == 0 {
					cm.mu.Unlock()
					break
				}
				cb := cm.callbacks[0]
				cm.callbacks = cm.callbacks[1:]
				cm.mu.Unlock()
				cb()
			}
			cm.wg.Done()
		}()
	}
	cm.wg.Wait()
}

// TraceCallback wraps a callback and adds tracing information
func TraceCallback(cb Callback, name string) Callback {
	return func() {
		log.Printf("Starting callback: %s\n", name)
		cb()
		log.Printf("Finished callback: %s\n", name)
	}
}

var raceConditionCount int64

// Simulate a race condition by incrementing a shared variable
func simulateRaceCondition() {
	atomic.AddInt64(&raceConditionCount, 1)
}

// Main function
func main() {
	cm := &CallbackManager{}

	// Adding callbacks with tracing and race condition simulation
	cm.Add(TraceCallback(func() {
		log.Printf("Callback 1: Sleeping for 1 second...\n")
		time.Sleep(1 * time.Second)
		simulateRaceCondition()
	}, "Callback 1"))

	cm.Add(TraceCallback(func() {
		log.Printf("Callback 2: Sleeping for 2 seconds...\n")
		time.Sleep(2 * time.Second)
		simulateRaceCondition()
	}, "Callback 2"))

	cm.Add(TraceCallback(func() {
		log.Printf("Callback 3: Sleeping for 3 seconds...\n")
		time.Sleep(3 * time.Second)
		simulateRaceCondition()
	}, "Callback 3"))

	// Execute all callbacks concurrently with 5 workers
	cm.Execute(5)

	// Log the race condition count
	fmt.Printf("Race condition count: %d\n", atomic.LoadInt64(&raceConditionCount))
}

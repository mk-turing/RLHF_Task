package main

import (
	"fmt"
	"log"
	"runtime"
	"sync"
	"time"
)

const (
	numCallbacks = 10
	concurrency  = 3 // Number of goroutines to simulate concurrency
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

// Execute executes all callbacks concurrently using goroutines
func (cm *CallbackManager) Execute() {
	cm.wg.Add(len(cm.callbacks))
	for _, cb := range cm.callbacks {
		go func(cb Callback) {
			defer cm.wg.Done()
			cb()
		}(cb)
	}
	cm.wg.Wait()
}

// TraceCallback wraps a callback and adds tracing information
func TraceCallback(cb Callback, name string) Callback {
	return func() {
		log.Printf("Starting callback: %s [Goroutine: %d]\n", name, runtime.NumGoroutine())
		cb()
		log.Printf("Finished callback: %s [Goroutine: %d]\n", name, runtime.NumGoroutine())
	}
}

// Main function
func main() {
	cm := &CallbackManager{}

	// Adding callbacks with tracing
	for i := 1; i <= numCallbacks; i++ {
		callbackName := fmt.Sprintf("Callback %d", i)
		cm.Add(TraceCallback(func() {
			time.Sleep(time.Duration(i*100) * time.Millisecond) // Simulate different execution times
		}, callbackName))
	}

	// Execute callbacks concurrently using goroutines
	cm.Execute()

	// Simulate further callback invocations after the initial set has completed
	go func() {
		time.Sleep(2 * time.Second) // Wait for some time before adding more callbacks
		cm.Add(TraceCallback(func() {
			log.Println("Additional Callback 1 executed.")
		}, "Additional Callback 1"))
		cm.Add(TraceCallback(func() {
			log.Println("Additional Callback 2 executed.")
		}, "Additional Callback 2"))
		cm.Execute()
	}()

	// Wait for all callbacks to complete before exiting
	cm.wg.Wait()
	log.Println("All callbacks completed.")
}

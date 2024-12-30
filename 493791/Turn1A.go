package main

import (
	"log"
	"sync"
	"time"
)

// Callback type for our callbacks
type Callback func()

// CallbackManager manages a list of callbacks
type CallbackManager struct {
	callbacks []Callback
	mu        sync.Mutex
}

// Add adds a callback to the manager
func (cm *CallbackManager) Add(cb Callback) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	cm.callbacks = append(cm.callbacks, cb)
}

// Execute executes all callbacks in sequence
func (cm *CallbackManager) Execute() {
	for _, cb := range cm.callbacks {
		cb()
	}
}

// TraceCallback wraps a callback and adds tracing information
func TraceCallback(cb Callback, name string) Callback {
	return func() {
		log.Printf("Starting callback: %s\n", name)
		cb()
		log.Printf("Finished callback: %s\n", name)
	}
}

// Main function
func main() {
	cm := &CallbackManager{}

	// Adding callbacks with tracing
	cm.Add(TraceCallback(func() {
		log.Printf("Callback 1: Sleeping for 1 second...\n")
		time.Sleep(1 * time.Second)
	}, "Callback 1"))

	cm.Add(TraceCallback(func() {
		log.Printf("Callback 2: Sleeping for 2 seconds...\n")
		time.Sleep(2 * time.Second)
	}, "Callback 2"))

	cm.Add(TraceCallback(func() {
		log.Printf("Callback 3: Sleeping for 3 seconds...\n")
		time.Sleep(3 * time.Second)
	}, "Callback 3"))

	// Execute all callbacks
	cm.Execute()
}

package main

import (
	"fmt"
	"log"
	"sort"
	"sync"
	"sync/atomic"
	"time"
)

// Callback type for our callbacks
type Callback func()

// PriorityCallback wraps a callback with a priority
type PriorityCallback struct {
	Callback Callback
	Priority int
}

// CallbackManager manages a list of callbacks with priorities
type CallbackManager struct {
	callbacks []*PriorityCallback
	mu        sync.Mutex
	wg        sync.WaitGroup
}

// Add adds a callback to the manager with a specified priority
func (cm *CallbackManager) Add(cb Callback, priority int) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	cm.callbacks = append(cm.callbacks, &PriorityCallback{Callback: cb, Priority: priority})
	// Re-sort the callbacks by priority
	sort.Sort(PriorityCallbackSlice(cm.callbacks))
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
				cb.Callback()
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
	}, "Callback 1"), 1)

	cm.Add(TraceCallback(func() {
		log.Printf("Callback 2: Sleeping for 2 seconds...\n")
		time.Sleep(2 * time.Second)
		simulateRaceCondition()
	}, "Callback 2"), 3)

	cm.Add(TraceCallback(func() {
		log.Printf("Callback 3: Sleeping for 3 seconds...\n")
		time.Sleep(3 * time.Second)
		simulateRaceCondition()
	}, "Callback 3"), 2)

	// Execute all callbacks concurrently with 5 workers
	cm.Execute(5)

	// Log the race condition count
	fmt.Printf("Race condition count: %d\n", atomic.LoadInt64(&raceConditionCount))

	// Simulate a change in priority
	cm.Add(TraceCallback(func() {
		log.Printf("Callback 4: Sleeping for 1 second...\n")
		time.Sleep(1 * time.Second)
		simulateRaceCondition()
	}, "Callback 4"), 4)

	// Execute again with updated priorities
	cm.Execute(5)

	// Log the race condition count again
	fmt.Printf("Race condition count: %d\n", atomic.LoadInt64(&raceConditionCount))
}

// Implement a custom sort for PriorityCallback
type PriorityCallbackSlice []*PriorityCallback

func (p PriorityCallbackSlice) Len() int           { return len(p) }
func (p PriorityCallbackSlice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p PriorityCallbackSlice) Less(i, j int) bool { return p[i].Priority > p[j].Priority }

package main

import (
	"log"
	"net/http"
	"sort"
	"sync"
	"sync/atomic"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Callback type for our callbacks
type Callback func()

// PriorityCallback wraps a callback with a priority and timestamps
type PriorityCallback struct {
	Callback Callback
	Priority int
	Start    time.Time
	End      time.Time
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

	// Corrected assignment: Remove :=, use = since cb is already declared.
	newCallback := &PriorityCallback{Callback: cb, Priority: priority}
	cm.callbacks = append(cm.callbacks, newCallback)

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
				cb.End = time.Now()
			}
			cm.wg.Done()
		}()
	}
	cm.wg.Wait()
}

// TraceCallback wraps a callback and adds tracing information
func TraceCallback(cb Callback, name string) Callback {
	return func() {
		cb()
		log.Printf("Finished callback: %s\n", name)
	}
}

var raceConditionCount int64

// Simulate a race condition by incrementing a shared variable
func simulateRaceCondition() {
	atomic.AddInt64(&raceConditionCount, 1)
}

// Simulate failure by introducing fault injection
var failureRate int64

func simulateFailure() {
	if atomic.LoadInt64(&failureRate) > 0 {
		if atomic.AddInt64(&failureRate, 1)%100 == 0 {
			log.Printf("Simulating failure: introducing a fault.\n")
			panic("Failed callback")
		}
	}
}

// Main function
func main() {
	cm := &CallbackManager{}

	// Create Prometheus metrics
	latencyHistogram := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "callback_execution_latency_seconds",
			Help: "Histogram of callback execution latency in seconds",
		},
		[]string{"priority"},
	)

	prometheus.MustRegister(latencyHistogram)

	// Adding callbacks with tracing and race condition simulation
	cm.Add(TraceCallback(func() {
		log.Printf("Callback 1: Sleeping for 1 second...\n")
		time.Sleep(1 * time.Second)
		simulateRaceCondition()
		simulateFailure()
	}, "Callback 1"), 1)

	cm.Add(TraceCallback(func() {
		log.Printf("Callback 2: Sleeping for 2 seconds...\n")
		time.Sleep(2 * time.Second)
		simulateRaceCondition()
		simulateFailure()
	}, "Callback 2"), 3)

	cm.Add(TraceCallback(func() {
		log.Printf("Callback 3: Sleeping for 3 seconds...\n")
		time.Sleep(3 * time.Second)
		simulateRaceCondition()
		simulateFailure()
	}, "Callback 3"), 2)

	// Execute all callbacks concurrently with 5 workers
	cm.Execute(5)

	// Analyze results
	log.Printf("Race condition count: %d\n", atomic.LoadInt64(&raceConditionCount))

	// Simulate a change in priority
	cm.Add(TraceCallback(func() {
		log.Printf("Callback 4: Sleeping for 1 second...\n")
		time.Sleep(1 * time.Second)
		simulateRaceCondition()
		simulateFailure()
	}, "Callback 4"), 4)

	// Execute again with updated priorities
	cm.Execute(5)

	// Analyze results again
	log.Printf("Race condition count: %d\n", atomic.LoadInt64(&raceConditionCount))

	// Enable fault injection for the next run
	atomic.StoreInt64(&failureRate, 100)

	// Execute again with fault injection
	cm.Execute(5)

	// Analyze results with fault injection
	log.Printf("Race condition count: %d\n", atomic.LoadInt64(&raceConditionCount))

	// Start the Prometheus HTTP server to export metrics
	http.Handle("/metrics", promhttp.Handler()) // Fixed Prometheus handler usage
	log.Println("Prometheus metrics are available at http://localhost:8080/metrics")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// Implement a custom sort for PriorityCallback
type PriorityCallbackSlice []*PriorityCallback

func (p PriorityCallbackSlice) Len() int           { return len(p) }
func (p PriorityCallbackSlice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p PriorityCallbackSlice) Less(i, j int) bool { return p[i].Priority > p[j].Priority }

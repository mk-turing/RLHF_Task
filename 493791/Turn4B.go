
package main

import (
	"fmt"
	"log"
	"net/http"
	"sort"
	"sync"
	"sync/atomic"
	"time"

	"github.com/braintree/manners"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// New metrics variables
var (
	callbackLatency = promauto.NewHistogram(prometheus.HistogramOpts{
		Name: "callback_latency_seconds",
		Help: "The latency distribution of callback executions.",
	})

	callbackOrderAccuracy = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "callback_order_accuracy",
		Help: "Indicates the accuracy of callback execution order.",
	}, []string{"callback_name"})

	callbackExecutionTime = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name: "callback_execution_time_seconds",
		Help: "The execution time distribution of callbacks.",
	}, []string{"callback_name"})
)

// Callback type for our callbacks
type Callback func()

// PriorityCallback wraps a callback with a priority and additional metadata for metrics
type PriorityCallback struct {
	Callback      Callback
	Priority      int
	CallbackName string
}

// CallbackManager manages a list of callbacks with priorities
type CallbackManager struct {
	callbacks []*PriorityCallback
	mu        sync.Mutex
	wg        sync.WaitGroup
}

func (cm *CallbackManager) Add(cb Callback, priority int, callbackName string) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	cm.callbacks = append(cm.callbacks, &PriorityCallback{Callback: cb, Priority: priority, CallbackName: callbackName})
	// Re-sort the callbacks by priority
	sort.Sort(PriorityCallbackSlice(cm.callbacks))
}

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
				startTime := time.Now()
				cb.Callback()
				executionTime := time.Since(startTime)
				callbackExecutionTime.WithLabelValues(cb.CallbackName).Observe(executionTime.Seconds())
			}
			cm.wg.Done()
		}()
	}
	cm.wg.Wait()
}

// TraceCallback wraps a callback and adds tracing information
func TraceCallback(cb Callback, name string) Callback {
	return func() {
		startTime := time.Now()
		log.Printf("Starting callback: %s\n", name)
		cb()
		latency := time.Since(startTime)
		log.Printf("Finished callback: %s with latency: %s\n", name, latency)
		callbackLatency.Observe(latency.Seconds())
	}
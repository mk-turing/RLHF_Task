package _91164

import (
	"math/rand"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// Stateful Goroutine
type ConcurrencyTest struct {
	mu           sync.Mutex
	value        int
	errors       int
	operations   int64
	stopped      bool
	wg           *sync.WaitGroup
	stressConfig *StressConfig
}

type StressConfig struct {
	NumWorkers    int
	NumOperations int
	Duration      time.Duration
	ErrorRate     float64
}

func (ct *ConcurrencyTest) Start(ch chan string) {
	ct.wg = &sync.WaitGroup{}
	ct.wg.Add(1)

	go func() {
		defer ct.wg.Done()
		for {
			select {
			case msg := <-ch:
				if msg == "inc" {
					ct.increment()
				} else if msg == "dec" {
					ct.decrement()
				} else if msg == "stop" {
					ct.mu.Lock()
					ct.stopped = true
					ct.mu.Unlock()
					return
				}
			default:
				time.Sleep(1 * time.Millisecond) // Simulate work
			}
		}
	}()
}

func (ct *ConcurrencyTest) increment() {
	ct.mu.Lock()
	ct.value++
	atomic.AddInt64(&ct.operations, 1)
	ct.mu.Unlock()
}

func (ct *ConcurrencyTest) decrement() {
	ct.mu.Lock()
	if ct.value > 0 {
		ct.value--
	} else {
		ct.errors++
	}
	atomic.AddInt64(&ct.operations, 1)
	ct.mu.Unlock()
}

func (ct *ConcurrencyTest) GetResults() (value, errors int, operations int64) {
	ct.mu.Lock()
	defer ct.mu.Unlock()
	return ct.value, ct.errors, ct.operations
}

func (ct *ConcurrencyTest) Wait() {
	ct.wg.Wait()
}

// Stress test function
func StressTest(t *testing.T, config *StressConfig) {
	ct := &ConcurrencyTest{stressConfig: config}
	ch := make(chan string)

	ct.Start(ch)

	// Start workers
	for i := 0; i < config.NumWorkers; i++ {
		go func() {
			for i := 0; i < config.NumOperations; i++ {
				if rand.Float64() < config.ErrorRate {
					ch <- "dec"
				} else {
					ch <- "inc"
				}
				time.Sleep(100 * time.Microsecond) // Simulate operation time
			}
		}()
	}

	// Wait for a duration or until all operations are done
	select {
	case <-time.After(config.Duration):
		ch <- "stop"
	case <-ch:
	}

	ct.Wait()

	value, errors, operations := ct.GetResults()

	t.Logf("Final Value: %d, Errors: %d, Operations: %d", value, errors, operations)

	// Assertions
	if errors > 0 {
		t.Errorf("Expected no errors, got %d", errors)
	}
}

func TestStatefulGoroutines(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	config := &StressConfig{
		NumWorkers:    100,
		NumOperations: 10000,
		Duration:      10 * time.Second,
		ErrorRate:     0.1,
	}

	StressTest(t, config)
}

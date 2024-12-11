package _91164

import (
	"math/rand"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// Stateful Goroutine
type FaultInjectionTest struct {
	mu          sync.Mutex
	value       int
	errors      int
	operations  int64
	stopped     bool
	wg          *sync.WaitGroup
	faultConfig *FaultConfig
}

type FaultConfig struct {
	InjectPanic   bool
	PanicRate     float64
	InjectDelay   bool
	DelayDuration time.Duration
}

func (ct *FaultInjectionTest) Start(ch chan string) {
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

func (ct *FaultInjectionTest) increment() {
	ct.mu.Lock()
	ct.value++
	atomic.AddInt64(&ct.operations, 1)
	ct.mu.Unlock()

	if ct.faultConfig.InjectPanic && rand.Float64() < ct.faultConfig.PanicRate {
		panic("Simulated panic")
	}

	if ct.faultConfig.InjectDelay {
		time.Sleep(ct.faultConfig.DelayDuration)
	}
}

func (ct *FaultInjectionTest) decrement() {
	ct.mu.Lock()
	if ct.value > 0 {
		ct.value--
	} else {
		ct.errors++
	}
	atomic.AddInt64(&ct.operations, 1)
	ct.mu.Unlock()

	if ct.faultConfig.InjectPanic && rand.Float64() < ct.faultConfig.PanicRate {
		panic("Simulated panic")
	}

	if ct.faultConfig.InjectDelay {
		time.Sleep(ct.faultConfig.DelayDuration)
	}
}

func (ct *FaultInjectionTest) GetResults() (value, errors int, operations int64) {
	ct.mu.Lock()
	defer ct.mu.Unlock()
	return ct.value, ct.errors, ct.operations
}

func (ct *FaultInjectionTest) Wait() {
	ct.wg.Wait()
}

// Fault injection test function
func FaultInjectionTestFunc(t *testing.T, config *FaultConfig) {
	ct := &FaultInjectionTest{faultConfig: config}
	ch := make(chan string)

	ct.Start(ch)

	// Start workers
	for i := 0; i < 100; i++ {
		go func() {
			for i := 0; i < 10000; i++ {
				if rand.Intn(2) == 0 {
					ch <- "inc"
				} else {
					ch <- "dec"
				}
				time.Sleep(100 * time.Microsecond) // Simulate operation time
			}
		}()
	}

	ch <- "stop"
	ct.Wait()

	value, errors, operations := ct.GetResults()

	t.Logf("Final Value: %d, Errors: %d, Operations: %d", value, errors, operations)

	// Assertions: Adjust these based on expected behavior post-faults
	if errors > 0 {
		t.Errorf("Expected no errors, got %d", errors)
	}
}

func TestFaultInjection(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	config := &FaultConfig{
		InjectPanic:   true,
		PanicRate:     0.01, // 1% chance of panic
		InjectDelay:   true,
		DelayDuration: 10 * time.Millisecond,
	}

	FaultInjectionTestFunc(t, config)
}

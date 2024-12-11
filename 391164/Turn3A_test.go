package _91164

import (
	"fmt"
	"math/rand"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// Stateful Goroutine
type StatefulService struct {
	mu         sync.Mutex
	value      int
	errors     int32
	operations int64
}

func (ss *StatefulService) Increment(fault bool) {
	ss.mu.Lock()
	defer ss.mu.Unlock()

	if fault && rand.Float64() < 0.5 {
		// Simulate a panic
		panic("simulated panic")
	}

	ss.value++
	atomic.AddInt64(&ss.operations, 1)
}

func (ss *StatefulService) Decrement(fault bool) {
	ss.mu.Lock()
	defer ss.mu.Unlock()

	if fault && rand.Float64() < 0.5 {
		ss.errors++
		return // Simulate an error
	}

	if ss.value > 0 {
		ss.value--
	} else {
		ss.errors++
	}
	atomic.AddInt64(&ss.operations, 1)
}

func (ss *StatefulService) GetState() (int, int32, int64) {
	ss.mu.Lock()
	defer ss.mu.Unlock()
	return ss.value, ss.errors, ss.operations
}

// Fault Injection Test
func FaultInjectionTest(t *testing.T, iterations int, faultRate float64) {
	ss := &StatefulService{}
	var wg sync.WaitGroup

	for i := 0; i < iterations; i++ {
		wg.Add(2) // Two operations: increment and decrement

		go func() {
			defer wg.Done()
			fault := rand.Float64() < faultRate
			defer func() {
				if r := recover(); r != nil {
					// Handle panic
					t.Log("Recovered from panic:", r)
				}
			}()
			ss.Increment(fault)
		}()

		go func() {
			defer wg.Done()
			fault := rand.Float64() < faultRate
			ss.Decrement(fault)
		}()
	}

	wg.Wait()

	value, errors, operations := ss.GetState()
	t.Logf("Final Value: %d, Errors: %d, Operations: %d", value, errors, operations)

	// Assertions
	if errors > 0 {
		t.Errorf("Expected no errors, but got %d", errors)
	}
}

func TestFaultInjection(t *testing.T) {
	rand.Seed(time.Now().UnixNano())

	// Parameters for the fault injection test
	iterations := 1000
	faultRate := 0.1 // 10% chance to inject a fault

	FaultInjectionTest(t, iterations, faultRate)
}

func main() {
	fmt.Println("Running fault injection tests on stateful Goroutines.")
}

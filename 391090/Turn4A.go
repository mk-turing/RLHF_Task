package main

import (
	"fmt"
	"log"
	"math/rand"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

var (
	operationCount uint64
	errorCount     uint64
	goroutineCount = 0
	maxGoroutines  = 100
)

type Workload struct {
	Type        string
	Generator   func()
	Count       int
	Concurrency int
}

type ResultCollector struct {
	responseTimes []time.Duration
	sync.Mutex
}

func (rc *ResultCollector) Add(d time.Duration) {
	rc.Lock()
	rc.responseTimes = append(rc.responseTimes, d)
	rc.Unlock()
}

func main() {
	rand.Seed(time.Now().UnixNano())
	runtime.GOMAXPROCS(8)

	workload := Workload{
		Type:        "mixed",
		Generator:   GenerateWorkload,
		Count:       1000,
		Concurrency: 100,
	}

	var wg sync.WaitGroup

	rc := &ResultCollector{}

	// Start workload generator
	for i := 0; i < workload.Concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < workload.Count; j++ {
				workload.Generator()
			}
		}()
	}

	go MonitorSystem(rc)

	wg.Wait()

	fmt.Printf("Total operations: %d, Total errors: %d\n", atomic.LoadUint64(&operationCount), atomic.LoadUint64(&errorCount))
	fmt.Printf("Average response time: %.2fms\n", float64(rc.CalculateAverageResponseTime())/1000)
}

func GenerateWorkload() {
	operation := rand.Intn(2)

	start := time.Now()

	if operation == 0 { // Simulate read operation
		time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
	} else { // Simulate write operation
		time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
		atomic.AddUint64(&errorCount, 1) // Simulate a potential error
	}

	end := time.Now()
	rc := &ResultCollector{}
	rc.Add(end.Sub(start))
	atomic.AddUint64(&operationCount, 1)
}

func MonitorSystem(rc *ResultCollector) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		currentGoroutines := runtime.NumGoroutine()
		if currentGoroutines > maxGoroutines {
			log.Printf("Detected Goroutine leak: %d Goroutines. Adjusting workload...\n", currentGoroutines)
			// Adjust workload to reduce Goroutine count
		}

		// Monitor other metrics like CPU, memory, etc.
		// Implement logic to respond to issues like resource contention
	}
}

func (rc *ResultCollector) CalculateAverageResponseTime() time.Duration {
	if len(rc.responseTimes) == 0 {
		return 0
	}
	total := time.Duration(0)
	for _, t := range rc.responseTimes {
		total += t
	}
	return total / time.Duration(len(rc.responseTimes))
}

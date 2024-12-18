package main

import (
	"fmt"
	"log"
	"math/rand"
	"runtime"
	"sync"
	"sync/atomic"
	"time"

	"golang.org/x/time/rate"
)

var (
	operationCount uint64
	errorCount     uint64
)

type Workload struct {
	Type       string
	Generator  func()
	Count      int
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
		Type:       "mixed",
		Generator:  GenerateWorkload,
		Count:      1000,
		Concurrency: 100,
	}

	var wg sync.WaitGroup

	rc := &ResultCollector{}

	// Start workload generator
	go func() {
		simulateWorkload(workload, rc)
	}()

	// Start real-time monitoring
	go func() {
		monitorGoroutines()
	}()

	wg.Wait()

	fmt.Printf("Total operations: %d, Total errors: %d\n", atomic.LoadUint64(&operationCount), atomic.LoadUint64(&errorCount))
	fmt.Printf("Average response time: %.2fms\n", float64(rc.CalculateAverageResponseTime())/1000)
}

func simulateWorkload(workload Workload, rc *ResultCollector) {
	// Create a rate limiter to control workload intensity
	limiter := rate.NewLimiter(rate.Every(time.Second/100), workload.Concurrency)

	for i := 0; i < workload.Count; i++ {
		// Wait for the next token from the rate limiter
		limiter.Wait(nil)

		workload.Generator()
	}
}

func GenerateWorkload() {
	operation := rand.Intn(2)

	start := time.Now()

	if operation == 0 { // Simulate read operation
		time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
	} else
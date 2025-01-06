package main

import (
	"fmt"
	"net/rpc"
	"sync"
	"time"
)

// Benchmark interface definition
type Benchmark interface {
	Run() (time.Duration, error)
}

// Test struct example
type Test1 struct{}

func (t *Test1) Run() (time.Duration, error) {
	// Benchmark logic here
	return 1 * time.Second, nil
}

// Master node RPC server
type MasterRPC struct {
	workers  []*rpc.Client
	results  chan Result
	sync.WaitGroup
}

func (m *MasterRPC) StartBenchmark(test Benchmark, reply *string) error {
	// Distribute test to all workers
	m.Add(len(m.workers))
	for _, worker := range m.workers {
		go func(w *rpc.Client) {
			var workerResult Result
			w.Call("WorkerRPC.RunBenchmark", test, &workerResult)
			m.results <- workerResult
			m.Done()
		}(worker)
	}
	m.Wait()
	return nil
}

// Worker node RPC server
type WorkerRPC struct {
	Results chan Result
}

func (w *WorkerRPC) RunBenchmark(test Benchmark, reply *Result) error {
	// Run benchmark in parallel on this worker
	var wg sync.WaitGroup
	numIterations := 10 // Assuming you want to run each test multiple times to get an average
	wg.Add(numIterations)
	durations := make([]time.Duration, numIterations)

	for i := 0; i < numIterations; i++ {
		go func() {
			duration, _ := test.Run()
			durations[i] = duration
			wg.Done()
		}()
	}
	wg.Wait()

	// Calculate mean of durations
	var sum time.Duration
	for _, duration := range durations {
		sum += duration
	}
	mean := sum / time.Duration(numIterations)

	// Send result back to master
	w.Results <- Result{Mean: mean}
	return nil
}

// Result struct to hold results from workers
type Result struct {
	Mean time.Duration
}

func main() {
	// Master node setup
	masterRPC := new(MasterRPC)
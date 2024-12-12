package main

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// Resource represents a monitored system resource.
type Resource struct {
	name       string
	current    float64 // Current usage (e.g., CPU percentage, memory used)
	limit      float64 // Maximum allowed usage
	limiter    chan struct{}
	adjustment time.Duration
}

// Task function simulates a workload with different resource demands.
func task(ctx context.Context, id int, cpu float64, mem float64, net float64, resources map[string]*Resource) {
	select {
	case <-time.After(time.Duration(rand.Intn(500)) * time.Millisecond):
		fmt.Printf("Task %d completed (CPU: %.2f, Mem: %.2f, Net: %.2f)\n", id, cpu, mem, net)
		// Simulate resource usage
		resources["cpu"].current += cpu
		resources["mem"].current += mem
		resources["net"].current += net
	case <-ctx.Done():
		fmt.Printf("Task %d cancelled\n", id)
	}
}

// Worker function processes tasks, using semaphores for resource throttling.
func worker(ctx context.Context, wg *sync.WaitGroup, id int, resources map[string]*Resource) {
	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			return
		default:
			// Acquire resources
			if acquireResources(resources) {
				// Simulate task workload
				task(ctx, id, rand.Float64()*10, rand.Float64()*10, rand.Float64()*10, resources)
				releaseResources(resources)
			}
		}
	}
}

// AcquireResources attempts to acquire resources based on their limits.
func acquireResources(resources map[string]*Resource) bool {
	for _, resource := range resources {
		select {
		case resource.limiter <- struct{}{}:
		default:
			fmt.Printf("Task dropped due to resource limit on %s: %.2f / %.2f\n", resource.name, resource.current, resource.limit)
			return false
		}
	}
	return true
}

// ReleaseResources releases the acquired resources.
func releaseResources(resources map[string]*Resource) {
	for _, resource := range resources {
		<-resource.limiter
	}
}

// MonitorResources updates resource usage dynamically.
func monitorResources(ctx context.Context, resources map[string]*Resource) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			// Simulate load fluctuation
			for _, resource := range resources {
				resource.current = rand.Float64()*resource.limit + 1 // Ensure some usage
				fmt.Printf("%s usage: %.2f / %.2f\n", resource.name, resource.current, resource.limit)
			}
			time.Sleep(2 * time.Second)
		}
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())

	const maxWorkers = 5
	const totalTasks = 20
	resources := make(map[string]*Resource)

	resources["cpu"] = &Resource{
		name:       "CPU",
		current:    0.0,
		limit:      100.0,
		limiter:    make(chan struct{}, 10), // Max concurrent CPU tasks
		adjustment: 1 * time.Second,
	}
	resources["mem"] = &Resource{
		name:       "Memory",
		current:    0.0,
		limit:      100.0,
		limiter:    make(chan struct{}, 8), // Max concurrent Memory tasks
		adjustment: 1 * time.Second,
	}
	resources["net"] = &Resource{
		name:       "Network",
		current:    0.0,
		limit:      100.0,
		limiter:    make(chan struct{}, 5), // Max concurrent Network tasks
		adjustment: 1 * time.Second,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var wg sync.WaitGroup

	// Start resource monitoring
	go monitorResources(ctx, resources)

	// Starting worker goroutines
	for i := 0; i < maxWorkers; i++ {
		wg.Add(1)
		go worker(ctx, &wg, i, resources)
	}

	// Sending tasks
	for i := 0; i < totalTasks; i++ {
		select {
		case <-ctx.Done():
			break
		default:
			// Tasks could be prioritized differently based on strategy here
			time.Sleep(100 * time.Millisecond) // Simulate incoming tasks
		}
	}

	// Wait for all workers to finish
	wg.Wait()
}

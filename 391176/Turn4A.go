package main

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

// Resource structure for managing resource limits in Redis.
type Resource struct {
	name   string
	limit  int
	client *redis.Client
}

// NewResource creates a new Resource with Redis client.
func NewResource(name string, limit int) *Resource {
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	return &Resource{name: name, limit: limit, client: rdb}
}

// Task function simulates workload.
func task(nodeID, id int, resourceName string) {
	select {
	case <-time.After(time.Duration(rand.Intn(500)) * time.Millisecond):
		fmt.Printf("Node %d: Task %d completed (Resource: %s)\n", nodeID, id, resourceName)
	}
}

// RequestResource attempts to acquire a resource from Redis.
func (r *Resource) RequestResource(nodeID int) bool {
	acquired, err := r.client.SetNX(ctx, r.name+fmt.Sprintf(":%d", nodeID), 1, 10*time.Second).Result()
	if err != nil {
		fmt.Println("Error requesting resource:", err)
		return false
	}
	return acquired
}

// ReleaseResource releases the requested resource.
func (r *Resource) ReleaseResource(nodeID int) {
	r.client.Del(ctx, r.name+fmt.Sprintf(":%d", nodeID))
}

// SimulateLoad adjusts resource limits based on the load.
func (r *Resource) SimulateLoad(adjustmentInterval time.Duration) {
	for {
		loadValue := rand.Intn(10) + 1 // Simulated load between 1 and 10
		newLimit := 5 - loadValue
		if newLimit < 1 {
			newLimit = 1
		}
		r.limit = newLimit
		fmt.Printf("Adjusting %s limit to %d based on load %d\n", r.name, r.limit, loadValue)
		time.Sleep(adjustmentInterval)
	}
}

// Node representation for each worker node.
type Node struct {
	id       int
	resource *Resource
}

func (n *Node) Start(wg *sync.WaitGroup, totalTasks int) {
	defer wg.Done()
	for i := 0; i < totalTasks; i++ {
		if n.resource.RequestResource(n.id) {
			task(n.id, i, n.resource.name)
			n.resource.ReleaseResource(n.id)
		} else {
			fmt.Printf("Node %d: Task %d dropped due to %s load\n", n.id, i, n.resource.name)
		}
		time.Sleep(100 * time.Millisecond) // Simulate incoming tasks
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())

	const totalTasks = 20
	cpuResource := NewResource("cpu", 5)
	memoryResource := NewResource("memory", 5)
	networkResource := NewResource("network", 5)

	var wg sync.WaitGroup

	// Simulate resource load adjustments
	go cpuResource.SimulateLoad(5 * time.Second)
	go memoryResource.SimulateLoad(5 * time.Second)
	go networkResource.SimulateLoad(5 * time.Second)

	// Create worker nodes and start them
	nodes := []Node{
		{id: 1, resource: cpuResource},
		{id: 2, resource: memoryResource},
		{id: 3, resource: networkResource},
	}

	for _, node := range nodes {
		wg.Add(1)
		go node.Start(&wg, totalTasks)
	}

	// Wait for all workers to finish
	wg.Wait()
}

package main

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// Resource represents resource state.
type Resource struct {
	limit   int
	current int
	lock    sync.RWMutex
}

// NodeStatus reports a node's resource status.
type NodeStatus struct {
	ID        string
	Resources map[string]*Resource
}

// CentralStore simulates a centralized key-value store.
type CentralStore struct {
	mu sync.Mutex

	globalLimits map[string]*Resource
	nodeStatuses map[string]*NodeStatus
}

func (c *CentralStore) GetGlobalLimit(resource string) *Resource {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.globalLimits[resource]
}

func (c *CentralStore) SetGlobalLimit(resource string, limit *Resource) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.globalLimits[resource] = limit
}

func (c *CentralStore) RegisterNode(nodeStatus *NodeStatus) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.nodeStatuses[nodeStatus.ID] = nodeStatus
}

func (c *CentralStore) GetAllNodeStatuses() []*NodeStatus {
	c.mu.Lock()
	defer c.mu.Unlock()
	var result []*NodeStatus
	for _, status := range c.nodeStatuses {
		result = append(result, status)
	}
	return result
}

// NewCentralStore creates a new CentralStore instance.
func NewCentralStore() *CentralStore {
	return &CentralStore{
		globalLimits: make(map[string]*Resource),
		nodeStatuses: make(map[string]*NodeStatus),
	}
}

// Task simulates work with a resource.
func task(ctx context.Context, id int, resource string, nodeID string, store *CentralStore) {
	select {
	case <-ctx.Done():
		return
	default:
		globalLimit := store.GetGlobalLimit(resource)
		if globalLimit == nil {
			fmt.Printf("Node %s: Task %d dropped due to no global limit for %s\n", nodeID, id, resource)
			return
		}

		nodeStatus := store.nodeStatuses[nodeID]
		if nodeStatus == nil {
			fmt.Printf("Node %s: Task %d dropped due to node not registered\n", nodeID, id)
			return
		}

		localResource := nodeStatus.Resources[resource]
		if localResource == nil {
			fmt.Printf("Node %s: Task %d dropped due to no local resource for %s\n", nodeID, id, resource)
			return
		}

		localResource.lock.Lock()
		defer localResource.lock.Unlock()
		//globalResource.lock.Lock()
		//defer globalResource.lock.Unlock()
		//
		//if localResource.current < localResource.limit && globalResource.current < globalResource.limit {
		//	localResource.current++
		//	globalResource.current++
		//	fmt.Printf("Node %s: Task %d started (Resource: %s, Local: %d/%d, Global: %d/%d)\n",
		//		nodeID, id, resource, localResource.current, localResource.limit, globalResource.current, globalResource.limit)
		//
		//	select {
		//	case <-time.After(time.Duration(rand.Intn(500)) * time.Millisecond):
		//		localResource.current--
		//		globalResource.current--
		//	case <-ctx.Done():
		//		fmt.Printf("Node %s: Task %d cancelled (Resource: %s)\n", nodeID, id, resource)
		//	}
		//} else {
		//	fmt.Printf("Node %s: Task %d dropped due to resource load (Local: %d/%d, Global: %d/%d)\n",
		//		nodeID, id, localResource.current, localResource.limit, globalResource.current, globalResource.limit)
		//}
	}
}

// SimulateNode simulates a node with resource handling and reporting.
func SimulateNode(ctx context.Context, nodeID string, totalTasks int, store *CentralStore) {
	resources := map[string]*Resource{
		"CPU":     {limit: 3, current: 0},
		"Memory":  {limit: 3, current: 0},
		"Network": {limit: 3, current: 0},
	}

	store.RegisterNode(&NodeStatus{ID: nodeID, Resources: resources})
	time.Sleep(time.Second) // Allow registration to propagate

	var wg sync.WaitGroup

	for _, resource := range resources {
		go func(resource *Resource) {
			ticker := time.NewTicker(2 * time.Second)
			defer ticker.Stop()
			for range ticker.C {
				resource.lock.Lock()
				if resource.current > 0 {
					resource.current--
				}
				resource.lock.Unlock()
			}
		}(resource)
	}

	for i := 0; i < totalTasks; i++ {
		select {
		case <-ctx.Done():
			break
		default:
			resourceType := []string{"CPU", "Memory", "Network"}[rand.Intn(3)]
			wg.Add(1)
			go task(ctx, i, resourceType, nodeID, store)
		}
		time.Sleep(100 * time.Millisecond) // Simulate incoming tasks
	}

	wg.Wait()
}

func main() {
	rand.Seed(time.Now().UnixNano())

	store := NewCentralStore()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Start simulation on multiple nodes
	go SimulateNode(ctx, "Node1", 10, store)
	go SimulateNode(ctx, "Node2", 10, store)

	time.Sleep(10 * time.Second)

	// Distributed global limit coordination (mocked here)
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()
	for range ticker.C {
		fmt.Println("Coordinating global limits:")
		nodeStatuses := store.GetAllNodeStatuses()
		var totalLimits map[string]int
		if len(nodeStatuses) == 0 {
			continue
		}

		totalLimits = make(map[string]int, len(nodeStatuses[0].Resources))
		for _, nodeStatus := range nodeStatuses {
			for resource, res := range nodeStatus.Resources {
				totalLimits[resource] += res.limit
			}
		}

		for resource, totalLimit := range totalLimits {
			store.SetGlobalLimit(resource, &Resource{limit: totalLimit, current: 0})
		}
	}
}

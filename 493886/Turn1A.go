package _93886

import (
	"fmt"
	"sync"
)

// DI struct containing the dependencies map and synchronization primitives
type DIContainer struct {
	dependencies map[string]interface{}
	mu           sync.RWMutex
}

// NewDIContainer creates a new DI container
func NewDIContainer() *DIContainer {
	return &DIContainer{dependencies: make(map[string]interface{})}
}

// Register registers a new dependency with the given key
func (container *DIContainer) Register(key string, dependency interface{}) {
	container.mu.Lock()
	defer container.mu.Unlock()
	container.dependencies[key] = dependency
	fmt.Println("Registered:", key)
}

// Resolve retrieves a dependency by its key
func (container *DIContainer) Resolve(key string) interface{} {
	container.mu.RLock()
	defer container.mu.RUnlock()
	return container.dependencies[key]
}

// Cleanup removes all dependencies from the container
func (container *DIContainer) Cleanup() {
	container.mu.Lock()
	defer container.mu.Unlock()
	delete(container.dependencies)
	container.dependencies = make(map[string]interface{})
}

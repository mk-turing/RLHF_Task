package main

import (
	"fmt"
	"sync"
)

// Dynamic map resizing with memory pooling
type DynamicMap struct {
	mu     sync.RWMutex
	data   map[string]DependencyInterface
	pool   *sync.Pool
	cap    int
	growth int
}

func NewDynamicMap(initialCap int, growth int) *DynamicMap {
	if initialCap <= 0 {
		initialCap = 4
	}
	if growth <= 0 {
		growth = initialCap
	}
	return &DynamicMap{
		data:   make(map[string]DependencyInterface, initialCap),
		pool:   &sync.Pool{New: func() interface{} { return make(map[string]DependencyInterface) }},
		cap:    initialCap,
		growth: growth,
	}
}

func (m *DynamicMap) Set(key string, value DependencyInterface) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Check if map needs resizing
	if len(m.data) >= m.cap {
		m.resize()
	}

	m.data[key] = value
}

func (m *DynamicMap) Get(key string) DependencyInterface {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.data[key]
}

func (m *DynamicMap) Delete(key string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.data, key)
}

func (m *DynamicMap) resize() {
	oldMap := m.data
	m.data = m.pool.Get().(map[string]DependencyInterface)
	m.cap += m.growth

	// Copy old entries to the new map
	for key, value := range oldMap {
		m.data[key] = value
	}

	// Return the old map to the pool
	m.pool.Put(oldMap)
}

// Dependency injection framework using dynamic map
type DependencyInjector struct {
	deps *DynamicMap
}

func NewDependencyInjector() *DependencyInjector {
	return &DependencyInjector{
		deps: NewDynamicMap(0, 16), // Initialize with reasonable growth factor
	}
}

func (d *DependencyInjector) Register(name string, dep DependencyInterface) {
	d.deps.Set(name, dep)
}

func (d *DependencyInjector) Get(name string) DependencyInterface {
	return d.deps.Get(name)
}

func main() {
	injector := NewDependencyInjector()
	dep1 := &MyDependency{name: "dep1"}
	dep2 := &MyDependency{name: "dep2"}

	// Register dependencies
	injector.Register("dep1", dep1)
	injector.Register("dep2", dep2)

	// Use dependencies
	fmt.Println(injector.Get("dep1").DoSomething()) // Output: dep1
	fmt.Println(injector.Get("dep2").DoSomething()) // Output: dep2

	// Delete a dependency
	injector.deps.Delete("dep2")
}

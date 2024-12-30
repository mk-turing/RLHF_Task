package main

import (
	"fmt"
	"sync"
)

type DependencyInterface interface {
	DoSomething() string
}

type MyDependency struct {
	name string
}

func (m *MyDependency) DoSomething() string {
	return m.name
}

type DependencyInjector struct {
	dependencies    sync.RWMutex
	dependenciesMap map[string]DependencyInterface
	pool            *Pool
}

func NewDependencyInjector(initialCapacity int) *DependencyInjector {
	return &DependencyInjector{
		dependenciesMap: make(map[string]DependencyInterface, initialCapacity),
		pool:            NewPool(),
	}
}

func (d *DependencyInjector) Register(name string, dep DependencyInterface) {
	d.dependencies.Lock()
	defer d.dependencies.Unlock()

	d.dependenciesMap[name] = dep
}

func (d *DependencyInjector) Get(name string) DependencyInterface {
	d.dependencies.RLock()
	defer d.dependencies.RUnlock()

	return d.dependenciesMap[name]
}

func (d *DependencyInjector) Release(dep DependencyInterface) {
	if dep, ok := dep.(*MyDependency); ok {
		d.pool.Put(dep)
	}
}

type Pool struct {
	mu      sync.Mutex
	cache   []*MyDependency
	maxSize int
}

func NewPool(maxSize ...int) *Pool {
	if len(maxSize) == 0 {
		maxSize = append(maxSize, 100) // Default max size
	}
	return &Pool{
		maxSize: maxSize[0],
	}
}

func (p *Pool) Get() *MyDependency {
	p.mu.Lock()
	defer p.mu.Unlock()

	if len(p.cache) == 0 {
		return new(MyDependency)
	}
	return p.cache[len(p.cache)-1]
}

func (p *Pool) Put(dep *MyDependency) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if len(p.cache) >= p.maxSize {
		return
	}

	dep.name = "" // Reset the name for reuse
	p.cache = append(p.cache, dep)
}

func main() {
	injector := NewDependencyInjector(10) // Initial capacity of 10

	// Register some dependencies concurrently
	go func() {
		dep1 := &MyDependency{name: "dep1"}
		injector.Register("dep1", dep1)
	}()

	go func() {
		dep2 := injector.pool.Get()
		dep2.name = "dep2"
		injector.Register("dep2", dep2)
	}()

	// Access the dependencies concurrently
	go func() {
		fmt.Println(injector.Get("dep1").DoSomething()) // Output: dep1
	}()

	go func() {
		fmt.Println(injector.Get("dep2").DoSomething()) // Output: dep2
	}()

	// Sleep for a moment to ensure all goroutines have finished
	select {}
}

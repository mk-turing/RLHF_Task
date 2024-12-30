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
	mu           sync.RWMutex
	dependencies map[string]func() DependencyInterface
	pool         *Pool
}

func NewDependencyInjector() *DependencyInjector {
	return &DependencyInjector{
		dependencies: make(map[string]func() DependencyInterface),
		pool:         NewPool(),
	}
}

func (d *DependencyInjector) Register(name string, depFunc func() DependencyInterface) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.dependencies[name] = depFunc
}

func (d *DependencyInjector) Get(name string) DependencyInterface {
	d.mu.RLock()
	depFunc, ok := d.dependencies[name]
	d.mu.RUnlock()

	if !ok {
		return nil // Handle non-existent dependency
	}

	// Lazy initialization of the dependency
	d.mu.Lock()
	defer d.mu.Unlock()
	if _, ok := d.dependencies[name]; ok { // Check again in case of race condition
		dep := depFunc()
		d.dependencies[name] = func() DependencyInterface { return dep }
		return dep
	}
	return nil // Another goroutine already initialized the dependency
}

func (d *DependencyInjector) Release(dep DependencyInterface) {
	d.mu.Lock()
	defer d.mu.Unlock()
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
	injector := NewDependencyInjector()

	// Register lazy dependency initialization functions
	injector.Register("dep1", func() DependencyInterface {
		return &MyDependency{name: "dep1"}
	})

	injector.Register("dep2", func() DependencyInterface {
		dep2 := injector.pool.Get()
		dep2.name = "dep2"
		return dep2
	})

	// Application startup time remains low as no dependencies are created yet.

	// Accessing dependencies will trigger their lazy initialization.
	fmt.Println(injector.Get("dep1").DoSomething()) // Output: dep1
	fmt.Println(injector.Get("dep2").DoSomething()) // Output: dep2

	// Release dependencies
	injector.Release(injector.Get("dep2"))
	injector.Release(injector.Get("dep1"))
}

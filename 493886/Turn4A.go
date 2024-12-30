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

type LazyDependency struct {
	initializer func() DependencyInterface
	once        sync.Once
	dep         DependencyInterface
}

func newLazyDependency(initializer func() DependencyInterface) *LazyDependency {
	return &LazyDependency{
		initializer: initializer,
	}
}

func (l *LazyDependency) Get() DependencyInterface {
	l.once.Do(func() {
		l.dep = l.initializer()
	})
	return l.dep
}

type DependencyInjector struct {
	mu           sync.RWMutex
	dependencies map[string]*LazyDependency
	pool         *Pool
}

func NewDependencyInjector(initialCapacity int) *DependencyInjector {
	return &DependencyInjector{
		dependencies: make(map[string]*LazyDependency, initialCapacity),
		pool:         NewPool(),
	}
}

func (d *DependencyInjector) Register(name string, initializer func() DependencyInterface) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.dependencies[name] = newLazyDependency(initializer)
}

func (d *DependencyInjector) Get(name string) DependencyInterface {
	d.mu.RLock()
	defer d.mu.RUnlock()
	if lazy, ok := d.dependencies[name]; ok {
		return lazy.Get()
	}
	return nil
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
	injector := NewDependencyInjector(10)

	// Register lazy dependencies
	injector.Register("dep1", func() DependencyInterface {
		return &MyDependency{name: "dep1"}
	})

	injector.Register("dep2", func() DependencyInterface {
		dep2 := injector.pool.Get()
		dep2.name = "dep2"
		return dep2
	})

	// Access the dependencies
	fmt.Println(injector.Get("dep1").DoSomething()) // Output: dep1

	// Wait for the goroutine to finish registering dep2
	// ...

	fmt.Println(injector.Get("dep2").DoSomething()) // Output: dep2

	// Release dependencies
	injector.Release(injector.Get("dep2"))
	injector.Release(injector.Get("dep1"))
}

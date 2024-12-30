package main

import "fmt"

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
	dependencies map[string]DependencyInterface
}

func NewDependencyInjector() *DependencyInjector {
	return &DependencyInjector{
		dependencies: make(map[string]DependencyInterface),
	}
}

func (d *DependencyInjector) Register(name string, dep DependencyInterface) {
	d.dependencies[name] = dep
}

func (d *DependencyInjector) Get(name string) DependencyInterface {
	return d.dependencies[name]
}

func main() {
	injector := NewDependencyInjector()
	dep1 := &MyDependency{name: "dep1"}
	dep2 := &MyDependency{name: "dep2"}
	injector.Register("dep1", dep1)
	injector.Register("dep2", dep2)
	fmt.Println(injector.Get("dep1").DoSomething()) // Output: dep1
	fmt.Println(injector.Get("dep2").DoSomething()) // Output: dep2
}

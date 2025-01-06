package main

import (
	"fmt"
	"sync"
	"time"
)

type Benchmark interface {
	Run() (time.Duration, error)
}

type Test1 struct{}

func (t *Test1) Run() (time.Duration, error) {
	start := time.Now()
	// Benchmark logic here
	end := time.Now()
	return end.Sub(start), nil
}

func main() {
	var wg sync.WaitGroup
	tests := []Benchmark{&Test1{}}

	for _, test := range tests {
		wg.Add(1)
		go func(test Benchmark) {
			defer wg.Done()
			for i := 0; i < 10; i++ { // Run test 10 times
				duration, err := test.Run()
				if err != nil {
					fmt.Println("Error:", err)
					continue
				}
				fmt.Println("Test duration:", duration)
			}
		}(test)
	}

	wg.Wait()
}

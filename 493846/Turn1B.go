package main

import (
	"fmt"
	"sync"
	"time"
)

type MyState struct {
	Count int
	sync.Mutex
}

func main() {
	msgChan := make(chan int)
	var s MyState
	go func() {
		for i := range msgChan {
			s.Lock()
			s.Count += i
			s.Unlock()
		}
	}()
	// Simulate callback functions sending updates
	msgChan <- 1
	msgChan <- 2
	msgChan <- 3
	close(msgChan)
	time.Sleep(time.Millisecond * 100)
	fmt.Println("Final count:", s.Count) // Output should be Final count: 6
}

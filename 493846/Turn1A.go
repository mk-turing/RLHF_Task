package main

import (
	"fmt"
	"sync"
)

type sharedState struct {
	Value int
	mu    sync.RWMutex
}

func callback1(state *sharedState, ch chan<- bool) {
	state.mu.Lock()
	defer state.mu.Unlock()
	state.Value++
	ch <- true
}

func callback2(state *sharedState, ch chan<- bool) {
	state.mu.Lock()
	defer state.mu.Unlock()
	state.Value--
	ch <- true
}

func main() {
	state := &sharedState{Value: 0}
	ch := make(chan bool)

	go callback1(state, ch)
	go callback2(state, ch)

	<-ch // Wait for callback1 to finish
	<-ch // Wait for callback2 to finish

	fmt.Printf("Final state value: %d\n", state.Value)
}

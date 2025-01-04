package main

import (
	"fmt"
	"log"
	"os"
	"sync"
)

// Callback function types
type Callback func(string)
type HarmfulCallback func(string)

// A function that accepts multiple callbacks and uses them
func processEvent(message string, callbacks []Callback, harmfulCallback HarmfulCallback) {
	var wg sync.WaitGroup
	wg.Add(len(callbacks) + 1)

	for _, callback := range callbacks {
		go func(cb Callback) {
			defer wg.Done()
			cb(message)
		}(callback)
	}

	go func() {
		defer wg.Done()
		harmfulCallback(message)
	}()

	wg.Wait()
}

func benignCallback1(msg string) {
	fmt.Println("Benign callback 1:", msg)
}

func benignCallback2(msg string) {
	fmt.Println("Benign callback 2:", msg)
}

func maliciousCallback(msg string) {
	// Simulate a harmful action, e.g., execute system commands or access sensitive data
	fmt.Println("Malicious callback executed!")
	// For demonstration purposes, let's just sleep for a few seconds
	// In a real-world scenario, this could be a call to a function that deletes files or exposes data
	// time.Sleep(5 * time.Second)
}

func main() {
	// Registering benign callbacks
	benignCallbacks := []Callback{benignCallback1, benignCallback2}

	// Registering the malicious callback
	var harmfulCallback HarmfulCallback = maliciousCallback

	// Initialize a logger
	logger := log.New(os.Stdout, "Callback Execution: ", log.LstdFlags)

	// Process events with logging
	processEventWithLogging("Hello, World!", benignCallbacks, harmfulCallback, logger)

	// Process another event without the malicious callback
	benignCallbacks = []Callback{benignCallback1, benignCallback2}
	harmfulCallback = nil
	processEventWithLogging("Hello, Again!", benignCallbacks, harmfulCallback, logger)
}

func processEventWithLogging(message string, callbacks []Callback, harmfulCallback HarmfulCallback, logger *log.Logger) {
	// Use a context to track callback execution
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Wrap the callback execution with logging
	loggedCallbacks := make([]Callback, len(callbacks))
	for i, cb := range callbacks {
		loggedCallbacks[i] = func(msg string) {
			logger.Printf("Executing benign callback: %s", msg)
			cb(msg)
		}
	}

	if harmfulCallback != nil {
		loggedHarmfulCallback := func(msg string) {
			logger.Printf("Executing malicious callback: %s", msg)
			harmfulCallback(msg)
		}
		processEvent(message, loggedCallbacks, loggedHarmfulCallback)
	} else {
		processEvent(message, loggedCallbacks, nil)
	}
}

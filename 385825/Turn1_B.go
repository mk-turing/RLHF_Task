package main

import (
	"fmt"
	"time"
)

// Shared channel between service1 and service2
var sharedChannel = make(chan string)

func main() {
	// Start service2 (consumer) in a separate goroutine
	go service2()

	// Start service1 (producer)
	service1()
}

// Service 1: The producer
func service1() {
	for i := 0; i < 5; i++ {
		data := fmt.Sprintf("Data from service1: %d", i)
		fmt.Println("Service1: Sending data:", data)
		sharedChannel <- data
		time.Sleep(time.Second) // Simulate some work
	}
	close(sharedChannel) // Signal the end of data transmission
}

// Service 2: The consumer
func service2() {
	for {
		select {
		case data, ok := <-sharedChannel:
			if !ok {
				fmt.Println("Service2: Channel closed. Exiting.")
				return
			}
			fmt.Println("Service2: Received data:", data)
		default:
			time.Sleep(10 * time.Millisecond)
		}
	}
}

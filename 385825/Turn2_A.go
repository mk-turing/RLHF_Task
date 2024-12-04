//package main
//
//import (
//	"fmt"
//	"time"
//)
//
//const bufferSize = 10 // Buffer size to prevent sender from blocking immediately
//
//// Shared channel between service1 and service2 with a buffer
//var sharedChannel = make(chan string, bufferSize)
//
//func main() {
//	// Start service2 (consumer) in a separate goroutine
//	go service2()
//
//	// Start service1 (producer)
//	service1()
//}
//
//// Service 1: The producer
//func service1() {
//	for i := 0; i < 5; i++ {
//		data := fmt.Sprintf("Data from service1: %d", i)
//		fmt.Println("Service1: Sending data:", data)
//
//		select {
//		case sharedChannel <- data:
//			// Data sent successfully
//		case <-time.After(1 * time.Second):
//			// Failed to send data within the timeout period, handle the error
//			fmt.Println("Service1: Error sending data: Channel is full or receiver is slow.")
//		}
//
//		time.Sleep(time.Second) // Simulate some work
//	}
//	close(sharedChannel) // Signal the end of data transmission
//}
//
//// Service 2: The consumer
//func service2() {
//	for {
//		select {
//		case data, ok := <-sharedChannel:
//			if !ok {
//				fmt.Println("Service2: Channel closed. Exiting.")
//				return
//			}
//			time.Sleep(500 * time.Millisecond) // Simulate some work by sleeping
//			fmt.Println("Service2: Received data:", data)
//		default:
//			time.Sleep(10 * time.Millisecond)
//		}
//	}
//}

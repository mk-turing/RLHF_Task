//package main
//
//import (
//	"fmt"
//	"time"
//)
//
//// sharedChannel is a global channel for communication between microservices
//var sharedChannel = make(chan string)
//
//func senderMicroservice() {
//	for {
//		message := fmt.Sprintf("Hello from Sender! Time: %s", time.Now())
//		fmt.Println("Sender: Sending message:", message)
//		sharedChannel <- message
//		time.Sleep(2 * time.Second)
//	}
//}
//
//func receiverMicroservice() {
//	for {
//		message := <-sharedChannel
//		fmt.Println("Receiver: Received message:", message)
//	}
//}
//
//func main() {
//	go senderMicroservice()
//	go receiverMicroservice()
//	select {} // Main thread blocks indefinitely to keep the program running
//}

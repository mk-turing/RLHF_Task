package main

import (
	"log"
	"time"
)

// Callback function
type Callback func(int)

// Function that accepts a callback function and calls it after a random delay
func asyncOperation(callback Callback) {
	time.Sleep(time.Duration(1+int(time.Now().UnixNano()%3)) * time.Second)
	callback(10)
}
func main() {
	// Create a logger with file output
	logger := log.New(nil, "", log.LstdFlags|log.Lmicroseconds)
	logger.Println("Start of the application.")
	// Callback function to be invoked asynchronously
	callback := func(result int) {
		logger.Println("Callback executed:", result)
	}
	// Call the asyncOperation with the callback
	asyncOperation(callback)
	logger.Println("End of the application.")
}

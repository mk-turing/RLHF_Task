package main

import (
	"fmt"
	"log"
	"net"
	"time"
)

func sendRequest(address string, message string) error {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return fmt.Errorf("failed to connect to server: %w", err)
	}
	defer conn.Close()

	// Send message
	_, err = conn.Write([]byte(message))
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	// Read response
	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	fmt.Printf("Server response: %s\n", string(buffer[:n]))
	return nil
}

func main() {
	serverAddress := "localhost:8080"
	message := "Hello, server!"

	var err error
	for i := 0; i < 5; i++ { // Retry mechanism
		err = sendRequest(serverAddress, message)
		if err == nil {
			break // Success
		}

		log.Printf("Error occurred: %v. Retrying in %d seconds...", err, i+1)
		time.Sleep(time.Duration((i + 1)) * time.Second) // Exponential backoff
	}

	if err != nil {
		log.Fatalf("Failed to send request after retries: %v", err)
	}
}

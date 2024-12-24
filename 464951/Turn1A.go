package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"time"
)

func connectToServer(ctx context.Context, address string) error {
	// Create a new context with a timeout
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	conn, err := net.Dial("tcp", address)
	if err != nil {
		log.Printf("Failed to connect to %s: %v\n", address, err)
		return err
	}
	defer conn.Close() // Ensure the connection is closed regardless of the outcome

	// Simulate some activity on the socket
	select {
	case <-ctx.Done():
		log.Println("Context canceled, exiting.")
		return ctx.Err()
	default:
		// Read some data from the server
		buf := make([]byte, 1024)
		n, err := conn.Read(buf)
		if err != nil {
			log.Printf("Error reading from socket: %v\n", err)
			return err
		}
		log.Printf("Received %d bytes: %s\n", n, string(buf[:n]))

		// Send some data to the server
		msg := "Hello from the client!"
		_, err = conn.Write([]byte(msg))
		if err != nil {
			log.Printf("Error writing to socket: %v\n", err)
			return err
		}
		log.Println("Sent message to server.")
	}

	return nil
}

func main() {
	// Create a context with a parent context
	parentCtx, cancelParent := context.WithCancel(context.Background())
	defer cancelParent()

	// Example usage
	if err := connectToServer(parentCtx, "localhost:9090"); err != nil {
		fmt.Println("An error occurred:", err)
		os.Exit(1)
	}

	fmt.Println("Connection successful.")
}

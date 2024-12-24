package main

import (
	"context"
	"fmt"
	"net"
	"sync"
	"time"
)

func handleConnection(ctx context.Context, conn net.Conn) {
	defer conn.Close() // Ensure connection is closed regardless of errors
	defer fmt.Println("Connection closed.")

	var wg sync.WaitGroup
	wg.Add(2) // Two goroutines for reading and writing

	// Ensure we wait for both read and write operations to complete
	defer wg.Wait()

	// Handle reading from the connection in a goroutine
	go func() {
		defer wg.Done()
		for {
			select {
			case <-ctx.Done():
				fmt.Println("Reading from connection interrupted, closing...")
				return
			default:
				buf := make([]byte, 1024)
				n, err := conn.Read(buf)
				if err != nil {
					logError("Error reading from connection:", err)
					return
				}
				if n == 0 {
					fmt.Println("Connection closed by client")
					return
				}
				fmt.Println("Received data:", string(buf[:n]))
			}
		}
	}()

	// Handle writing to the connection (echo the received data) in a goroutine
	go func() {
		defer wg.Done()
		for {
			select {
			case <-ctx.Done():
				fmt.Println("Writing to connection interrupted, closing...")
				return
			default:
				buf := make([]byte, 1024)
				n, err := conn.Read(buf)
				if err != nil {
					logError("Error reading from connection for echo:", err)
					return
				}
				if n == 0 {
					fmt.Println("Connection closed by client for echo")
					return
				}
				_, err = conn.Write(buf[:n])
				if err != nil {
					logError("Error writing to connection:", err)
					return
				}
			}
		}
	}()
}

func main() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		logError("Error listening:", err)
		return
	}
	defer listener.Close()

	fmt.Println("Server listening on :8080")
	for {
		conn, err := listener.Accept()
		if err != nil {
			logError("Error accepting connection:", err)
			continue
		}

		// Create a context with a timeout of 5 seconds for each connection
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		go func() {
			handleConnection(ctx, conn)
		}()
	}
}

// Log an error and return it
func logError(message string, err error) error {
	fmt.Println(message, err)
	return err
}

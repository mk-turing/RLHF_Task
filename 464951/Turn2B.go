package main

import (
	"context"
	"fmt"
	"net"
	"sync"
	"time"
)

func handleConnection(ctx context.Context, conn net.Conn) {
	// Use a WaitGroup to ensure we wait for all goroutines to finish before closing the connection
	var wg sync.WaitGroup
	defer wg.Wait()
	defer conn.Close()

	// Handle reading from the connection in a goroutine
	wg.Add(1)
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
					fmt.Println("Error reading from connection:", err)
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

	// Handle writing to the connection (in this example, we just echo the received data)
	wg.Add(1)
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
					fmt.Println("Error reading from connection for echo:", err)
					return
				}
				if n == 0 {
					fmt.Println("Connection closed by client for echo")
					return
				}
				_, err = conn.Write(buf[:n])
				if err != nil {
					fmt.Println("Error writing to connection:", err)
					return
				}
			}
		}
	}()
}
func main() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Error listening:", err)
		return
	}
	defer listener.Close()

	fmt.Println("Server listening on :8080")
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}
		go func(conn net.Conn) {
			// Create a context with a timeout of 5 seconds for each connection
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			handleConnection(ctx, conn)
		}(conn)
	}
}

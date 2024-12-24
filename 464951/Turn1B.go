package main

import (
	"context"
	"fmt"
	"net"
	"time"
)

func handleConnection(ctx context.Context, conn net.Conn) {
	defer conn.Close() // Defer statement to ensure the connection is closed after the function returns

	for {
		select {
		case <-ctx.Done():
			fmt.Println("Handling connection interrupted, closing...")
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
}

func main() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Error listening:", err)
		return
	}
	defer listener.Close() // Defer statement to ensure the listener is closed after the main function returns

	fmt.Println("Server listening on :8080")
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		// Create a context with a timeout of 5 seconds
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel() // Defer statement to ensure the context is canceled after the function returns

		// Handle the connection in a goroutine
		go handleConnection(ctx, conn)
	}
}

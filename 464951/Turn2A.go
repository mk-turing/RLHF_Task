package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"
)

func handleConnection(ctx context.Context, conn net.Conn) {
	defer conn.Close() // Ensures the connection is closed on return

	for {
		select {
		case <-ctx.Done():
			fmt.Println("Connection interrupted, closing...")
			return
		default:
			buf := make([]byte, 1024)
			n, err := conn.Read(buf)
			if err != nil {
				log.Printf("Error reading from connection: %v", err)
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
		log.Fatalf("Error listening: %v", err)
	}
	defer listener.Close() // Ensures the listener is closed on program exit

	fmt.Println("Server listening on :8080")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Error accepting connection: %v", err)
			continue
		}

		// Create a context with a timeout of 5 seconds for each connection
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel() // Ensures the context is canceled on return

		// Start a new goroutine to handle the connection
		go handleConnection(ctx, conn)
	}
}

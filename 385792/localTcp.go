package main

import (
	"fmt"
	"log"
	"net"
	"os"
)

func handleTCPConnection(conn net.Conn) {
	defer conn.Close()
	fmt.Println("Client connected:", conn.RemoteAddr())

	// Send a welcome message
	_, err := conn.Write([]byte("Hello from server!\n"))
	if err != nil {
		log.Println("Error writing to client:", err)
		return
	}

	// Read incoming data from the client
	buf := make([]byte, 1024)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			log.Println("Error reading from client:", err)
			break
		}

		// Print the received data
		fmt.Printf("Received from client: %s\n", string(buf[:n]))
	}
}

func main() {
	// Listen on TCP port 8080
	listen, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("Error starting server: %v\n", err)
		os.Exit(1)
	}
	defer listen.Close()

	fmt.Println("Server is listening on port 8080...")

	for {
		// Accept a new client connection
		conn, err := listen.Accept()
		if err != nil {
			log.Println("Error accepting connection:", err)
			continue
		}

		// Handle the connection
		go handleTCPConnection(conn)
	}
}

package main

import (
	"io"
	"log"
	"net"
)

// TCP server function to handle incoming client connections
func startTCPServer() {
	// Start listening on localhost:9090
	listen, err := net.Listen("tcp", "localhost:9090")
	if err != nil {
		log.Fatalf("Error starting server: %v\n", err)
	}
	defer listen.Close()

	log.Println("Server is listening on localhost:9090")

	// Accept incoming connections
	for {
		conn, err := listen.Accept()
		if err != nil {
			log.Printf("Error accepting connection: %v\n", err)
			continue
		}

		// Handle each connection in a new goroutine
		go handleConnection(conn)
	}
}

// Handle client connections
func handleConnection(conn net.Conn) {
	defer conn.Close()
	log.Println("New client connected:", conn.RemoteAddr())

	// Read incoming data from the client
	buf := make([]byte, 1024)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			if err != io.EOF {
				log.Printf("Error reading from connection: %v\n", err)
			}
			break
		}

		log.Printf("Received %d bytes: %s\n", n, string(buf[:n]))

		// Send a response to the client
		msg := "Hello from the server!"
		_, err = conn.Write([]byte(msg))
		if err != nil {
			log.Printf("Error writing to connection: %v\n", err)
			break
		}
		log.Println("Sent message to client.")
	}
}

func main() {
	startTCPServer()
}

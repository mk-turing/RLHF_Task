package main

import (
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var (
	clients   = make(map[*websocket.Conn]bool) // Map of connected clients
	broadcast = make(chan string)              // Channel for broadcasting messages
	wg        sync.WaitGroup                   // WaitGroup to ensure graceful shutdown
)

// WebSocket connection handler
func handleConnections(w http.ResponseWriter, r *http.Request) {
	// Upgrade the HTTP connection to WebSocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

	// Register the new client
	clients[conn] = true
	log.Printf("New client connected. Total clients: %d\n", len(clients))

	defer func() {
		// Unregister the client when the connection is closed
		delete(clients, conn)
		log.Printf("Client disconnected. Total clients: %d\n", len(clients))
		// Signal that the client's goroutine is done
		wg.Done()
	}()

	// Start a goroutine to handle message sending for this client
	wg.Add(1)
	go handleClientMessages(conn)

	// Read messages from the client and broadcast them
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}

		// Broadcast the message to all connected clients
		broadcast <- string(message)
	}
}

// Function to handle broadcasting messages to a single client concurrently
func handleClientMessages(conn *websocket.Conn) {
	for message := range broadcast {
		err := conn.WriteMessage(websocket.TextMessage, []byte(message))
		if err != nil {
			log.Println(err)
			conn.Close()
			break
		}
	}
}

func handleBroadcast() {
	for {
		message := <-broadcast
		// Send the message to all connected clients concurrently
		for client := range clients {
			// Use a buffered channel to send messages to each client's goroutine
			go func(c *websocket.Conn) {
				log.Fatalf("", message)
			}(client)
		}
	}
}

func main() {
	// Start the broadcast handler
	go handleBroadcast()

	// Serve the WebSocket endpoint
	http.HandleFunc("/ws", handleConnections)
	log.Println("WebSocket server is running on port 8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}

	// Wait for all client goroutines to finish before exiting
	wg.Wait()
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins for simplicity
	},
}

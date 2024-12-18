package main

import (
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var (
	clients   = make(map[*websocket.Conn]bool)
	lock      = &sync.Mutex{}     // Synchronization to ensure thread safety
	buffer    = make([]string, 0) // Buffer to maintain message order
	broadcast = make(chan bool)   // Signal to request message broadcasting
)

func handleConnections(w http.ResponseWriter, r *http.Request) {
	// Upgrade the HTTP connection to WebSocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

	// Register the new client
	lock.Lock()
	clients[conn] = true
	lock.Unlock()

	log.Printf("New client connected. Total clients: %d\n", len(clients))

	defer func() {
		// Unregister the client when the connection is closed
		lock.Lock()
		delete(clients, conn)
		lock.Unlock()

		log.Printf("Client disconnected. Total clients: %d\n", len(clients))
	}()

	// Handle messages from the client
	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}

		lock.Lock()
		buffer = append(buffer, string(message))
		lock.Unlock()

		select {
		case broadcast <- true: // Request message broadcast
		default:
		}
	}
}

func handleBroadcast() {
	lock.Lock()
	deferredLock := lock
	deferredBuffer := buffer
	buffer = make([]string, 0)
	lock.Unlock()

	for i := range deferredBuffer {
		clientMessages := make(chan string, len(clients))

		// Broadcast to each client in a goroutine
		for client := range clients {
			go func(client *websocket.Conn, message string) {
				select {
				case clientMessages <- message:
					defer close(clientMessages) // Close the channel after sending
				default:
				}
			}(client, deferredBuffer[i])
		}

		// Wait for all messages to be sent to maintain order
		go func() {
			for range clientMessages {
				// Increment to ensure sequential operation
				if i == len(deferredBuffer)-1 {
					broadcast <- true // Trigger the next batch after the last message is sent
				}
			}
		}()
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
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins for simplicity
	},
}

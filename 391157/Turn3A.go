package main

import (
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

var (
	clients         = make(map[*websocket.Conn]bool) // Map of connected clients
	broadcastChan   = make(chan string)              // Channel for broadcasting messages
	broadcastMutex  = sync.Mutex{}                   // Mutex to ensure order of messages
	broadcastQueue  = []string{}                     // Queue to hold messages before sending
	broadcastTicker = time.NewTicker(50 * time.Millisecond)
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
	}()

	// Read messages from the client and broadcast them
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}

		// Broadcast the message to all connected clients
		broadcastChan <- string(message)
	}
}

// Function to handle broadcasting messages to all connected clients
func handleBroadcast() {
	for {
		message := <-broadcastChan

		// Queue the message to ensure order of delivery
		broadcastMutex.Lock()
		broadcastQueue = append(broadcastQueue, message)
		broadcastMutex.Unlock()

		// Send messages from the queue to clients
		select {
		case <-broadcastTicker.C:
			broadcastMutex.Lock()
			for len(broadcastQueue) > 0 {
				msg := broadcastQueue[0]
				broadcastQueue = broadcastQueue[1:]
				broadcastMutex.Unlock()

				// Send the message to all connected clients
				for client := range clients {
					err := client.WriteMessage(websocket.TextMessage, []byte(msg))
					if err != nil {
						log.Println(err)
						client.Close()
						delete(clients, client)
					}
				}

				// Lock the mutex back for the next message
				broadcastMutex.Lock()
			}
			broadcastMutex.Unlock()
		default:
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
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins for simplicity
	},
}

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
	broadcastQueue  = sync.Map{}                     // Priority queue (map of priorities to slices of messages)
	broadcastTicker = time.NewTicker(50 * time.Millisecond)
	maxQueueSize    = 1000 // Limit the queue size
	cond            = sync.NewCond(&broadcastMutex)
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

		// Add the message to the priority queue
		priority := 1 // Assign default priority
		if message == "URGENT" {
			priority = 0 // Assign higher priority for urgent messages
		}

		broadcastMutex.Lock()
		// Check if there are existing messages for the same priority
		if msgs, ok := broadcastQueue.Load(priority); ok {
			msgSlice := msgs.([]string)
			msgSlice = append(msgSlice, message)
			broadcastQueue.Store(priority, msgSlice)
		} else {
			// If no existing messages, create a new slice for the priority
			broadcastQueue.Store(priority, []string{message})
		}

		// Check if the queue has reached the maximum size
		count := 0
		broadcastQueue.Range(func(key, value interface{}) bool {
			count++
			return true
		})

		if count > maxQueueSize {
			log.Printf("Queue size limit reached, removing oldest messages")
			// Delete the oldest entry (first priority level)
			broadcastQueue.Range(func(key, value interface{}) bool {
				broadcastQueue.Delete(key)
				return false // Only delete one entry
			})
		}

		// Notify all goroutines that a new message has been added
		cond.Signal()
		broadcastMutex.Unlock()
	}
}

// Function to send messages from the queue to clients
func sendMessages() {
	for {
		broadcastMutex.Lock()
		broadcastQueue.Range(func(key, value interface{}) bool {
			msgSlice := value.([]string)
			for _, msg := range msgSlice {
				// Remove the message after sending
				broadcastQueue.Delete(key)

				// Send the message to all connected clients
				for client := range clients {
					err := client.WriteMessage(websocket.TextMessage, []byte(msg))
					if err != nil {
						log.Println(err)
						client.Close()
						delete(clients, client)
					}
				}
			}
			return false // Stop after the first priority level is processed
		})

		cond.Wait() // Wait for more messages to be added
		broadcastMutex.Unlock()
	}
}

func main() {
	// Start the broadcast handler
	go handleBroadcast()
	go sendMessages()

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

package main

import (
	"container/heap"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

const (
	maxQueueSize = 1000 // Maximum size of the broadcast queue
	lowPriority  = 1
	highPriority = 5 // Higher priority value for urgent updates
)

// message struct with priority and data
type message struct {
	data     string
	priority int
}

// messageHeap implements heap.Interface for priority queue
type messageHeap []*message

func (h messageHeap) Len() int           { return len(h) }
func (h messageHeap) Less(i, j int) bool { return h[i].priority > h[j].priority }
func (h messageHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h *messageHeap) Push(x interface{}) {
	*h = append(*h, x.(*message))
}
func (h *messageHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

var (
	clients         = make(map[*websocket.Conn]bool)
	broadcastChan   = make(chan *message)
	broadcastMutex  = sync.Mutex{}
	broadcastQueue  messageHeap
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
		// Assign priority based on message type (you can implement a more complex logic here)
		priority := lowPriority
		if string(message) == "urgent" {
			priority = highPriority
		}
		// Broadcast the message to all connected clients
		broadcastChan <- &message{string(message), priority}
	}
}

// Function to handle broadcasting messages to all connected clients
func handleBroadcast() {
	for {
		msg := <-broadcastChan
		// Queue the message to ensure order of delivery
		broadcastMutex.Lock()
		if len(broadcastQueue) >= maxQueueSize {
			// Queue is full, discard the message or take appropriate action
			broadcastMutex.Unlock()
			log.Println("Broadcast queue is full, discarding message.")
			continue
		}
		heap.Push(&broadcastQueue, msg)
		broadcastMutex.Unlock()
package main

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"sync"
)

var upgrader = websocket.Upgrader{} // Use default upgrader

// Define a struct for chat messages
type Message struct {
	User    string `json:"user"`
	Content string `json:"content"`
}

// A Hub maintains the set of active connections and broadcasts messages to clients.
type Hub struct {
	clients   map[*websocket.Conn]bool
	broadcast chan Message
	mu        sync.Mutex
}

// NewHub creates a new Hub
func NewHub() *Hub {
	return &Hub{
		clients:   make(map[*websocket.Conn]bool),
		broadcast: make(chan Message),
	}
}

// Run starts the Hub to listen for incoming messages and broadcast them to all clients
func (h *Hub) Run() {
	for {
		msg := <-h.broadcast
		h.mu.Lock()
		for client := range h.clients {
			if err := client.WriteJSON(msg); err != nil {
				log.Println("Error writing message to client:", err)
				client.Close()
				delete(h.clients, client)
			}
		}
		h.mu.Unlock()
	}
}

// ChatHandler handles WebSocket connections and communication
func chatHandler(hub *Hub, w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	claims := &jwt.StandardClaims{}
	_, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte("your_secret_key"), nil
	})
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "Could not upgrade connection", http.StatusInternalServerError)
		return
	}
	defer conn.Close()

	hub.mu.Lock()
	hub.clients[conn] = true
	hub.mu.Unlock()

	// Listen for incoming messages
	for {
		var msg Message
		err := conn.ReadJSON(&msg)
		if err != nil {
			break // Exit on error
		}
		hub.broadcast <- msg // Broadcast the message to all clients
	}

	// Remove client when done
	hub.mu.Lock()
	delete(hub.clients, conn)
	hub.mu.Unlock()
}

func main() {
	hub := NewHub()
	go hub.Run()

	http.HandleFunc("/chat", func(w http.ResponseWriter, r *http.Request) {
		chatHandler(hub, w, r)
	})

	log.Println("Server running on :8080")
	http.ListenAndServe(":8080", nil)
}

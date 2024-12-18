package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"sync"
	"time"
)

type InventoryUpdate struct {
	ProductID string `json:"product_id"`
	Quantity  int    `json:"quantity"`
}

type User struct {
	ID       string
	Role     string
	Token    string
	Expires  time.Time
	Conn     *websocket.Conn
	LastPing time.Time
}

var (
	upGrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true // Allow all origins; customize as needed
		},
	}
	inventoryUpdates = make(chan InventoryUpdate)
	clients          = make(map[*websocket.Conn]*User) // Map to track connected clients with their user data
	clientsMutex     = sync.Mutex{}                    // Mutex for safe access to the clients map
	pingInterval     = time.Second * 10                // Send a ping every 10 seconds
	pongInterval     = time.Second * 5                 // Expect a pong response within 5 seconds
)

func authenticate(token string) (*User, error) {
	// Placeholder authentication logic. In a real app, verify JWT and check user database.
	if token == "admin:12345" { // For demonstration purposes only
		return &User{
			ID:      "admin",
			Role:    "admin",
			Token:   token,
			Expires: time.Now().Add(time.Hour),
		}, nil
	}
	return nil, fmt.Errorf("invalid token")
}

func main() {
	go handleInventoryUpdates()
	go handlePings()
	http.HandleFunc("/ws", handleConnections)
	fmt.Println("http server started on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	if token == "" {
		http.Error(w, "Authorization header required", http.StatusUnauthorized)
		return
	}

	user, err := authenticate(token)
	if err != nil {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	conn, err := upGrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

	user.Conn = conn
	user.LastPing = time.Now()

	// Add the new connection to the clients map
	clientsMutex.Lock()
	clients[conn] = user
	clientsMutex.Unlock()

	defer func() {
		// Remove the connection when it closes
		clientsMutex.Lock()
		delete(clients, conn)
		clientsMutex.Unlock()
	}()

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		var update InventoryUpdate
		// Parse the incoming message
		if err := json.Unmarshal(message, &update); err != nil {
			log.Println("json unmarshal:", err)
			continue
		}

		// Check if the user has permission to update inventory
		if user.Role != "admin" {
			conn.WriteMessage(websocket.TextMessage, []byte("Not authorized to update inventory"))
			continue
		}

		// Send the update to the inventoryUpdates channel
		inventoryUpdates <- update
	}
}

func handleInventoryUpdates() {
	for update := range inventoryUpdates {
		// Broadcast the update to all connected clients
		clientsMutex.Lock()
		for conn, user := range clients {
			if user.Role != "admin" { // Only broadcast to non-admin clients
				if err := conn.WriteJSON(update); err != nil {
					// If the client has disconnected, remove it from the map
					log.Printf("Error writing to client: %v. Removing client.", err)
					delete(clients, conn)
					conn.Close()
				}
			}
		}
		clientsMutex.Unlock()
	}
}

func handlePings() {
	ticker := time.NewTicker(pingInterval)
	defer ticker.Stop()

	for range ticker.C {
		clientsMutex.Lock()
		for conn, user := range clients {
			if time.Since(user.LastPing) > pongInterval {
				log.Printf("Closing connection to inactive client %v", user.ID)
				conn.Close()
				delete(clients, conn)
			} else {
				err := conn.WriteMessage(websocket.PingMessage, nil)
				if err != nil {
					log.Printf("Error sending ping to client %v: %v", user.ID, err)
					conn.Close()
					delete(clients, conn)
				} else {
					user.LastPing = time.Now()
				}
			}
		}
		clientsMutex.Unlock()
	}
}

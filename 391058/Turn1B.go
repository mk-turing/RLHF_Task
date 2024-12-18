package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"sync"
)

type InventoryUpdate struct {
	ProductID string `json:"product_id"`
	Quantity  int    `json:"quantity"`
}

var (
	upGrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true // Allow all origins; customize as needed
		},
	}
	inventoryUpdates = make(chan InventoryUpdate)
	clients          = make(map[*websocket.Conn]bool) // Map to track connected clients
	clientsMutex     = sync.Mutex{}                   // Mutex for safe access to the clients map
)

func main() {
	go handleInventoryUpdates()
	http.HandleFunc("/ws", handleConnections)
	fmt.Println("http server started on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	conn, err := upGrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

	// Add the new connection to the clients map
	clientsMutex.Lock()
	clients[conn] = true
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
		// Send the update to the inventoryUpdates channel
		inventoryUpdates <- update
	}
}

func handleInventoryUpdates() {
	for update := range inventoryUpdates {
		// Broadcast the update to all connected clients
		clientsMutex.Lock()
		for conn := range clients {
			if err := conn.WriteJSON(update); err != nil {
				// If the client has disconnected, remove it from the map
				log.Printf("Error writing to client: %v. Removing client.", err)
				delete(clients, conn)
				conn.Close()
			}
		}
		clientsMutex.Unlock()
	}
}

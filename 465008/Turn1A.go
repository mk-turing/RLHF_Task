package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

// Constants to define the WebSocket endpoint and max message size
const (
	upgradeHeaderKey   = "Upgrade"
	upgradeHeaderValue = "websocket"
	connectHeaderKey   = "Connection"
	connectHeaderValue = "upgrade"
	wsProtocol         = "ws"
	maxMessageSize     = 1024
)

var (
	// Define a WebSocket upgradable handler
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true // Allow connections from any origin for simplicity
		},
	}

	// A channel to broadcast messages to all connected clients
	broadcast = make(chan string)
)

// wsUpgradeHeaders now uses http.Header type instead of map[string]string
var wsUpgradeHeaders = http.Header{
	upgradeHeaderKey: {upgradeHeaderValue},
	connectHeaderKey: {connectHeaderValue},
}

// Handle the WebSocket connections
func websocketHandler(w http.ResponseWriter, r *http.Request) {
	// Upgrade the HTTP connection to a WebSocket connection
	conn, err := upgrader.Upgrade(w, r, wsUpgradeHeaders)
	if err != nil {
		log.Println("upgrade:", err)
		return
	}
	defer conn.Close()

	// Join the broadcast channel
	go sendBroadcast(conn)

	// Handle incoming messages
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Println("read:", err)
			}
			return
		}

		// Log the incoming message
		log.Printf("Received: %s\n", message)

		// Broadcast the formatted message to all clients
		formattedMessage := fmt.Sprintf("<div>%s</div>", message)
		broadcast <- formattedMessage
	}
}

// Send broadcast messages to all connected clients
func sendBroadcast(conn *websocket.Conn) {
	for message := range broadcast {
		if err := conn.WriteMessage(websocket.TextMessage, []byte(message)); err != nil {
			log.Println("write:", err)
			return
		}
	}
}

func main() {
	// Register the WebSocket handler on the "/ws" endpoint
	http.HandleFunc("/ws", websocketHandler)

	// Start the server
	log.Println("Starting websocket server on http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

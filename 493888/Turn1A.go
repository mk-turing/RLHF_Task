package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

type Document struct {
	Content string `json:"content"`
}

type DocumentChange struct {
	Type    string `json:"type"`
	Index   int    `json:"index"`
	Content string `json:"content"`
}

var (
	mu          sync.RWMutex
	doc         *Document                = &Document{Content: ""}
	connections map[*websocket.Conn]bool = make(map[*websocket.Conn]bool)
)

func updateDocument(change DocumentChange) {
	mu.Lock()
	defer mu.Unlock()

	switch change.Type {
	case "insert":
		doc.Content = doc.Content[:change.Index] + change.Content + doc.Content[change.Index:]
	case "delete":
		doc.Content = doc.Content[:change.Index] + doc.Content[change.Index+len(change.Content):]
	default:
		log.Printf("Unsupported change type: %v", change.Type)
	}
}

func sendChangeToAll(change DocumentChange) {
	bytes, err := json.Marshal(change)
	if err != nil {
		log.Printf("Error marshalling change: %v", err)
		return
	}

	for conn := range connections {
		err := conn.WriteMessage(websocket.TextMessage, bytes)
		if err != nil {
			log.Printf("Error writing message: %v", err)
			delete(connections, conn)
			conn.Close()
		}
	}
}

func handleConnection(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Error upgrading connection: %v", err)
		return
	}
	defer conn.Close()

	connections[conn] = true
	defer delete(connections, conn)

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Error reading message: %v", err)
			break
		}

		var change DocumentChange
		err = json.Unmarshal(message, &change)
		if err != nil {
			log.Printf("Error unmarshalling change: %v", err)
			continue
		}

		updateDocument(change)
		sendChangeToAll(change)
	}
}

func main() {
	http.HandleFunc("/ws", handleConnection)
	log.Println("Server starting on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

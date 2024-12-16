package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Message struct {
	Type     string      `json:"type"`
	Data     interface{} `json:"data"`
	Username string      `json:"username"`
}

func websocketHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error upgrading connection:", err)
		return
	}
	defer conn.Close()

	users := make(map[string]*websocket.Conn)

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error reading message:", err)
			break
		}

		var incomingMessage Message
		if err := json.Unmarshal(msg, &incomingMessage); err != nil {
			log.Println("Error unmarshaling message:", err)
			continue
		}

		switch incomingMessage.Type {
		case "join":
			users[incomingMessage.Username] = conn
			fmt.Println(incomingMessage.Username, "joined the chat")
			broadcastMessage(users, &Message{
				Type:     "join",
				Data:     incomingMessage.Username,
				Username: "server",
			})
		case "message":
			broadcastMessage(users, &Message{
				Type:     "message",
				Data:     incomingMessage.Data,
				Username: incomingMessage.Username,
			})
		default:
			log.Println("Unknown message type:", incomingMessage.Type)
		}
	}

	delete(users, conn.RemoteAddr().String())
}

func broadcastMessage(users map[string]*websocket.Conn, message *Message) {
	jsonMessage, err := json.Marshal(message)
	if err != nil {
		log.Println("Error marshaling broadcast message:", err)
		return
	}

	for _, userConn := range users {
		if err := userConn.WriteMessage(websocket.TextMessage, jsonMessage); err != nil {
			log.Println("Error writing broadcast message:", err)
		}
	}
}

func main() {
	http.HandleFunc("/ws", websocketHandler)
	log.Println("Listening on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("Error starting server:", err)
	}
}

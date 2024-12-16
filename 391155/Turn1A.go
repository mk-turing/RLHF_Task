package main

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"sync"
)

// Message struct to define the structure of messages
type Message struct {
	Type    string `json:"type"`
	Content string `json:"content"`
}

// Server holds the active connections
type Server struct {
	clients   map[*websocket.Conn]bool
	broadcast chan Message
	mu        sync.Mutex
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (s *Server) handleConnections(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

	s.mu.Lock()
	s.clients[conn] = true
	s.mu.Unlock()

	for {
		var msg Message
		err := conn.ReadJSON(&msg)
		if err != nil {
			log.Println(err)
			break
		}
		s.broadcast <- msg
	}

	s.mu.Lock()
	delete(s.clients, conn)
	s.mu.Unlock()
}

func (s *Server) handleMessages() {
	for {
		msg := <-s.broadcast

		s.mu.Lock()
		for client := range s.clients {
			err := client.WriteJSON(msg)
			if err != nil {
				log.Println(err)
				client.Close()
				delete(s.clients, client)
			}
		}
		s.mu.Unlock()
	}
}

func main() {
	server := &Server{
		clients:   make(map[*websocket.Conn]bool),
		broadcast: make(chan Message),
	}

	go server.handleMessages()

	http.HandleFunc("/ws", server.handleConnections)
	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

package main

import (
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/websocket"
)

var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	connections = make(map[*websocket.Conn]chan string)
	messageCh   chan string
)

func startWebSocketServer(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}
	defer conn.Close()

	connections[conn] = make(chan string)
	messageCh <- strings.Join([]string{"new client connected: ", conn.RemoteAddr().String()}, "")

	go handleMessages(conn)
	go handleConnecionClose(conn)
}

func handleMessages(conn *websocket.Conn) {
	defer close(connections[conn])

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error reading message:", err)
			return
		}

		messageCh <- strings.Join([]string{"receive from ", conn.RemoteAddr().String(), ": ", string(message)}, "")

		for otherConn := range connections {
			if otherConn != conn {
				err := otherConn.WriteMessage(websocket.TextMessage, message)
				if err != nil {
					log.Println("Error writing message to other clients:", err)
					return
				}
			}
		}
	}
}

func handleConnecionClose(conn *websocket.Conn) {
	select {
	case _ = <-connections[conn]:
		messageCh <- strings.Join([]string{"client disconnected: ", conn.RemoteAddr().String()}, "")
		delete(connections, conn)
	case msg := <-messageCh:
		err := conn.WriteMessage(websocket.TextMessage, []byte(msg))
		if err != nil {
			log.Println("Error writing disconnection message to client:", err)
		}
	}
}

func main() {
	http.HandleFunc("/ws", startWebSocketServer)

	log.Println("Starting WebSocket server on :8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("ListenAndServe error:", err)
	}
}

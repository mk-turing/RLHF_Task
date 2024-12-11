package main

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"os"
	"sync"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Secret key for JWT signing
var jwtKey = []byte("your_secret_key")

// Manage WebSocket connections
type Chat struct {
	Clients map[*websocket.Conn]string // map websocket connection to username
	Send    chan Message               // Channel for sending messages
	Lock    sync.RWMutex
}

type Message struct {
	From    string `json:"from"`
	Message string `json:"message"`
}

func (c *Chat) Broadcast(msg Message) {
	c.Lock.RLock()
	defer c.Lock.RUnlock()

	for conn := range c.Clients {
		if err := conn.WriteJSON(msg); err != nil {
			log.Println("Error writing to client:", err)
			c.Remove(conn)
		}
	}
}

func (c *Chat) Add(conn *websocket.Conn, username string) {
	c.Lock.Lock()
	defer c.Lock.Unlock()

	c.Clients[conn] = username
	log.Println("Client", username, "connected")
}

func (c *Chat) Remove(conn *websocket.Conn) {
	c.Lock.Lock()
	defer c.Lock.Unlock()

	delete(c.Clients, conn)
	//close(conn)
	log.Println("Client disconnected")
}

func authenticate(tokenStr string) (bool, string, error) {
	claims := &jwt.StandardClaims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		return false, "", err
	}
	if !token.Valid {
		return false, "", jwt.ErrSignatureInvalid
	}
	// Here, typically you would extract the username from the token
	username := "anonymous"
	return true, username, nil
}

func chatHandler(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	isValid, username, err := authenticate(token)
	if !isValid || err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}
	defer conn.Close()

	chat.Add(conn, username)
	defer chat.Remove(conn)

	for {
		var message Message
		// Read message from client
		if err := conn.ReadJSON(&message); err != nil {
			log.Println("Read error:", err)
			break
		}

		message.From = username
		chat.Broadcast(message)
	}
}

func main() {
	// Ensure the JWT verification key is provided (you can load it from an environment variable or .env file)
	if _, exists := os.LookupEnv("JWT_SECRET"); !exists {
		log.Fatal("JWT_SECRET environment variable is required")
	}

	chat := &Chat{
		Clients: make(map[*websocket.Conn]string),
		Send:    make(chan Message),
	}

	go func() {
		for msg := range chat.Send {
			chat.Broadcast(msg)
		}
	}()

	http.HandleFunc("/chat", chatHandler)

	log.Println("Server running on :8080")
	http.ListenAndServe(":8080", nil)
}

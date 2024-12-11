package main

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

var upgrader = websocket.Upgrader{}

// Secret key for JWT signing
var jwtKey = []byte("your_secret_key")

type Message struct {
	Message string `json:"message"`
}

func authenticate(tokenStr string) (bool, error) {
	claims := &jwt.StandardClaims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		return false, err
	}
	return token.Valid, nil
}

func chatHandler(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	isValid, err := authenticate(token)
	if !isValid || err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "Could not upgrade connection", http.StatusInternalServerError)
		return
	}
	defer conn.Close()

	// Welcome message
	welcome := Message{Message: "Welcome to the chat!"}
	conn.WriteJSON(welcome)

	for {
		var msg Message
		// Read message from client
		err := conn.ReadJSON(&msg)
		if err != nil {
			break
		}

		// Echo the message back
		if err := conn.WriteJSON(msg); err != nil {
			break
		}
	}
}

func main() {
	// Ensure the JWT verification key is provided (you can load it from an environment variable or .env file)
	//if _, exists := os.LookupEnv("JWT_SECRET"); !exists {
	//	log.Fatal("JWT_SECRET environment variable is required")
	//}

	http.HandleFunc("/chat", chatHandler)

	// Serve OpenAPI documentation
	http.HandleFunc("/openapi.yaml", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "turn3Aopenapi.yaml")
	})

	log.Println("Server running on :8080")
	http.ListenAndServe(":8080", nil)
}

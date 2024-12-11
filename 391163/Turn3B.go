package main

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		u, err := url.Parse(r.Header.Get("Origin"))
		if err != nil {
			return false
		}
		return u.Host == "localhost:8080"
	},
}

type Message struct {
	Message string `json:"message"`
}

var secretKey = []byte("your-secret-key-here") // Replace with a secure secret

func validateJWT(tokenString string) bool {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodECDSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secretKey, nil
	})

	return err == nil && token.Valid
}

func chatHandler(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, "Unauthorized: No Authorization header provided", http.StatusUnauthorized)
		return
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" {
		http.Error(w, "Unauthorized: Invalid Authorization header", http.StatusUnauthorized)
		return
	}

	if !validateJWT(parts[1]) {
		http.Error(w, "Unauthorized: Invalid JWT", http.StatusUnauthorized)
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
			break // Closed connection or error
		}

		// Echo the message back
		if err := conn.WriteJSON(msg); err != nil {
			break
		}
	}
}

func main() {
	http.HandleFunc("/chat", chatHandler)

	// Serve OpenAPI documentation
	http.HandleFunc("/openapi.yaml", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "turn3Bopenapi.yaml")
	})

	http.ListenAndServe(":8080", nil)
}

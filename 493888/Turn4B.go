package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/gomodule/redigo/redis"
	"github.com/gorilla/websocket"
	_ "github.com/lib/pq"
)

// Notification model
type Notification struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	Type      string    `json:"type"` // e.g., "update", "conflict"
	Content   string    `json:"content"`
	Timestamp time.Time `json:"timestamp"`
}

// Database connection pool
var db *sql.DB
var redisPool *redis.Pool

func initDatabase() error {
	dbConnString := "postgres://username:password@localhost/dbname?sslmode=disable"
	var err error
	db, err = sql.Open("postgres", dbConnString)
	if err != nil {
		return err
	}
	if err = db.Ping(); err != nil {
		return err
	}
	return nil
}

func initRedisPool() error {
	redisPool = &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", "localhost:6379")
		},
	}
	return nil
}

func generateToken(username string, role string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["username"] = username
	claims["role"] = role
	claims["exp"] = time.Now().Add(time.Hour).Unix() // Token expires in 1 hour

	return token.SignedString(privateKey)
}

func authenticateMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authorizationHeader := r.Header.Get("Authorization")
		if strings.HasPrefix(authorizationHeader, "Bearer ") {
			tokenString := authorizationHeader[7:]
			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
				}
				return privateKey, nil
			})
			if err != nil || !token.Valid {
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}
			claims := token.Claims.(jwt.MapClaims)
			username := claims["username"].(string)
			role := claims["role"].(string)
			fmt.Println("Authenticated user:", username, ", role:", role)
			next(w, r)
		} else {
			http.Error(w, "Missing authorization token", http.StatusUnauthorized)
		}
	}
}

func updateDocument(change DocumentChange, username string, role string) error {
	mu.Lock()
	defer mu.Unlock()

	if role != "editor" {
		return fmt.Errorf("User %s does not have permission to edit the document", username)
	}

	doc.Version++
	doc.Operations = append(doc.Operations, changeToOperation(change))

	return nil
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

		// Here, you should add more checks for specific sections if needed
		err = updateDocument(change, "user1", "editor") // Example user; replace with actual user
		if err != nil {
			log.Printf("Error updating document: %v", err)
			conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("Error: %s", err.Error())))
			continue
		}

		sendChangeToAll(change)
	}
}

func changeToOperation(change DocumentChange) string {
	switch change.Type {
	case "insert":
		return fmt.Sprintf("I%d%s", change.Index, change.Content)
	case "delete":
		return fmt.Sprintf("D%d", change.Index)
	default:
		return ""
	}
}

func main() {
	// Use authenticateMiddleware to wrap the handleConnection function
	http.HandleFunc("/ws", authenticateMiddleware(handleConnection))
	http.HandleFunc("/auth/login", authenticateLogin)

	log.Println("Server starting on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

func authenticateLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		return
	}

	var username string
	err := json.NewDecoder(r.Body).Decode(&username)
	if err != nil {
		http.Error(w, "Error decoding request", http.StatusBadRequest)
		return
	}

	role, ok := userRoles[username]
	if !ok {
		http.Error(w, "Username not found", http.StatusNotFound)
		return
	}

	tokenString, err := generateToken(username, role)
	if err != nil {
		http.Error(w, "Error generating token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(map[string]string{"token": tokenString}); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	}
}

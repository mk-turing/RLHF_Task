package main

import (
	"fmt"
	"github.com/go-redis/redis"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"sync"
	"time"
)

type InventoryUpdate struct {
	ProductID string `json:"product_id"`
	Quantity  int    `json:"quantity"`
}

type User struct {
	ID      string
	Role    string
	Token   string
	Expires time.Time
}

var (
	upGrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true // Allow all origins; customize as needed
		},
	}
	// Redis client
	rdb *redis.Client
	// Channel for inventory updates
	inventoryUpdates = make(chan InventoryUpdate)
	// Map to track connected clients
	clients      = make(map[*websocket.Conn]*User)
	clientsMutex = sync.Mutex{}
)

func authenticate(token string) (*User, error) {
	// Placeholder authentication logic. In a real app, verify JWT and check user database.
	// ...
	return &User{ID: "user", Role: "user"}, nil
}

func initRedis() {
	// Redis connection setup
	rdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // Redis server address
		Password: "",               // No password set
		DB:       0,                // Use default DB
	})

	_, err := rdb.Ping().Result()
	if err != nil {
		log.Fatal("Error connecting to Redis:", err)
	}
}

func main() {
	initRedis()
	go handleInventoryUpdates()
	go handleRedisSubscription()
	http.HandleFunc("/ws", handleConnections)
	fmt.Println("http server started on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	// ... (Same as before)

	// Subscribe the client to the inventory updates channel in Redis
	ps := rdb.PSubscribe("inventory_updates")
	defer ps.Close()

	for {
		msg, err := ps.ReceiveMessage()
		if err != nil {
			log.Println("PSubscribe:", err)
			break
		}
		if msg.Channel == "inventory_updates" {
			conn.WriteMessage(websocket.TextMessage, []byte(msg.Payload))
		}
	}
}

func handleInventoryUpdates() {
	for update := range inventoryUpdates {
		// Broadcast the update to all connected clients using Redis
		rdb.Publish("inventory_updates", string(update)).Result()
	}
}

func handleRedisSubscription() {
	// Handle incoming inventory updates from clients and publish them to the Redis channel
	for {
		update := <-inventoryUpdates
		// ... (Apply business logic here if needed)
		rdb.Publish("inventory_updates", string(update)).Result()
	}
}

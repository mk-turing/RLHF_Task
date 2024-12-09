package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/go-redis/redis"
	"github.com/gorilla/websocket"
)

const (
	// Maximum message size allowed from client
	maxMessageSize = 512

	// Time allowed to write message to peer before closing the connection
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Time allowed to write message to the redis channel
	redisWriteWait = 60 * time.Second
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type client struct {
	wg          *sync.WaitGroup
	conn        *websocket.Conn
	done        chan struct{}
	send        chan []byte
	redisClient *redis.Client
}

func (c *client) read() {
	defer c.wg.Done()
	defer func() {
		close(c.done) // Signal the `write` method to stop
		c.conn.Close()
	}()

	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error reading: %v", err)
			}
			return
		}

		// Send the message to Redis channel
		if err := c.redisClient.RPush("chat", message).Err(); err != nil {
			log.Printf("error pushing message to Redis: %v", err)
		}

		c.conn.SetReadDeadline(time.Now().Add(pongWait))
	}
}

func (c *client) write() {
	defer c.wg.Done()

	t := time.NewTicker(pingPeriod)
	defer t.Stop()

	for {
		select {
		case message := <-c.send:
			err := c.conn.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				log.Printf("error writing: %v", err)
				return
			}
		case <-t.C:
			err := c.conn.WriteMessage(websocket.PingMessage, []byte{})
			if err != nil {
				log.Printf("error writing ping: %v", err)
				return
			}
		case <-c.done:
			return
		}
	}
}

func newClient(conn *websocket.Conn, redisClient *redis.Client) *client {
	var wg sync.WaitGroup
	c := &client{
		wg:          &wg,
		conn:        conn,
		send:        make(chan []byte),
		redisClient: redisClient,
	}

	wg.Add(2)

	go c.write()
	go c.read()

	return c
}

func serveWS(w http.ResponseWriter, r *http.Request, redisClient *redis.Client) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	defer conn.Close()

	c := newClient(conn, redisClient)

	// Wait for all clients connections to terminate
	c.wg.Wait()
}

func main() {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	_, err := redisClient.Ping().Result()
	if err != nil {
		log.Fatalf("error connecting to Redis: %v", err)
	}

	fmt.Println("Server is running on :8080")
	log.Fatal(http.ListenAndServe(":8080", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/ws" {
			serveWS(w, r, redisClient)
		} else {
			http.ServeFile(w, r, "index.html")
		}
	})))
}

package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/gomodule/redigo/redis"
)

var wg sync.WaitGroup

type QueryParam struct {
	URL    string            `json:"url"`
	Params map[string]string `json:"params"`
}

var qps = []QueryParam{}
var c *redis.Pool

func main() {
	// Set up Redis connection pool
	redisAddr := os.Getenv("REDIS_ADDR")
	c = &redis.Pool{
		MaxIdle:     8,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", redisAddr)
		},
	}

	wg.Add(1)
	go processMonitorEvents()

	http.HandleFunc("/monitor", handleMonitor)
	http.HandleFunc("/status", handleStatus)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func processMonitorEvents() {
	defer wg.Done()
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		queryParameters := queryFromRedis()
		qps = append(qps, queryParameters)
	}
}

func queryFromRedis() QueryParam {
	var q QueryParam
	conn := c.Get()
	defer conn.Close()
	reply, err := conn.Do("RPOP", "query_params")
	if err != nil {
		log.Printf("Error retrieving query parameter from Redis: %v\n", err)
		return q
	}

	err = json.Unmarshal([]byte(reply.(string)), &q)
	if err != nil {
		log.Printf("Error unmarshalling JSON from Redis: %v\n", err)
		return q
	}

	return q
}

func handleMonitor(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		var q QueryParam
		err := json.NewDecoder(r.Body).Decode(&q)
		if err != nil {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		conn := c.Get()
		defer conn.Close()
		b, err := json.Marshal(q)
		if err != nil {
			log.Printf("Error marshalling JSON to Redis: %v\n", err)
			return
		}
		_, err = conn.Do("LPUSH", "query_params", string(b))
		if err != nil {
			log.Printf("Error inserting query parameter into Redis: %v\n", err)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func handleStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(qps)
	}
}

// main.go

package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/gorilla/mux"
)

var (
	redisPool *redis.Pool
)

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/query", handleQuery)

	log.Println("Starting server on port 8000")
	http.ListenAndServe(":8000", router)
}

func handleQuery(w http.ResponseWriter, r *http.Request) {
	queryParam := r.URL.Query().Get("param")
	if queryParam == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Parameter 'param' is required"))
		return
	}

	conn := redisPool.Get()
	defer conn.Close()

	_, err := conn.Do("SET", "query:"+queryParam, 1)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to save query parameter"))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Query parameter saved successfully"))
}

func init() {
	redisPool = &redis.Pool{
		MaxIdle:   3,
		MaxActive: 10,
		Dial:      func() (redis.Conn, error) { return redis.Dial("tcp", "redis:6379") },
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
}

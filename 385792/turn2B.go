package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"
)

const (
	serverAddress   = ":8081"
	healthcheckPath = "/healthz"
)

type Healthcheck struct {
	Status string `json:"status"`
}

func main() {
	// Start the HTTP server for healthchecks
	go startHealthcheckServer()

	// Start the TCP server
	listenAndServe()
}

func startHealthcheckServer() {
	http.HandleFunc(healthcheckPath, func(w http.ResponseWriter, r *http.Request) {
		healthcheck := Healthcheck{Status: "healthy"}
		json.NewEncoder(w).Encode(healthcheck)
		w.WriteHeader(http.StatusOK)
	})
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Error starting healthcheck server: %v", err)
	}
}

func listenAndServe() {
	listener, err := net.Listen("tcp", serverAddress)
	if err != nil {
		log.Fatalf("Error listening: %v", err)
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			if nerr, ok := err.(net.Error); ok && nerr.Temporary() {
				log.Printf("Temporary error accepting connection: %v", err)
				continue
			}
			log.Fatalf("Error accepting connection: %v", err)
		}

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	for {
		data := make([]byte, 1024)
		_, err := conn.Read(data)
		if err != nil {
			if nerr, ok := err.(net.Error); ok && nerr.Temporary() {
				log.Printf("Temporary error reading from connection: %v", err)
				continue
			}
			log.Printf("Error reading from connection: %v", err)
			return
		}

		// Process the received data
		processedData := processData(string(data))

		// Write data back to the client
		if _, err := conn.Write([]byte(processedData)); err != nil {
			if nerr, ok := err.(net.Error); ok && nerr.Temporary() {
				log.Printf("Temporary error writing to connection: %v", err)
				continue
			}
			log.Printf("Error writing to connection: %v", err)
			return
		}
	}
}

func processData(input string) string {
	// Simulate data processing
	time.Sleep(1 * time.Second)
	return fmt.Sprintf("Processed: %s", input)
}

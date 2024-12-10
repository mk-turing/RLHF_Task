package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"sync"
	"time"
)

type Request struct {
	Service string `json:"service"`
}

type Response struct {
	Service string `json:"service"`
	Message string `json:"message"`
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	var req Request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	var wg sync.WaitGroup
	responses := make(chan Response, 1)

	wg.Add(1)
	go func() {
		defer wg.Done()
		// Simulate a call to a downstream service
		response := processRequest(req.Service)
		responses <- response
	}()

	wg.Wait()
	close(responses)

	for resp := range responses {
		json.NewEncoder(w).Encode(resp)
	}
}

func processRequest(service string) Response {
	time.Sleep(time.Duration(rand.Intn(500)) * time.Millisecond) // Simulate processing time
	return Response{Service: service, Message: fmt.Sprintf("Processed by service %s", service)}
}

func main() {
	http.HandleFunc("/api/process", handleRequest)
	fmt.Println("API Gateway starting on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Error starting server:", err)
	}
}

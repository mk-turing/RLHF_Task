package main

import (
	"fmt"
	"net/http"
)

// handleStatus handles the health check requests.
func handleStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		w.WriteHeader(http.StatusOK) // Respond with 200 OK if the server is healthy
	}
}

func main() {
	// Handle the health check route
	http.HandleFunc("/status", handleStatus)

	// Start the server
	fmt.Println("Server running on port 8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}

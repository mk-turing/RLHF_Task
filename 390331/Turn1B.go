package main

import (
	"fmt"
	"log"
	"net/http"
)

// Custom error type
type AppError struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

func (e *AppError) Error() string {
	return fmt.Sprintf("code: %d, message: %s", e.Code, e.Message)
}

func readFile(filename string) ([]byte, error) {
	// Simulate an error
	if filename == "nonexistent.txt" {
		return nil, &AppError{Message: "file not found", Code: 404}
	}
	return []byte("file content"), nil
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	filename := r.URL.Query().Get("file")

	data, err := readFile(filename)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("Error handling request: %v", err)
		return
	}

	w.Write(data)
}

func main() {
	http.HandleFunc("/read", handleRequest)

	log.Print("Server starting on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}

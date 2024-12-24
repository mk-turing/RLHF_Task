package main

import (
	"log"
	"net/http"
)

// ErrorHandlerMiddleware is a middleware function that handles errors by logging them
// and sending a 500 internal server error response.
func ErrorHandlerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("Panic: %v", err)
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte("Internal Server Error"))
			}
		}()

		// Directly use the next handler without assigning return value
		next.ServeHTTP(w, r)
	})
}

// SimpleHandler is a mock handler function that might cause an error.
func SimpleHandler(w http.ResponseWriter, r *http.Request) error {
	if r.URL.Path != "/simple" {
		// Return a 404 status directly
		http.NotFound(w, r)
		return nil
	}

	// Simulate an error
	_, err := w.Write([]byte("This might fail sometimes"))
	if err != nil {
		return err
	}

	return nil
}

func main() {
	// Wrap SimpleHandler with http.HandlerFunc to make it a valid http.Handler
	http.Handle("/", ErrorHandlerMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := SimpleHandler(w, r); err != nil {
			log.Printf("Error: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	})))

	log.Println("Server is running on http://localhost:8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

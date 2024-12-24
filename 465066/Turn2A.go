package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

// Response represents the API response.
type Response struct {
	Data  interface{} `json:"data,omitempty"`
	Error *Error      `json:"error,omitempty"`
}

// Error represents an API error.
type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// Implement the error interface for the Error struct
func (e *Error) Error() string {
	return e.Message
}

// HandlerFunc is the type of a handler function.
type HandlerFunc func(w http.ResponseWriter, r *http.Request) (interface{}, error)

// ErrorHandler wraps a HandlerFunc and handles errors, including logging.
func ErrorHandler(next HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		data, err := next(w, r)
		duration := time.Since(startTime)

		if err != nil {
			handleError(w, err, duration)
			return
		}
		handleResponse(w, data, duration)
	}
}

func handleError(w http.ResponseWriter, err error, duration time.Duration) {
	code := http.StatusInternalServerError
	message := err.Error()

	switch e := err.(type) {
	case *Error:
		code = e.Code
		message = e.Message
	}

	log.Printf("Error: %+v (%dms)\n", err, duration.Milliseconds())
	response := Response{Error: &Error{Code: code, Message: message}}
	handleResponse(w, response, duration)
}

func handleResponse(w http.ResponseWriter, data interface{}, duration time.Duration) {
	// Set response headers and encode the response.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Use an encoder to handle JSON encoding, e.g., json.NewEncoder(w).Encode(data)
	resp := Response{Data: data}
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("Failed to encode response: %v (%dms)\n", err, duration.Milliseconds())
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"error": {"code": 500, "message": "Internal Server Error"}}`)
		return
	}

	log.Printf("Response sent successfully: %v (%dms)\n", resp, duration.Milliseconds())
}

func main() {
	http.HandleFunc("/example", ErrorHandler(exampleHandler))
	http.ListenAndServe(":8080", nil)
}

func exampleHandler(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	// Simulate some business logic that may return errors.
	// In a real application, you would have your own error types and handling logic.
	data, err := someBusinessLogic()
	if err != nil {
		return nil, err // Return the error to the ErrorHandler.
	}
	return data, nil // Return the data successfully.
}

func someBusinessLogic() (interface{}, error) {
	// Simulate a business logic function that can return an error.
	return nil, fmt.Errorf("some error occurred")
}

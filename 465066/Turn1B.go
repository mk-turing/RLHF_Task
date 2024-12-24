package main

import (
	"fmt"
	"net/http"
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

// ErrorHandler wraps a HandlerFunc and handles errors.
func ErrorHandler(next HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := next(w, r)
		if err != nil {
			handleError(w, err)
			return
		}
		handleResponse(w, data)
	}
}

func handleError(w http.ResponseWriter, err error) {
	// Handle different types of errors or use custom error structs with predefined codes.
	code := http.StatusInternalServerError
	message := err.Error()

	// You can customize the error handling based on your specific use case.
	switch e := err.(type) {
	case *Error:
		code = e.Code
		message = e.Message
	}

	response := Response{Error: &Error{Code: code, Message: message}}
	handleResponse(w, response)
}

func handleResponse(w http.ResponseWriter, data interface{}) {
	// Set response headers and encode the response.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	// Use an encoder to handle JSON encoding, e.g., json.NewEncoder(w).Encode(data)
	fmt.Fprintf(w, `{"data": %+v}`, data)
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

package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

// Custom Error Types
type AppError struct {
	Code    int
	Message string
	Err     error
}

func (e *AppError) Error() string {
	return fmt.Sprintf("Code: %d, Message: %s, Error: %v", e.Code, e.Message, e.Err)
}

func NewAppError(code int, message string, err error) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

// Error Logger
var (
	errorLog = log.New(os.Stderr, "ERROR: ", log.LstdFlags)
)

func LogError(err error) {
	errorLog.Println(err)
}

// Service Layer
type UserService struct{}

func (s *UserService) GetUser(id string) (string, error) {
	// Simulating an error
	if id == "" {
		return "", NewAppError(400, "user ID cannot be empty", fmt.Errorf("invalid input"))
	}
	return "UserData", nil
}

// Handler with Centralized Error Handling
type Handler struct {
	userService UserService
}

func (h *Handler) GetUserHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	userData, err := h.userService.GetUser(id)
	if err != nil {
		// Centralized error handling
		appError, ok := err.(*AppError)
		if ok {
			LogError(err) // log the error
			http.Error(w, appError.Message, appError.Code)
			return
		}
		// Fallback for unexpected errors
		LogError(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Write([]byte(userData))
}

func main() {
	handler := &Handler{
		userService: UserService{},
	}

	http.HandleFunc("/user", handler.GetUserHandler)

	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

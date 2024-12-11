package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

// Custom Error Types
type AppError interface {
	Error() string
	GetCode() int
	GetType() string
}

type appError struct {
	Code    int
	Message string
	Err     error
	Type    string
}

func (e *appError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("Code: %d, Type: %s, Message: %s, Error: %v", e.Code, e.Type, e.Message, e.Err)
	}
	return fmt.Sprintf("Code: %d, Type: %s, Message: %s", e.Code, e.Type, e.Message)
}

func (e *appError) GetCode() int {
	return e.Code
}

func (e *appError) GetType() string {
	return e.Type
}

func NewApplicationError(code int, message string, err error) AppError {
	return &appError{
		Code:    code,
		Message: message,
		Err:     err,
		Type:    "ApplicationError",
	}
}

func NewNetworkError(code int, message string, err error) AppError {
	return &appError{
		Code:    code,
		Message: message,
		Err:     err,
		Type:    "NetworkError",
	}
}

func NewSystemError(code int, message string, err error) AppError {
	return &appError{
		Code:    code,
		Message: message,
		Err:     err,
		Type:    "SystemError",
	}
}

// Enhanced Error Logger
var (
	errorLog = log.New(os.Stderr, "ERROR: ", log.LstdFlags)
)

func LogError(err AppError) {
	errorLog.Printf("[%s] %s", err.GetType(), err.Error())
}

// Service Layer
type UserService struct{}

func (s *UserService) GetUser(id string) (string, AppError) {
	if id == "" {
		return "", NewApplicationError(400, "user ID cannot be empty", nil)
	}

	// Simulate a network error
	if id == "network" {
		return "", NewNetworkError(502, "failed to reach user service", nil)
	}

	// Simulate a system error
	if id == "system" {
		return "", NewSystemError(500, "database connection failed", nil)
	}

	return "UserData", nil
}

// Handler Layer
type Handler struct {
	userService UserService
}

func (h *Handler) GetUserHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	userData, err := h.userService.GetUser(id)
	if err != nil {
		// Centralized error handling
		LogError(err) // log the error

		http.Error(w, err.Error(), err.GetCode())
		return
	}
	w.Write([]byte(userData))
}

// Main Application Entry Point
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

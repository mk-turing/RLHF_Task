package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog"
	"net/http"
	"os"
	"strconv"
)

// Example custom error type
type DataNotFoundError struct {
	msg string
}

func (e *DataNotFoundError) Error() string {
	return e.msg
}

// Simulates a service that fetches data
func fetchData(ctx context.Context, id int) (*struct {
	Name string `json:"name"`
}, error) {
	// Simple lookup simulation
	data := map[int]struct {
		Name string `json:"name"`
	}{
		1: {Name: "Alice"},
		2: {Name: "Bob"},
	}

	if entry, ok := data[id]; ok {
		return &entry, nil
	}

	// Return custom error for data not found
	return nil, &DataNotFoundError{fmt.Sprintf("data not found for id: %d", id)}
}

// HTTP handler function
func fetchDataHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	data, err := fetchData(ctx, id)
	if err != nil {
		// Log the error with context using zerolog
		log := zerolog.New(os.Stdout).With().Timestamp().Logger()
		log.Error().Err(err).Msg("error fetching data")

		// Propagate the error to the client
		writeError(w, http.StatusNotFound, err)
		return
	}

	// Write the response
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(data); err != nil {
		// Log the error writing response using zerolog
		log := zerolog.New(os.Stdout).With().Timestamp().Logger()
		log.Error().Err(err).Msg("error writing response")
		writeError(w, http.StatusInternalServerError, err)
		return
	}
}

// Utility function to write error responses
func writeError(w http.ResponseWriter, statusCode int, err error) {
	w.WriteHeader(statusCode)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"error": err.Error(),
	})
}

func main() {
	// Set up a logger with context using zerolog
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()

	// Register HTTP handler
	r := mux.NewRouter()
	r.HandleFunc("/data/{id}", fetchDataHandler).Methods("GET")

	// Start HTTP server with zerolog logger
	log := logger
	log.Fatal().Err(http.ListenAndServe(":8080", r)).Msg("server failed")
}

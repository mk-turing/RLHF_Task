package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Logger struct{}

func (l *Logger) Log(msg string) {
	fmt.Println("Log:", msg)
}

type Middleware[T any] struct {
	Logger *Logger
}

func (m *Middleware[T]) Handle(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		var param T
		paramBytes, err := json.Marshal(q)
		if err != nil {
			http.Error(w, "Error marshaling query parameters to JSON", http.StatusBadRequest)
			return
		}
		err = json.Unmarshal(paramBytes, &param)
		if err != nil {
			http.Error(w, "Error unmarshaling JSON to struct", http.StatusBadRequest)
			return
		}
		m.Logger.Log(fmt.Sprintf("Query Param: %v", param))
		h.ServeHTTP(w, r)
	})
}

type NestedParams struct {
	User     User     `json:"user"`
	Location Location `json:"location"`
}

type User struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

type Location struct {
	City    string `json:"city"`
	Country string `json:"country"`
}

func main() {
	l := &Logger{}
	mNested := &Middleware[NestedParams]{Logger: l}

	http.Handle("/nested", mNested.Handle(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		params := r.URL.Query().Get("params")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if params != "" {
			json.NewEncoder(w).Encode(params)
		}
	})))

	fmt.Println("Server is running on port 8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}

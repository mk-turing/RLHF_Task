package main

import (
	"fmt"
	"net/http"
	"reflect"
)

type Logger struct{}

func (l *Logger) Log(msg string) {
	fmt.Println("Log:", msg)
}

type Middleware struct {
	Logger *Logger
}

func (m *Middleware) Handle(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		params := r.URL.Query()
		m.Logger.Log("Query Parameters:")

		for key, values := range params {
			// Extract the first value for logging
			value := values[0]
			m.Logger.Log(fmt.Sprintf("%s: %v (%T)", key, value, reflect.TypeOf(value)))
		}

		h.ServeHTTP(w, r)
	})
}

func main() {
	l := &Logger{}
	m := &Middleware{Logger: l}

	http.Handle("/", m.Handle(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Request processed successfully.")
	})))

	fmt.Println("Server is running on port 8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}

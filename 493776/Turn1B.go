package main

import (
	"fmt"
	"net/http"
	"strconv"
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
		// extract query param from request
		if val, ok := q["param"]; ok {
			// Handle string type
			if _, ok := any(param).(string); ok {
				param = any(val[0]).(T) // type assertion to T
			} else if _, ok := any(param).(int); ok {
				i, err := strconv.Atoi(val[0])
				if err != nil {
					http.Error(w, "Invalid query parameter type", http.StatusBadRequest)
					return
				}
				param = any(i).(T) // type assertion to T
			}
			// add more types if needed
		}
		m.Logger.Log(fmt.Sprintf("Query Param: %v", param))
		h.ServeHTTP(w, r)
	})
}

func main() {
	l := &Logger{}
	mStr := &Middleware[string]{Logger: l}
	mInt := &Middleware[int]{Logger: l}

	http.Handle("/string", mStr.Handle(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, %s!", r.URL.Query().Get("param"))
	})))

	http.Handle("/int", mInt.Handle(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		param := r.URL.Query().Get("param")
		i, _ := strconv.Atoi(param)
		fmt.Fprintf(w, "The value of param is %d", i)
	})))

	fmt.Println("Server is running on port 8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}

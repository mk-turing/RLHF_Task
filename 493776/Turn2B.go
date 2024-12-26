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
		if err := unmarshalQuery(q, &param); err != nil {
			http.Error(w, "Invalid query parameter type", http.StatusBadRequest)
			return
		}
		m.Logger.Log(fmt.Sprintf("Query Param: %v", param))
		h.ServeHTTP(w, r)
	})
}

func unmarshalQuery(q map[string][]string, v interface{}) error {
	b, err := json.Marshal(q)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, v)
}

type User struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

type Address struct {
	City    string `json:"city"`
	Country string `json:"country"`
}

type UserRequest struct {
	User    User     `json:"user"`
	Address *Address `json:"address"`
}

func main() {
	l := &Logger{}
	mUser := &Middleware[UserRequest]{Logger: l}

	http.Handle("/user", mUser.Handle(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var user UserRequest
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		fmt.Fprintf(w, "Hello, %s! You are %d years old. You live in %s, %s.",
			user.User.Name, user.User.Age, user.Address.City, user.Address.Country)
	})))

	fmt.Println("Server is running on port 8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}

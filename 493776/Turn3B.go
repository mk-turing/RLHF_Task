package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
)

type Logger struct {
	Format map[reflect.Kind]string
}

func (l *Logger) Log(msg string) {
	fmt.Println("Log:", msg)
}

func (l *Logger) LogParams(params interface{}) {
	v := reflect.ValueOf(params)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	t := v.Type()
	msg := ""
	for i := 0; i < v.NumField(); i++ {
		f := v.Field(i)
		name := t.Field(i).Name
		kind := f.Kind()
		format, ok := l.Format[kind]
		if !ok {
			format = "%v"
		}
		msg += fmt.Sprintf("%s=%s, ", name, fmt.Sprintf(format, f.Interface()))
	}
	if msg != "" {
		l.Log(msg[:len(msg)-2])
	}
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
		m.Logger.LogParams(param)
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
	Name string
	Age  int
}

type Address struct {
	City    string
	Country string
}

type UserRequest struct {
	User    User     `json:"user"`
	Address *Address `json:"address"`
}

func main() {
	l := &Logger{
		Format: map[reflect.Kind]string{
			reflect.String: "%q",
			reflect.Int:    "%d",
			reflect.Ptr:    "%+v",
		},
	}
	mUser := &Middleware[UserRequest]{Logger: l}

	http.Handle("/user", mUser.Handle(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello!")
	})))

	fmt.Println("Server is running on port 8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}

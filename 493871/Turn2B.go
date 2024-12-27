package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/gorilla/mux"
)

func sanitizeQueryParameters(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var sanitizedParams string
		// Check for different HTTP methods
		switch r.Method {
		case "GET", "DELETE":
			// Parse the query string for GET and DELETE methods
			q := r.URL.Query()
			sanitizedParams = sanitizeQueryString(q)
		case "POST", "PUT", "PATCH":
			// For POST, PUT, and PATCH methods, check the form values
			r.ParseForm()
			sanitizedParams = sanitizeFormValues(r.Form)
		default:
			log.Printf("Unsupported HTTP method: %s", r.Method)
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		// Log potential injection attempts
		if strings.Contains(sanitizedParams, "--") {
			log.Printf("Possible SQL injection attempt in method: %s, Path: %s, Sanitized Params: %s", r.Method, r.URL.Path, sanitizedParams)
		}

		// Rebuild the request accordingly based on the method
		if r.Method == "GET" {
			r.URL.RawQuery = sanitizedParams
		} else {
			r.PostForm = make(url.Values)
			for key, value := range r.Form {
				r.PostForm[key] = []string{strings.Join(value, ",")}
			}
		}

		// Call the next handler
		next.ServeHTTP(w, r)
	})
}

func sanitizeQueryString(q url.Values) string {
	sanitizedValues := []string{}
	for key, values := range q {
		for _, value := range values {
			sanitizedValue := sanitize(value)
			sanitizedValues = append(sanitizedValues, fmt.Sprintf("%s=%s", key, sanitizedValue))
		}
	}
	return strings.Join(sanitizedValues, "&")
}

func sanitizeFormValues(form url.Values) string {
	sanitizedValues := []string{}
	for key, values := range form {
		for _, value := range values {
			sanitizedValue := sanitize(value)
			sanitizedValues = append(sanitizedValues, fmt.Sprintf("%s=%s", key, sanitizedValue))
		}
	}
	return strings.Join(sanitizedValues, "&")
}

func sanitize(input string) string {
	sanitized := strings.Map(func(r rune) rune {
		if strings.ContainsRune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_.- ", r) {
			return r
		}
		return -1
	}, input)

	return string(sanitized)
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/example", exampleHandler).Methods("GET", "POST", "PUT")

	r.Use(sanitizeQueryParameters)

	log.Fatal(http.ListenAndServe(":8080", r))
}

func exampleHandler(w http.ResponseWriter, r *http.Request) {
	var param string
	// Access sanitized parameters based on the method
	switch r.Method {
	case "GET":
		q := r.URL.Query()
		param = q.Get("param")
	case "POST", "PUT":
		r.ParseForm()
		param = r.FormValue("param")
	}

	// Respond with the sanitized parameter
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Sanitized Parameter: %s", param)
}

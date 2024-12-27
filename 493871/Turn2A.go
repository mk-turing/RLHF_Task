package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

// Enhanced sanitization method
func sanitize(input string) string {
	// Sanity check: filter out harmful characters
	return strings.Map(func(r rune) rune {
		if strings.ContainsRune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_-@.%=:?", r) {
			return r
		}
		return -1
	}, input)
}

func sanitizeQueryParameters(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		loggingHandler(w, r, func() {
			// Sanitate query parameters for GET requests
			q := r.URL.Query()
			for key, values := range q {
				sanitizedValues := []string{}
				for _, value := range values {
					sanitizedValue := sanitize(value)
					sanitizedValues = append(sanitizedValues, sanitizedValue)
				}
				q.Set(key, strings.Join(sanitizedValues, ","))
			}
			r.URL.RawQuery = q.Encode()

			// Sanitate body parameters for POST, PUT, and PATCH requests
			if r.Method == http.MethodPost || r.Method == http.MethodPut || r.Method == http.MethodPatch {
				decoder := json.NewDecoder(r.Body)
				defer r.Body.Close()
				var m map[string]interface{}
				if err := decoder.Decode(&m); err != nil {
					http.Error(w, err.Error(), http.StatusBadRequest)
					return
				}

				for key, value := range m {
					v := sanitize(value.(string))
					m[key] = v
				}

				b, err := json.Marshal(m)
				if err != nil {
					http.Error(w, err.Error(), http.StatusBadRequest)
					return
				}

				r.Body = http.NoBody
				r.Header.Set("Content-Type", "application/json; charset=UTF-8")
				r.Header.Set("Content-Length", strconv.Itoa(len(b)))
				w.WriteHeader(http.StatusUnprocessableEntity)
				w.Write(b)
			}
		})

		next.ServeHTTP(w, r)
	})
}

func loggingHandler(w http.ResponseWriter, r *http.Request, h func()) {
	// Pretend to log and identify potential injection attempts
	h()

	// Example of logging sanitized requests (possible injection attempt identification)
	for key, value := range r.URL.Query() {
		for _, v := range value {
			if strings.Contains(v, "SQL;") || strings.Contains(v, "XSS:") {
				// You could log here or take additional actions like rate limiting
				log.Printf("Potential injection detected in query param %s: %s\n", key, v)
			}
		}
	}

	if r.Method == http.MethodPost || r.Method == http.MethodPut || r.Method == http.MethodPatch {
		decoder := json.NewDecoder(r.Body)
		defer r.Body.Close()
		var m map[string]interface{}
		if err := decoder.Decode(&m); err != nil {
			log.Printf("Error decoding request body: %s\n", err.Error())
			return
		}

		for key, value := range m {
			v := value.(string)
			if strings.Contains(v, "SQL;") || strings.Contains(v, "XSS:") {
				// You could log here or take additional actions like rate limiting
				log.Printf("Potential injection detected in body param %s: %s\n", key, v)
			}
		}
	}
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/example", exampleHandler).Methods("GET", "POST", "PUT", "DELETE")

	r.Use(sanitizeQueryParameters)

	http.ListenAndServe(":8080", r)
}

func exampleHandler(w http.ResponseWriter, r *http.Request) {
	// Access sanitized query parameters and body (for POST, PUT, PATCH)
	q := r.URL.Query()
	param := q.Get("param")

	if r.Method == http.MethodPost || r.Method == http.MethodPut || r.Method == http.MethodPatch {
		decoder := json.NewDecoder(r.Body)
		defer r.Body.Close()
		var body struct {
			BodyParam string `json:"bodyParam"`
		}
		if err := decoder.Decode(&body); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		// Respond with sanitized body parameter
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Sanitized Body Parameter: " + body.BodyParam))
	}

	// Respond with sanitized query parameter
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Sanitized Parameter: " + param))
}

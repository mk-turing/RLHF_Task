package main

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

type Config struct {
	AllowedOrigins map[string]struct{}
	AllowedMethods map[string]struct{}
	AllowedHeaders map[string]struct{}
}

var (
	configs = map[string]Config{
		"api-key-1": Config{
			AllowedOrigins: map[string]struct{}{
				"https://example1.com": {},
				"https://sub.example1.com": {},
			},
			AllowedMethods: map[string]struct{}{
				"GET":  {},
				"POST": {},
			},
			AllowedHeaders: map[string]struct{}{
				"Content-Type":   {},
				"Authorization":  {},
				"X-Custom-Header": {},
			},
		},
		"api-key-2": Config{
			AllowedOrigins: map[string]struct{}{
				"https://example2.com": {},
			},
			AllowedMethods: map[string]struct{}{
				"GET":  {},
				"PUT":  {},
				"DELETE": {},
			},
			AllowedHeaders: map[string]struct{}{
				"Content-Type": {},
			},
		},
	}
	// Max number of workers to handle requests concurrently
	numWorkers = 100
)

func queryParameterInspector(w http.ResponseWriter, r *http.Request) {
	// Dispatch the request to a worker goroutine for asynchronous processing
	workerCh <- request{w, r}
}

type request struct {
	w http.ResponseWriter
	r *http.Request
}

var workerCh = make(chan request, numWorkers)

func worker() {
	for req := range workerCh {
		processRequest(req.w, req.r)
	}
}

func processRequest(w http.ResponseWriter, r *http.Request) {
	// Parse the URL and handle errors
	parsedURL, err := url.Parse(r.URL.String())
	if err != nil {
		handleError(w, err, http.StatusBadRequest, "Invalid URL format. Please check the URL syntax.")
		return
	}

	// Extract API key from header
	apiKey := r.Header.Get("X-Api-Key")
	config, ok := configs[apiKey]
	if !ok {
		handleError(w, errors.New("Invalid API key"), http.StatusUnauthorized, "Invalid API key. Please provide a valid API key.")
		return
	}

	// Perform CORS checks using set operations for efficient lookups
	origin := r.Header.Get("Origin")
	if _, allowed := config.AllowedOrigins[origin]; !allowed {
		handleError(w, errors.New("Unsupported CORS policy"), http.StatusForbidden, fmt.Sprintf("Origin %s not allowed. Please check your CORS configuration.", origin))
		return
	}

	method := r.Method
	if _, allowed := config.AllowedMethods[method]; !allowed {
		handleError(w, errors.New("Method not allowed"), http.StatusMethodNotAllowed, fmt.Sprintf("Method %s is not allowed for this endpoint.", method))
		return
	}

	// Preflight request handling
	if method == "OPTIONS" {
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Methods", strings.Join(config.AllowedMethods, ", "))
		w.Header().Set("Access-Control-Allow-Headers", strings.Join(config.AllowedHeaders, ", "))
		w.WriteHeader(http.StatusOK)
		return
	}

	// Parse query parameters from both URL and form
	query := parsedURL.Query()
	query.Add("form-param", r.FormValue("form-param"))

	// Validate query parameters
	requiredFields := map[string]bool{"param1": true, "form-param": true}
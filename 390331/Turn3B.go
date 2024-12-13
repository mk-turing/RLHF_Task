package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"
)

// Custom error type for traceable errors
type TraceError struct {
	error
	RequestID string `json:"request_id"`
	Service   string `json:"service"`
}

func (e *TraceError) Error() string {
	return fmt.Sprintf("[%s:%s] %s", e.Service, e.RequestID, e.error)
}

// HTTP Client with retry and backoff
type retryingClient struct {
	client *http.Client
}

func newRetryingClient() *retryingClient {
	return &retryingClient{client: &http.Client{Timeout: 5 * time.Second}}
}

func (rc *retryingClient) Do(ctx context.Context, req *http.Request) (*http.Response, error) {
	const (
		maxAttempts = 3
		backoff     = time.Second * 2
	)

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		resp, err := rc.client.Do(req.WithContext(ctx))
		if err != nil {
			log.Printf("Attempt %d failed: %v, retrying...", attempt, err)
			if attempt < maxAttempts {
				time.Sleep(backoff * time.Duration(attempt))
				continue
			}
		} else {
			defer resp.Body.Close()
			return resp, nil
		}
	}
	return nil, nil
}

// Handler for the remote service
func remoteServiceHandler(w http.ResponseWriter, r *http.Request) {
	time.Sleep(100 * time.Millisecond)
	if r.URL.Query().Get("fail") == "true" {
		http.Error(w, "Remote service failure", http.StatusInternalServerError)
		return
	}
	w.Write([]byte(`{"message": "Remote service successful"}`))
}

// Local handler that calls the remote service
func localHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	reqID := r.Header.Get("X-Request-ID")
	if reqID == "" {
		reqID = fmt.Sprintf("%08x", uint64(time.Now().UnixNano()))
	}

	remoteURL, _ := url.Parse("http://localhost:8081/remote")
	q := remoteURL.Query()
	q.Set("delay", r.URL.Query().Get("delay"))
	q.Set("fail", r.URL.Query().Get("fail"))
	remoteURL.RawQuery = q.Encode()

	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, remoteURL.String(), nil)
	req.Header.Set("X-Request-ID", reqID)

	client := newRetryingClient()
	resp, err := client.Do(ctx, req)
	if err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		errorResp, _ := json.Marshal(&TraceError{err, reqID, "local"})
		w.Write(errorResp)
		return
	}
	defer resp.Body.Close()

	var respData map[string]string
	if err := json.NewDecoder(resp.Body).Decode(&respData); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		errorResp, _ := json.Marshal(&TraceError{err, reqID, "local"})
		w.Write(errorResp)
		return
	}

	w.Write([]byte(fmt.Sprintf("Local handler got: %#v", respData)))
}

func main() {
	http.HandleFunc("/remote", remoteServiceHandler)
	http.HandleFunc("/local", localHandler)

	log.Print("Starting services...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}

	if err := http.ListenAndServe(":8081", nil); err != nil {
		log.Fatal(err)
	}
}

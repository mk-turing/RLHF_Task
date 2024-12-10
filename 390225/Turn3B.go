package main

import (
	"log"
	"net/http"
)

func recoveryMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("Panic recovered: %v\n", err)
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
		}()
		next(w, r)
	}
}

func main() {
	http.HandleFunc("/gateway", recoveryMiddleware(func(w http.ResponseWriter, r *http.Request) {
		panic("Simulated Panic") // This will be caught by the recovery middleware
	}))

	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

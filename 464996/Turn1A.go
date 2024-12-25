package main

import (
	"net/http"
	"time"
)

var counter int

func handleDefer(w http.ResponseWriter, r *http.Request) {
	counter++
	defer func() {
		counter--
	}()
	time.Sleep(1 * time.Millisecond)
}

func handleWithoutDefer(w http.ResponseWriter, r *http.Request) {
	counter++
	defer func() {
		counter--
	}()
	time.Sleep(1 * time.Millisecond)
}

func main() {
	http.HandleFunc("/defer", handleDefer)
	http.HandleFunc("/without-defer", handleWithoutDefer)

	go func() {
		for {
			time.Sleep(10 * time.Second)
			println("Current counter:", counter)
		}
	}()

	http.ListenAndServe(":8080", nil)
}

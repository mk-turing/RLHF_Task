package main

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"
)

const (
	variantCookieName = "abctest"
	defaultVariant    = "a"
	redisAddr         = "localhost:6379" // Redis server address
	cachePrefix       = "abtest:variant:"
	logFileName       = "turn4B_variant_log.txt"
)

var (
	variants = []string{"a", "b", "c"}
	// In-memory session store for demonstration purposes
	sessions    = make(map[string]string, 100)
	sessionLock = sync.Mutex{}
	rdb         *redis.Client
)

func init() {
	// Initialize Redis client
	rdb = redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: "", // No password set
		DB:       0,  // Use default DB
	})

	rand.Seed(time.Now().UnixNano())
	// Clear existing log file on startup
	clearLogFile()
}

func main() {
	rand.Seed(time.Now().UnixNano())
	http.HandleFunc("/", handleRequest)
	http.HandleFunc("/set-variant", setVariant)
	http.HandleFunc("/get-variant", getVariant)
	http.HandleFunc("/logs", logHandler)

	// Clear existing log file on startup
	clearLogFile()

	fmt.Println("Server starting on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	variant := getVariantFromCookieOrSession(r)
	serveVariantContent(w, variant)
	logVariantInteraction(r.RemoteAddr, variant)
}

func setVariant(w http.ResponseWriter, r *http.Request) {
	queryParams, _ := url.ParseQuery(r.URL.Query().Encode())
	variant := queryParams.Get("variant")
	if variant != "" && containsVariant(variant) {
		setVariantCookie(w, variant)
		setVariantSession(r, variant)
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Variant set to %s\n", variant)
	} else {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Invalid variant parameter\n")
	}
}

func getVariant(w http.ResponseWriter, r *http.Request) {
	variant := getVariantFromCookieOrSession(r)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Variant: %s\n", variant)
}

func logHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, logFileName)
}

func getVariantFromCache(sessionID string) string {
	cacheKey := cachePrefix + sessionID
	variant, err := rdb.Get(context.Background(), cacheKey).Result()
	if err == nil {
		return variant
	}
	return ""
}

func setVariantCache(sessionID, variant string) {
	cacheKey := cachePrefix + sessionID
	rdb.Set(context.Background(), cacheKey, variant, 7*24*time.Hour) // Cache variant for 7 days
}

func getVariantFromCookieOrSession(r *http.Request) string {
	if variant := getVariantFromCookie(r); variant != "" {
		return variant
	}
	sessionID := r.RemoteAddr
	if variant := getVariantFromCache(sessionID); variant != "" {
		return variant
	}
	return chooseRandomVariant()
}

func getVariantFromCookie(r *http.Request) string {
	cookie, err := r.Cookie(variantCookieName)
	if err == nil {
		return cookie.Value
	}
	return ""
}

func setVariantCookie(w http.ResponseWriter, variant string) {
	cookie := http.Cookie{
		Name:     variantCookieName,
		Value:    variant,
		Path:     "/",
		Expires:  time.Now().AddDate(0, 0, 7),
		HttpOnly: true,
	}
	http.SetCookie(w, &cookie)
}

func setVariantSession(r *http.Request, variant string) {
	sessionLock.Lock()
	defer sessionLock.Unlock()
	sessions[r.RemoteAddr] = variant
}

func chooseRandomVariant() string {
	return variants[rand.Intn(len(variants))]
}

func serveVariantContent(w http.ResponseWriter, variant string) {
	switch variant {
	case "a":
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "<html><body><h1>Variant A</h1><p>This is the content for variant A.</p></body></html>")
	case "b":
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "<html><body><h1>Variant B</h1><p>This is the content for variant B.</p></body></html>")
	case "c":
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "<html><body><h1>Variant C</h1><p>This is the content for variant C.</p></body></html>")
	default:
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "<html><body><h1>404 Not Found</h1><p>Invalid variant parameter.</p></body></html>")
	}
}

func containsVariant(variant string) bool {
	return strings.Contains(strings.Join(variants, ","), variant)
}

func logVariantInteraction(ip, variant string) {
	logEntry := fmt.Sprintf("%s,%s,%d\n", time.Now().Format("2006-01-02 15:04:05"), ip, variant)
	logFile, err := os.OpenFile(logFileName, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		fmt.Printf("Error logging variant: %v\n", err)
		return
	}
	defer logFile.Close()
	if _, err := logFile.WriteString(logEntry); err != nil {
		fmt.Printf("Error writing to log file: %v\n", err)
	}
}

func clearLogFile() {
	if _, err := os.Stat(logFileName); err != nil {
		if os.IsNotExist(err) {
			return // Log file does not exist, return
		} else {
			fmt.Printf("Error checking log file: %v\n", err)
			return
		}
	}
	if err := os.Remove(logFileName); err != nil {
		fmt.Printf("Error removing log file: %v\n", err)
	}
}

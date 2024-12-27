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

	"github.com/go-redis/redis/v8"
)

const (
	variantCookieName = "abctest"
	defaultVariant    = "a"
	logFileName       = "turn4A_variant_log.txt"
	cacheKeyPrefix    = "variant_"
	cacheExpiration   = 5 * time.Minute
)

var (
	variants    = []string{"a", "b", "c"}
	sessions    = make(map[string]string, 100)
	sessionLock = sync.Mutex{}
	redisClient *redis.Client
)

func init() {
	ctx := context.Background()
	redisClient = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		DB:       0,
		Password: "",
	})
	_, err := redisClient.Ping(ctx).Result()
	if err != nil {
		fmt.Printf("Error connecting to Redis: %v\n", err)
		os.Exit(1)
	}
	rand.Seed(time.Now().UnixNano())
	clearLogFile()
}

func main() {
	http.HandleFunc("/", handleRequest)
	http.HandleFunc("/set-variant", setVariant)
	http.HandleFunc("/get-variant", getVariant)
	http.HandleFunc("/logs", logHandler)

	fmt.Println("Server starting on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	variant := getVariantFromCacheOrSession(r)
	serveVariantContent(w, variant)
	logVariantInteraction(r.RemoteAddr, variant)
}

func setVariant(w http.ResponseWriter, r *http.Request) {
	queryParams, _ := url.ParseQuery(r.URL.Query().Encode())
	variant := queryParams.Get("variant")
	if variant != "" && containsVariant(variant) {
		setVariantCookie(w, variant)
		setVariantSession(r, variant)
		invalidateVariantCache(r)
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Variant set to %s\n", variant)
	} else {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Invalid variant parameter\n")
	}
}

func getVariant(w http.ResponseWriter, r *http.Request) {
	variant := getVariantFromCacheOrSession(r)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Variant: %s\n", variant)
}

func logHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, logFileName)
}

func getVariantFromCacheOrSession(r *http.Request) string {
	ctx := context.Background()
	sessionID := r.RemoteAddr
	variant, err := redisClient.Get(ctx, cacheKeyPrefix+sessionID).Result()
	if err == redis.Nil {
		if variant, ok := sessions[sessionID]; ok {
			return variant
		}
		return chooseRandomVariant()
	} else if err != nil {
		fmt.Printf("Error retrieving variant from cache: %v\n", err)
		return chooseRandomVariant()
	}
	return variant
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

func invalidateVariantCache(r *http.Request) {
	ctx := context.Background()
	sessionID := r.RemoteAddr
	err := redisClient.Del(ctx, cacheKeyPrefix+sessionID).Err()
	if err != nil {
		fmt.Printf("Error invalidating variant cache: %v\n", err)
	}
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

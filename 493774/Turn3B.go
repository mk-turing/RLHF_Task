package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/gorilla/sessions"
	"github.com/sirupsen/logrus"
)

const (
	sessionKey        = "session-key"
	variantCookieName = "abctest"
	defaultVariant    = "a"
)

var (
	variants = []string{"a", "b", "c"}
	store    = sessions.NewCookieStore([]byte("secret-key"))
	logger   = logrus.New()
)

func init() {
	rand.Seed(time.Now().UnixNano())
	logger.SetLevel(logrus.InfoLevel)
}

func main() {
	http.HandleFunc("/", handleRequest)
	fmt.Println("Server starting on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, sessionKey)
	if err != nil {
		logger.Errorf("Error getting session: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	variant := getVariantFromSession(session)
	if variant == "" {
		variant = chooseRandomVariant()
		setVariantSession(session, variant)
		setVariantCookie(w, variant)
	}

	logUserVariant(session, variant)
	serveVariantContent(w, variant)
}

func getVariantFromSession(session *sessions.Session) string {
	if variant, ok := session.Values["variant"].(string); ok {
		return variant
	}
	return ""
}

func setVariantSession(session *sessions.Session, variant string) {
	session.Values["variant"] = variant
	if err := session.Save(r, w); err != nil {
		logger.Errorf("Error saving session: %v", err)
	}
}

func logUserVariant(session *sessions.Session, variant string) {
	logger.Infof("User %v received variant: %v", session.ID, variant)
}

func chooseRandomVariant() string {
	// Select a random variant from the available variants
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

func setVariantCookie(w http.ResponseWriter, variant string) {
	cookie := http.Cookie{
		Name:     variantCookieName,
		Value:    variant,
		Path:     "/",
		Expires:  time.Now().AddDate(0, 0, 7), // Set cookie expiration for 7 days
		HttpOnly: true,
	}
	http.SetCookie(w, &cookie)
}

package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"sync"
)

type Document struct {
	Content string `json:"content"`
	Version int    `json:"version"`
}

type DocumentChange struct {
	Type    string `json:"type"`
	Index   int    `json:"index"`
	Content string `json:"content"`
	Version int    `json:"version"`
}

var (
	mu          sync.RWMutex
	doc         *Document                = &Document{Content: "", Version: 0}
	connections map[*websocket.Conn]bool = make(map[*websocket.Conn]bool)
	userRoles   map[string]string        = map[string]string{
		"user1": "editor",
		"user2": "viewer",
	}
	privateKey = []byte("your-private-key-here")
)

func generateToken(username string, role string) (string, error) {
	// ... (Same as before)
}

func authenticateMiddleware(next http.HandlerFunc) http.HandlerFunc {
	// ... (Same as before)
}

func updateDocument(change DocumentChange, username string, role string) error {
	mu.Lock()
	defer mu.Unlock()

	if role != "editor" {
		return fmt.Errorf("User %s does not have permission to edit the document", username)
	}

	if change.Version != doc.Version {
		return fmt.Errorf("Concurrent edit detected. Please refresh the page.")
	}

	switch change.Type {
	case "insert":
		doc.Content = doc.Content[:change.Index] + change.Content + doc.Content[change.Index:]
	case "delete":
		doc.Content = doc.Content[:change.Index] + doc.Content[change.Index+len(change.Content):]
	default:
		return fmt.Errorf("Unsupported change type: %v", change.Type)
	}

	doc.Version++
	return nil
}

func mergeChanges(change DocumentChange) {
	mu.Lock()
	defer mu.Unlock()

	if change.Version > doc.Version {
		// The server's document is outdated, update it
		doc.Content = change.Content
		doc.Version = change.Version
	} else if change.Version < doc.Version {
		// The client's document is outdated, send the updated document
		sendChangeToAll(DocumentChange{
			Type:    "update",
			Content: doc.Content,
			Version: doc.Version,
		})
	} else {
		// The versions are the same, apply the change
		switch change.Type {
		case "insert":
			doc.Content = doc.Content[:change.Index] + change.Content + doc.Content[change.Index:]
		case "delete":
			doc.Content = doc.Content[:change.Index] + doc.Content[change.Index+len(change.Content):]
		default:
			log.Printf("Unsupported change type: %v", change.Type)
		}

		doc.Version++
	}
}

func sendChangeToAll(change DocumentChange) {
	// ... (Same as before)
}

func handleConnection(w http.ResponseWriter, r *http.Request) {
	// ... (Same as before)
}

func main() {
	// ... (Same as before)
}

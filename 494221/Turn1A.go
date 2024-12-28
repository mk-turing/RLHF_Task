package main

import (
	"log"
	"os"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
)

type User struct {
	Name  string
	Roles []string
}

type Role struct {
	Name  string
	Perms []string // Permissions like "read", "write", "execute"
}

type ACL struct {
	Path      string
	Users     []User // Users allowed access
	Permitted bool
}

var (
	aclList = make(map[string]ACL)
	mutex   sync.Mutex
	roles   = []Role{
		{Name: "admin", Perms: []string{"read", "write", "execute"}},
		{Name: "user", Perms: []string{"read"}},
	}
)

func isUserAllowed(user User, path string, perm string) bool {
	acl, exists := aclList[path]
	if !exists {
		return false // No ACL found, deny by default
	}

	for _, u := range acl.Users {
		if u.Name == user.Name {
			for _, role := range u.Roles {
				for _, p := range getRolePermissions(role) {
					if perm == p {
						return acl.Permitted
					}
				}
			}
		}
	}
	return false
}

func getRolePermissions(role string) []string {
	for _, r := range roles {
		if r.Name == role {
			return r.Perms
		}
	}
	return []string{} // Return empty if role not found
}

func logAccessAttempt(user User, path string, perm string, allowed bool) {
	log.Printf("User: %s, Path: %s, Permission: %s, Allowed: %t\n", user.Name, path, perm, allowed)
	if !allowed {
		log.Printf("Unauthorized access attempt by user %s on %s with permission %s\n", user.Name, path, perm)
	}
}

func watchFile(user User, path string) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatalf("Error creating watcher: %v\n", err)
	}
	defer watcher.Close()

	err = watcher.Add(path)
	if err != nil {
		log.Fatalf("Error watching %s: %v\n", path, err)
	}
	defer watcher.Remove(path)

	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				log.Println("Watcher closed")
				return
			}
			log.Printf("Event: %s\n", event.Name)

			// Check permissions for operation
			allowed := false
			switch event.Op {
			case fsnotify.Write, fsnotify.Create:
				allowed = isUserAllowed(user, path, "write")
			case fsnotify.Remove:
				allowed = isUserAllowed(user, path, "execute") // Simulating need for execute to remove
			default:
				allowed = isUserAllowed(user, path, "read")
			}

			logAccessAttempt(user, path, event.Op.String(), allowed)
		case err := <-watcher.Errors:
			log.Fatalf("Watcher error: %v\n", err)
		}
	}
}

func main() {
	users := []User{
		{Name: "alice", Roles: []string{"admin"}},
		{Name: "bob", Roles: []string{"user"}},
	}

	// Define ACLs
	aclList["test.txt"] = ACL{Path: "test.txt", Users: []User{users[0]}, Permitted: true}

	wg := sync.WaitGroup{}

	// Simulate watching a file with a specific user
	for _, user := range users {
		path := "test.txt"
		wg.Add(1)
		go func() {
			defer wg.Done()
			watchFile(user, path)
		}()
		time.Sleep(1 * time.Second)

		// Simulate actions
		if _, err := os.Stat(path); err == nil {
			f, err := os.OpenFile(path, os.O_WRONLY, 0644)
			if err != nil {
				log.Fatalf("Error opening file: %v\n", err)
			}
			f.Close()
		}
	}

	wg.Wait()
}

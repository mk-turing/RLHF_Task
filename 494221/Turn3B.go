package main

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
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

type LogEntry struct {
	Timestamp time.Time `json:"timestamp"`
	User      string    `json:"user"`
	Path      string    `json:"path"`
	Permission string    `json:"permission"`
	Allowed    bool       `json:"allowed"`
}

const aclStreamKey = "acl-stream"

var (
	aclList  = make(map[string]ACL)
	mutex    sync.Mutex
	roles    = []Role{
		{Name: "admin", Perms: []string{"read", "write", "execute"}},
		{Name: "user", Perms: []string{"read"}},
	}
	logEntries = []LogEntry{}
	mutexLog   sync.Mutex
	rdb        *redis.Client
)

func initRedis() {
	rdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // Redis server address
		Password: "",              // Redis password (set if required)
		DB:       0,                // Redis database number
	})

	pong, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		log.Fatalf("Error pinging Redis: %v", err)
	}
	log.Printf("Connected to Redis: %s", pong)
}

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
	logEntry := LogEntry{
		Timestamp: time.Now(),
		User:      user.Name,
		Path:      path,
		Permission: perm,
		Allowed:    allowed,
	}

	mutexLog.Lock()
	defer mutexLog.Unlock()

	logEntries = append(logEntries, logEntry)

	log.Printf("User: %s, Path: %s, Permission: %s, Allowed: %t\n", user.Name, path, perm, allowed)
	if !allowed {
		log.Printf("Unauthorized access attempt by user %s on %s with permission %s\n", user.Name, path, perm)
	}
}

func watchFile(user User, path string) {
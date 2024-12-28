package main

import (
	"encoding/json"
	"log"
	"os"
	"sync"
	"time"

	"github.com/asdine/storm"
	"github.com/asdine/storm/q"
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

type AccessLog struct {
	ID        int64  `storm:"id,increment"`
	User      string
	Path      string
	Permission string
	Allowed    bool
	Timestamp  time.Time
}

var (
	aclList      = make(map[string]ACL)
	mutex        sync.Mutex
	roles        = []Role{...}
	logDB        *storm.DB
	logMutex     sync.Mutex
)

func init() {
	var err error
	logDB, err = storm.Open("access_logs.db")
	if err != nil {
		log.Fatalf("Error opening log database: %v\n", err)
	}
	if err := logDB.Init(&AccessLog{}); err != nil {
		log.Fatalf("Error initializing log database: %v\n", err)
	}
}

func isUserAllowed(user User, path string, perm string) bool {
	// ... (Same code as before)
}

func logAccessAttempt(user User, path string, perm string, allowed bool) {
	logMutex.Lock()
	defer logMutex.Unlock()

	logEntry := AccessLog{
		User:      user.Name,
		Path:      path,
		Permission: perm,
		Allowed:    allowed,
		Timestamp:  time.Now(),
	}
	if err := logDB.Save(&logEntry); err != nil {
		log.Fatalf("Error saving access log: %v\n", err)
	}
}

func queryLogs(user string, path string, permission string) []AccessLog {
	logMutex.Lock()
	defer logMutex.Unlock()

	var logs []AccessLog
	query := logDB.Select()
	if user != "" {
		query = query.Where("User = ?", user)
	}
	if path != "" {
		query = query.Where("Path = ?", path)
	}
	if permission != "" {
		query = query.Where("Permission = ?", permission)
	}
	if err := query.Find(&logs); err != nil {
		log.Fatalf("Error querying access logs: %v\n", err)
	}
	return logs
}

func main() {
	// ... (Same code as before)

	// Example of logging and querying access logs
	logs := queryLogs("alice", "test.txt", "write")
	for _, log := range logs {
		logBytes, _ := json.Marshal(log)
		log.Println(string(logBytes))
	}
}
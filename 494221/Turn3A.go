package main

import (
	"log"
	"net"
	"net/rpc"
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

type LogEntry struct {
	Timestamp  time.Time `json:"timestamp"`
	User       string    `json:"user"`
	Path       string    `json:"path"`
	Permission string    `json:"permission"`
	Allowed    bool      `json:"allowed"`
}

type Server struct {
	mutex      sync.Mutex
	aclList    map[string]ACL
	roles      []Role
	logEntries []LogEntry
}

func (s *Server) UpdateACL(args *ACL, reply *bool) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.aclList[args.Path] = *args
	*reply = true
	return nil
}

func (s *Server) GetACL(args *string, reply *ACL) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	*reply = s.aclList[*args]
	return nil
}

func (s *Server) ListACLs(args *struct{}, reply *[]ACL) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// Convert map values to slice
	aclSlice := make([]ACL, 0, len(s.aclList))
	for _, acl := range s.aclList {
		aclSlice = append(aclSlice, acl)
	}
	*reply = aclSlice
	return nil
}

func (s *Server) LogAccessAttempt(args *LogEntry, reply *bool) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.logEntries = append(s.logEntries, *args)
	*reply = true
	return nil
}

func (s *Server) ListLogs(args *struct{}, reply *[]LogEntry) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	*reply = s.logEntries
	return nil
}

var (
	aclList = make(map[string]ACL)
	mutex   sync.Mutex
	roles   = []Role{
		{Name: "admin", Perms: []string{"read", "write", "execute"}},
		{Name: "user", Perms: []string{"read"}},
	}
	logEntries = []LogEntry{}
	mutexLog   sync.Mutex
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
	logEntry := LogEntry{
		Timestamp:  time.Now(),
		User:       user.Name,
		Path:       path,
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

	// Start RPC server
	go func() {
		listener, err := net.Listen("tcp", ":1234")
		if err != nil {
			log.Fatalf("Error starting server: %v", err)
		}
		defer listener.Close()

		rpc.Register(&Server{aclList: aclList, roles: roles})
		for {
			conn, err := listener.Accept()
			if err != nil {
				log.Fatalf("Error accepting connection: %v", err)
			}
			go rpc.ServeConn(conn)
		}
	}()

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

	// Example of updating ACL dynamically
	client, err := rpc.Dial("tcp", ":1234")
	if err != nil {
		log.Fatalf("Error dialing RPC client: %v\n", err)
	}
	defer client.Close()

	var updated bool
	updatedACL := ACL{Path: "test.txt", Users: []User{users[1]}, Permitted: true}
	if err := client.Call("Server.UpdateACL", &updatedACL, &updated); err != nil {
		log.Fatalf("Error calling UpdateACL: %v\n", err)
	}
	if updated {
		log.Printf("ACL for 'test.txt' updated successfully.\n")
	} else {
		log.Printf("Failed to update ACL for 'test.txt'.\n")
	}

	// Example of querying ACL
	var queriedACL ACL
	if err := client.Call("Server.GetACL", "test.txt", &queriedACL); err != nil {
		log.Fatalf("Error calling GetACL: %v\n", err)
	}
	log.Printf("ACL for 'test.txt': %v\n", queriedACL)
}

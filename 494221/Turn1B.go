package main

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"log"
	"os"
	"strings"
)

// Role represents the user's role
type Role string

// Permission represents the permissions on a file
type Permission string

const (
	Reader  Permission = "read"
	Writer  Permission = "write"
	Executor Permission = "execute"
)

// User represents the user with a role and permissions on a file
type User struct {
	Role        Role
	FilePerms map[string][]Permission
}

// ACL represents the access control list for a file
type ACL struct {
	Path         string
	AllowedRoles []Role
}

var (
	// users map from user name to User
	users = make(map[string]User)
	// acls map from file path to ACL
	acls = make(map[string]ACL)
)

func main() {
	// Initialize some users and their permissions
	initializeUsersAndPermissions()
	// Initialize some ACLs for files
	initializeACLs()
	// Watch the specified directory for changes
	watchDirectory()
}

func initializeUsersAndPermissions() {
	// Example users and their permissions
	users["alice"] = User{
		Role: "reader",
		FilePerms: map[string][]Permission{
			"/path/to/file1.txt": []Permission{Reader},
			"/path/to/file2.txt": []Permission{Reader, Writer},
		},
	}
	users["bob"] = User{
		Role: "writer",
		FilePerms: map[string][]Permission{
			"/path/to/file2.txt": []Permission{Reader, Writer},
			"/path/to/file3.txt": []Permission{Reader, Writer},
		},
	}
	users["chris"] = User{
		Role: "admin",
		FilePerms: map[string][]Permission{
			"/path/to/file1.txt": []Permission{Reader, Writer, Executor},
			"/path/to/file2.txt": []Permission{Reader, Writer, Executor},
			"/path/to/file3.txt": []Permission{Reader, Writer, Executor},
		},
	}
}

func initializeACLs() {
	// Example ACLs for files
	acls["/path/to/file1.txt"] = ACL{
		Path:         "/path/to/file1.txt",
		AllowedRoles: []Role{"reader", "writer", "admin"},
	}
	acls["/path/to/file2.txt"] = ACL{
		Path:         "/path/to/file2.txt",
		AllowedRoles: []Role{"reader", "writer", "admin"},
	}
	acls["/path/to/file3.txt"] = ACL{
		Path:         "/path/to/file3.txt",
		AllowedRoles: []Role{"reader", "admin"},
	}
}

func watchDirectory() {
	// Create a new file watcher
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal("Error creating file watcher:", err)
	}
	defer watcher.Close()

	// Add the directory to watch for changes
	err = watcher.Add("/path/to/directory")
	if err != nil {
		log.Fatal("Error adding directory to watch:", err)
	}

	done := make(chan bool)

	go func() {
		for {
			select {
			case event := <-watcher.Events:
				handleFileEvent(event)
			case err := <-watcher.Errors:
				log.Println("Error watching file:", err)
			}
		}
	}()

	<-done
}

func handleFileEvent(event fsnotify.Event) {
	// Get the user name and permission from the event
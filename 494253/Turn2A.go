package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

// Role structure to represent roles and permissions
type Role struct {
	Name        string   `json:"name"`
	Permissions []string `json:"permissions"`
}

// User structure to represent users and their roles
type User struct {
	Username string   `json:"username"`
	Roles    []string `json:"roles"`
}

func main() {
	// Define roles and permissions
	roles := []Role{
		{Name: "admin", Permissions: []string{"read", "write", "manage"}},
		{Name: "user", Permissions: []string{"read"}},
		{Name: "guest", Permissions: []string{"none"}},
	}

	// Define users and their roles
	users := []User{
		{Username: "alice", Roles: []string{"admin"}},
		{Username: "bob", Roles: []string{"user"}},
		{Username: "charlie", Roles: []string{"guest"}},
	}

	// Save roles to a file
	if err := saveToFile(roles, "roles.json"); err != nil {
		fmt.Println("Error saving roles:", err)
		return
	}

	// Save users to a file
	if err := saveToFile(users, "users.json"); err != nil {
		fmt.Println("Error saving users:", err)
		return
	}

	// Example of checking permissions for a user
	user := "bob"
	if hasPermission(user, "write", roles, users) {
		fmt.Println(user, "has permission to write.")
	} else {
		fmt.Println(user, "does not have permission to write.")
	}
}

// Helper function to save a struct to a JSON file
func saveToFile(data interface{}, filename string) error {
	dataBytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal data: %w", err)
	}
	return ioutil.WriteFile(filename, dataBytes, 0600)
}

// Helper function to check if a user has a specific permission
func hasPermission(username string, permission string, roles []Role, users []User) bool {
	// Find user roles
	userRoleMap := make(map[string]bool)
	for _, user := range users {
		if user.Username == username {
			for _, role := range user.Roles {
				userRoleMap[role] = true
			}
			break
		}
	}

	// Check user roles against permissions
	for _, role := range roles {
		if userRoleMap[role.Name] {
			for _, perm := range role.Permissions {
				if perm == permission {
					return true
				}
			}
		}
	}

	return false
}

package main

import (
	"log"
	"sync"
	"time"
)

// A simple database-like structure
type Database map[string]string

func (db Database) Set(key, value string) {
	db[key] = value
}

func (db Database) Get(key string) string {
	return db[key]
}

// A callback function type that interacts with the database
type DatabaseCallback func(Database)

func processDatabaseEvent(message string, callback DatabaseCallback) {
	// Simulate database access
	database := make(Database)
	database["key"] = "initial value"

	if callback != nil {
		log.Printf("Executing callback for event: %s\n", message)
		callback(database)
	}
}

func benignCallback(db Database) {
	log.Printf("Benign callback: Reading database value: %s\n", db.Get("key"))
}

func maliciousCallback(db Database) {
	log.Printf("Malicious callback: Overwriting database value\n")
	db.Set("key", "malicious value")
}

func main() {
	log.SetPrefix("app: ")

	// Registering predefined callbacks
	processDatabaseEvent("Event A", benignCallback)

	// Simulate user-defined callback input
	callback, err := getCallbackFromString("maliciousCallback")
	if err != nil {
		log.Fatalf("Error getting callback: %v", err)
	}

	// Process user-defined event (with potential attack)
	processDatabaseEvent("User Event", callback)

	// Add a delay to see the database value change
	time.Sleep(2 * time.Second)

	// Check the database value to detect corruption
	database := make(Database)
	database["key"] = "initial value"
	log.Printf("Final database value: %s\n", database.Get("key"))

	// Implement safeguards like access control and resource isolation
	// Here, we use a mutex to ensure thread safety for database access
	var databaseMutex sync.Mutex

	safeBenignCallback := func(db Database) {
		databaseMutex.Lock()
		defer databaseMutex.Unlock()
		benignCallback(db)
	}

	safeMaliciousCallback := func(db Database) {
		databaseMutex.Lock()
		defer databaseMutex.Unlock()
		maliciousCallback(db)
	}

	// Process database events with safe callbacks (mutex protection)
	processDatabaseEvent("Safe Event A", safeBenignCallback)
	processDatabaseEvent("Safe Event B", safeMaliciousCallback)
}

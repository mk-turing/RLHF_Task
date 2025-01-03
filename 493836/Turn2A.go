package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// Tenant represents a tenant in the system
type Tenant struct {
	ID   int
	Name string
}

// Event represents an event that can be triggered by different tenants
type Event struct {
	TenantID int
	Type     string
	Data     string
}

// tenants stores the list of tenants in the system
var tenants map[int]Tenant

// eventChannels stores the event channels for each tenant
var eventChannels map[int]chan Event

func init() {
	tenants = make(map[int]Tenant)
	eventChannels = make(map[int]chan Event)

	// Initialize some default tenants
	tenants[1] = Tenant{ID: 1, Name: "Tenant A"}
	eventChannels[1] = make(chan Event, 100)

	tenants[2] = Tenant{ID: 2, Name: "Tenant B"}
	eventChannels[2] = make(chan Event, 100)
}

// registerTenantHandler dynamically registers a new tenant
func registerTenantHandler(w http.ResponseWriter, r *http.Request) {
	tenantIDStr := r.URL.Query().Get("tenant_id")
	tenantName := r.URL.Query().Get("tenant_name")

	// Validate and parse tenant ID
	tenantID, err := strconv.Atoi(tenantIDStr)
	if err != nil {
		http.Error(w, "Invalid tenant_id parameter", http.StatusBadRequest)
		return
	}

	// Validate that the tenant ID is unique
	if _, ok := tenants[tenantID]; ok {
		http.Error(w, fmt.Sprintf("Tenant with ID %d already exists", tenantID), http.StatusConflict)
		return
	}

	// Create a new tenant
	tenants[tenantID] = Tenant{ID: tenantID, Name: tenantName}

	// Create an event channel for the new tenant
	eventChannels[tenantID] = make(chan Event, 100)

	fmt.Fprintf(w, "New tenant registered: %s (ID: %d)", tenantName, tenantID)
}

// eventHandler handles incoming events and routes them to the appropriate tenant's event channel
func eventHandler(w http.ResponseWriter, r *http.Request) {
	tenantIDStr := r.URL.Query().Get("tenant_id")
	eventType := r.URL.Query().Get("type")
	eventData := r.URL.Query().Get("data")

	// Validate and parse tenant ID
	tenantID, err := strconv.Atoi(tenantIDStr)
	if err != nil {
		http.Error(w, "Invalid tenant_id parameter", http.StatusBadRequest)
		return
	}

	// Validate tenant exists
	tenant, ok := tenants[tenantID]
	if !ok {
		http.Error(w, fmt.Sprintf("Tenant with ID %d not found", tenantID), http.StatusNotFound)
		return
	}

	// Create an event and send it to the tenant's event channel
	event := Event{
		TenantID: tenant.ID,
		Type:     eventType,
		Data:     eventData,
	}
	eventChannels[tenant.ID] <- event

	fmt.Fprintf(w, "Event received for tenant %s: %+v", tenant.Name, event)
}

// processEventsForTenant processes events for a specific tenant
func processEventsForTenant(tenant Tenant) {
	for event := range eventChannels[tenant.ID] {
		fmt.Printf("Processing event for tenant %s: %+v\n", tenant.Name, event)
		// Handle the event for the specific tenant here
		// You can perform tenant-specific data operations and business logic here
	}
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/event", eventHandler).Methods("POST")
	r.HandleFunc("/register_tenant", registerTenantHandler).Methods("POST")

	// Start goroutines to process events for each existing tenant
	for _, tenant := range tenants {
		go processEventsForTenant(tenant)
	}

	log.Fatal(http.ListenAndServe(":8080", r))
}

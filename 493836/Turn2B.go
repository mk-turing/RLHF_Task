package main
import (
	"fmt"
	"net/http"
	"strconv"
	"sync"

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

var (
	tenants       map[int]Tenant
	eventChannels map[int]chan Event
	mutex         sync.Mutex // Mutex to protect concurrent access to tenants and eventChannels maps
)

func init() {
	tenants = make(map[int]Tenant)
	eventChannels = make(map[int]chan Event)

	// Initialize existing tenants (if any)
	tenants[1] = Tenant{ID: 1, Name: "Tenant A"}
	eventChannels[1] = make(chan Event, 100)

	tenants[2] = Tenant{ID: 2, Name: "Tenant B"}
	eventChannels[2] = make(chan Event, 100)
}

// registerTenantHandler handles tenant registration requests
func registerTenantHandler(w http.ResponseWriter, r *http.Request) {
	tenantName := r.URL.Query().Get("tenant_name")

	if tenantName == "" {
		http.Error(w, "tenant_name parameter is required", http.StatusBadRequest)
		return
	}

	mutex.Lock()
	defer mutex.Unlock()

	// Check if the tenant already exists
	for _, tenant := range tenants {
		if tenant.Name == tenantName {
			http.Error(w, fmt.Sprintf("Tenant with name %s already exists", tenantName), http.StatusConflict)
			return
		}
	}

	// Generate a new tenant ID
	newTenantID := len(tenants) + 1

	// Create the new tenant and its event channel
	tenants[newTenantID] = Tenant{ID: newTenantID, Name: tenantName}
	eventChannels[newTenantID] = make(chan Event, 100)

	fmt.Fprintf(w, "Tenant with name %s registered successfully (ID: %d)", tenantName, newTenantID)
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

	mutex.Lock()
	tenant, ok := tenants[tenantID]
	mutex.Unlock()

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
		// Handle the event for the specific tenant here.
	}
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/event", eventHandler).Methods("POST")
	r.HandleFunc("/register_tenant", registerTenantHandler).Methods("POST")

	// Start goroutines to process events for each existing tenant
	for _, tenant := range tenants {
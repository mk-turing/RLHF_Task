package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
)

// const databaseData = []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z"}
const (
	defaultPageSize  = 10
	maxPageSize      = 100
	paginationParam  = "page"
	pageSizeParam    = "size"
	totalCountHeader = "X-Total-Count"
)

func handlePagination(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	// Parse page and size from query parameters
	pageStr := r.URL.Query().Get(paginationParam)
	pageSizeStr := r.URL.Query().Get(pageSizeParam)
	page := 1
	size := defaultPageSize

	if pageStr != "" {
		var err error
		page, err = strconv.Atoi(pageStr)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "Invalid page parameter: %v", err)
			return
		}
		if page <= 0 {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "Page number must be positive.")
			return
		}
	}

	if pageSizeStr != "" {
		var err error
		size, err = strconv.Atoi(pageSizeStr)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "Invalid page size parameter: %v", err)
			return
		}
		if size <= 0 || size > maxPageSize {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "Page size must be between 1 and %d.", maxPageSize)
			return
		}
	}

	// In a real scenario, you would replace this with your actual data fetching mechanism
	//data, err := fetchDataFromDatabase(page, size)
	data := generateSampleData(page, size)
	//if err != nil {
	//	w.WriteHeader(http.StatusInternalServerError)
	//	fmt.Fprintf(w, "Error fetching data: %v", err)
	//	return
	//}

	totalCount := len(data) // Assuming totalCount is available from the database or other data source
	w.Header().Set(totalCountHeader, strconv.Itoa(totalCount))
	end := time.Now()
	responseTime := end.Sub(start).Milliseconds()

	// Write response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"data": %s, "response_time_ms": %d}`, data, responseTime)
}

func generateSampleData(page, size int) string {
	// simulate database retrieval
	time.Sleep(time.Duration(size) * 100 * time.Millisecond)
	start := (page - 1) * size
	end := min(start+size, len(sampleData))
	return fmt.Sprintf("[\"%s\"]", sampleData[start:end])
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func main() {
	http.HandleFunc("/data", handlePagination)
	log.Println("Server is running on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

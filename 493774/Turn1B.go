package main

import (
	"fmt"
	"net/http"
)

const (
	// TestName1 is the name of the first A/B test
	TestName1 = "test1"
	// Variant1A is the first variant for TestName1
	Variant1A = "variant1a"
	// Variant1B is the second variant for TestName1
	Variant1B = "variant1b"
)

func getVariant(testName, queryParam string) string {
	variants := map[string]string{
		TestName1: Variant1A,
	}

	// Get the variant from the query parameter if provided
	if queryVariant, ok := queryParamMap[testName]; ok {
		switch queryVariant {
		case Variant1A, Variant1B:
			return queryVariant
		}
	}

	// Default to the first variant if not specified or invalid
	return variants[testName]
}

func testHandler(w http.ResponseWriter, r *http.Request) {
	testName := r.URL.Query().Get("test")
	variant := getVariant(testName, r.URL.Query().Get("variant"))

	fmt.Fprintf(w, "Test Name: %s, Variant: %s\n", testName, variant)
}

func main() {
	http.HandleFunc("/", testHandler)
	fmt.Println("Listening on port 8080...")
	http.ListenAndServe(":8080", nil)
}

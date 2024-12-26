package main

import (
	"fmt"
	"log"
	"net/http"
	"sort"
	"strings"
)

type searchParameters struct {
	q        string
	author   string
	category string
	sort     string
}

func main() {
	http.HandleFunc("/search", searchHandler)
	fmt.Println("Server is running on http://localhost:8080/search")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	// Parse URL query parameters
	params := r.URL.Query()
	searchParams := searchParameters{
		q:        params.Get("q"),
		author:   params.Get("author"),
		category: params.Get("category"),
		sort:     params.Get("sort"),
	}

	// Perform the search based on the query parameters
	results := performSearch(searchParams)

	// Display the search results
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, "<h1>Search Results</h1>")
	fmt.Fprintf(w, "<ul>")
	for _, result := range results {
		fmt.Fprintf(w, "<li>%s - %s</li>", result.Title, result.Author)
	}
	fmt.Fprintf(w, "</ul>")
}

func performSearch(params searchParameters) []searchResult {
	// Example search data
	searchData := []searchResult{
		{
			Title:    "The Great Gatsby",
			Author:   "F. Scott Fitzgerald",
			Category: "Classic Literature",
		},
		{
			Title:    "To Kill a Mockingbird",
			Author:   "Harper Lee",
			Category: "Classic Literature",
		},
		{
			Title:    "1984",
			Author:   "George Orwell",
			Category: "Dystopian Fiction",
		},
		{
			Title:    "Harry Potter and the Sorcerer's Stone",
			Author:   "J.K. Rowling",
			Category: "Fantasy",
		},
		{
			Title:    "The Alchemist",
			Author:   "Paulo Coelho",
			Category: "Fantasy",
		},
	}

	var results []searchResult

	// Apply filters and sorting based on query parameters
	for _, result := range searchData {
		if strings.Contains(result.Title, params.q) || strings.Contains(result.Author, params.q) {
			if params.author != "" && result.Author != params.author {
				continue
			}
			if params.category != "" && result.Category != params.category {
				continue
			}
			results = append(results, result)
		}
	}

	switch params.sort {
	case "title":
		sortByTitle(results)
	case "author":
		sortByAuthor(results)
	case "category":
		sortByCategory(results)
	}

	return results
}

type searchResult struct {
	Title    string
	Author   string
	Category string
}

func sortByTitle(results []searchResult) {
	sort.Slice(results, func(i, j int) bool {
		return results[i].Title < results[j].Title
	})
}

func sortByAuthor(results []searchResult) {
	sort.Slice(results, func(i, j int) bool {
		return results[i].Author < results[j].Author
	})
}

func sortByCategory(results []searchResult) {
	sort.Slice(results, func(i, j int) bool {
		return results[i].Category < results[j].Category
	})
}

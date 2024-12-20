package main

import (
	"fmt"
	"sync"
)

type cache struct {
	data map[string]string
	mu   sync.RWMutex
}

func newCache() *cache {
	return &cache{
		data: make(map[string]string),
	}
}

func (c *cache) get(key string) string {
	value, ok := c.data[key]
	if !ok {
		fmt.Println("Cache miss for key:", key)
		// Perform a read-through to fetch the data from the underlying data source
		// For simplicity, let's assume we have a function to fetch data from the source
		value = fetchDataFromSource(key)
		// Store the fetched data in the cache
		c.mu.Lock()
		c.data[key] = value
		c.mu.Unlock()
	}
	return value
}

func (c *cache) put(key string, value string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data[key] = value
	// Perform a write-through to update the underlying data source
	// For simplicity, let's assume we have a function to update the source
	updateDataSource(key, value)
}

func fetchDataFromSource(key string) string {
	// Simulate fetching data from the underlying data source
	return "Value for " + key
}

func updateDataSource(key string, value string) {
	// Simulate updating the underlying data source
	fmt.Println("Updating data source for key:", key, "with value:", value)
}

func main() {
	c := newCache()

	// Initial read-through to populate the cache
	fmt.Println("Cached value for key 'abc':", c.get("abc"))

	// Update the cache
	c.put("abc", "Updated value")

	// Read from the cache
	fmt.Println("Cached value for key 'abc':", c.get("abc"))

	// Read from the cache (misses, triggers read-through)
	fmt.Println("Cached value for key 'xyz':", c.get("xyz"))
}

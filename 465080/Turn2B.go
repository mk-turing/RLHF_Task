package main

import (
	"fmt"
	"github.com/coreos/etcd/proxy/grpcproxy/cache"
	"time"
)

type Cache struct {
	c *cache.Cache
}

func NewCache(ttl time.Duration) *Cache {
	c := cache.New(ttl, 30*time.Second) // Cleanup expired entries every 30 seconds
	return &Cache{c: c}
}

func (c *Cache) get(key string) (string, bool) {
	value, ok := c.c.Get(key)
	if !ok {
		return "", false
	}
	return value.(string), true
}

func (c *Cache) put(key string, value string) {
	c.c.Set(key, value, cache.DefaultExpiration)
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
	cacheTTL := 5 * time.Second // Cache entries will expire after 5 seconds
	c := NewCache(cacheTTL)

	// Initial read-through to populate the cache
	fmt.Println("Cached value for key 'abc':", c.get("abc")) // Miss, triggers read-through

	// Update the cache
	c.put("abc", "Updated value")

	// Read from the cache
	fmt.Println("Cached value for key 'abc':", c.get("abc")) // Hit

	// Wait for the cache entry to expire
	time.Sleep(cacheTTL + 2*time.Second)

	// Check if the expired entry is still in the cache
	value, found := c.get("abc")
	if found {
		fmt.Println("Cached value for key 'abc':", value) // Expired entry should be gone
	} else {
		fmt.Println("Cache miss for key 'abc' (expired)")
	}

	// Trigger a read-through after expiration to refresh the cache
	fmt.Println("Cached value for key 'abc':", c.get("abc")) // Miss, triggers read-through
}

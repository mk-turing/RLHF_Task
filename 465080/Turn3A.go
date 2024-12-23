package main

import (
	"fmt"
	"sync"
	"time"
)

type cache struct {
	data        map[string]string
	expirations map[string]time.Time
	mu          sync.RWMutex
	ttl         time.Duration
	cleanup     *sync.WaitGroup
	stop        chan struct{}
}

func newCache(ttl time.Duration) *cache {
	c := &cache{
		data:        make(map[string]string),
		expirations: make(map[string]time.Time),
		ttl:         ttl,
		cleanup:     &sync.WaitGroup{},
		stop:        make(chan struct{}),
	}

	// Start a background goroutine to clean up expired entries
	c.cleanup.Add(1)
	go c.cleanupExpiredEntries()

	return c
}

func (c *cache) cleanupExpiredEntries() {
	defer c.cleanup.Done()

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-c.stop:
			return
		case <-ticker.C:
			now := time.Now()
			c.mu.Lock()
			for key, expiration := range c.expirations {
				if expiration.Before(now) {
					delete(c.data, key)
					delete(c.expirations, key)
					fmt.Println("Expired key:", key)
				}
			}
			c.mu.Unlock()
		}
	}
}

func (c *cache) get(key string) string {
	c.mu.RLock()
	value, ok := c.data[key]
	c.mu.RUnlock()
	if !ok {
		fmt.Println("Cache miss for key:", key)
		// Perform a read-through to fetch the data from the underlying data source
		value = fetchDataFromSource(key)
		// Store the fetched data in the cache with a TTL
		c.mu.Lock()
		c.data[key] = value
		c.expirations[key] = time.Now().Add(c.ttl)
		c.mu.Unlock()
		return value
	}

	expiration := c.expirations[key]
	if expiration.Before(time.Now()) {
		c.mu.Lock()
		delete(c.data, key)
		delete(c.expirations, key)
		c.mu.Unlock()
		fmt.Println("Expired key:", key)
		return c.get(key) // Recursively retrieve the value after handling expiration
	}

	return value
}

func (c *cache) put(key string, value string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.data[key] = value
	c.expirations[key] = time.Now().Add(c.ttl)
	// Perform a write-through to update the underlying data source
	updateDataSource(key, value)
}

func (c *cache) Close() {
	close(c.stop)
	c.cleanup.Wait()
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
	c := newCache(5 * time.Second)
	defer c.Close()

	// Initial read-through to populate the cache
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		fmt.Println("Cached value for key 'abc':", c.get("abc"))
		wg.Done()
	}()
	wg.Add(1)
	go func() {
		fmt.Println("Cached value for key 'abc':", c.get("abc"))
		wg.Done()
	}()
	wg.Wait()

	// Update the cache
	wg.Add(1)
	go func() {
		c.put("abc", "Updated value")
		wg.Done()
	}()
	wg.Wait()

	// Read from the cache
	wg.Add(1)
	go func() {
		fmt.Println("Cached value for key 'abc':", c.get("abc"))
		wg.Done()
	}()
	wg.Wait()

	// Read from the cache (misses, triggers read-through)
	wg.Add(1)
	go func() {
		fmt.Println("Cached value for key 'xyz':", c.get("xyz"))
		wg.Done()
	}()
	wg.Wait()

	// Wait for some time to see cache expiration
	time.Sleep(6 * time.Second)

	// Attempt to read an expired key
	wg.Add(1)
	go func() {
		fmt.Println("Cached value for key 'abc':", c.get("abc"))
		wg.Done()
	}()
	wg.Wait()
}
package main

import (
	"fmt"
	"sync"
	"time"
)

type cache struct {
	data        map[string]*entry
	expirations map[string]time.Time
	ttl         time.Duration
	capacity    int
	mu          sync.RWMutex
	cleanup     *sync.WaitGroup
	stop        chan struct{}
	head, tail  *entry
}

type entry struct {
	key        string
	value      string
	expiration time.Time
	next, prev *entry
}

func newCache(ttl time.Duration, capacity int) *cache {
	c := &cache{
		data:        make(map[string]*entry),
		expirations: make(map[string]time.Time),
		ttl:         ttl,
		capacity:    capacity,
		mu:          sync.RWMutex{},
		cleanup:     &sync.WaitGroup{},
		stop:        make(chan struct{}),
	}

	// Initialize head and tail of the doubly linked list
	c.head = &entry{}
	c.tail = &entry{}
	c.head.next = c.tail
	c.tail.prev = c.head

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
					c.removeEntry(c.data[key])
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
	entry, ok := c.data[key]
	c.mu.RUnlock()
	if !ok {
		fmt.Println("Cache miss for key:", key)
		// Perform a read-through to fetch the data from the underlying data source
		value := fetchDataFromSource(key)
		// Store the fetched data in the cache with a TTL
		c.mu.Lock()
		entry = c.addEntry(key, value)
		c.mu.Unlock()
		return value
	}

	expiration := entry.expiration
	if expiration.Before(time.Now()) {
		c.mu.Lock()
		c.removeEntry(entry)
		delete(c.data, key)
		delete(c.expirations, key)
		fmt.Println("Expired key:", key)
		c.mu.Unlock()
		return c.get(key) // Recursively retrieve the value after handling expiration
	}

	// Move the accessed entry to the front of the list
	c.mu.Lock()
	c.moveToHead(entry)
	c.mu.Unlock()

	return entry.value
}

func (c *cache) put(key string, value string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if entry, ok := c.data[key]; ok {
		entry.value = value
		entry.expiration = time.Now().Add(c.ttl)
		c.moveToHead(entry)
		// Perform a write-through to update the underlying data source
		updateDataSource(key, value)
	} else {
		entry := c.addEntry(key, value)
		entry.expiration = time.Now().Add(c.ttl)
		if c.capacity > 0 && len(c.data) > c.capacity {
			oldest := c.removeTail()
			delete(c.data, oldest.key)
			delete(c.expirations, oldest.key)
		}
		// Perform a write-through to update the underlying data source
		updateDataSource(key, value)
	}
}

func (c *cache) addEntry(key, value string) *entry {
	entry := &entry{key: key, value: value}
	c.data[key] = entry
	c.addToHead(entry)
	return entry
}

func (c *cache) addToHead(entry *entry) {
	entry.next = c.head.next
	entry.prev = c.head
	c.head.next.prev = entry
	c.head.next = entry
}

func (c *cache) moveToHead(entry *entry) {
	c.removeEntry(entry)
	c.addToHead(entry)
}

func (c *cache) removeEntry(entry *entry) {
	entry.prev.next = entry.next
	entry.next.prev = entry.prev
}

func (c *cache) removeTail() *entry {
	tail := c.tail.prev
	c.removeEntry(tail)
	return tail
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
	c := newCache(5*time.Second, 2)
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
		fmt.Println("Cached value for key 'def':", c.get("def"))
		wg.Done()
	}()
	wg.Wait()

	// Update the cache
	wg.Add(1)
	go func() {
		c.put("abc", "Updated value for abc")
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
		fmt.Println("Cached value for key 'ghi':", c.get("ghi"))
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

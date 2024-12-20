package main

import (
	"fmt"
	"sync"
	"time"
)

// Doubly linked list node to track the order of access
type cacheNode struct {
	key   string
	value string
	prev  *cacheNode
	next  *cacheNode
}

type cache struct {
	size        int
	ttl         time.Duration
	data        map[string]*cacheNode
	expirations map[string]time.Time
	head        *cacheNode
	tail        *cacheNode
	mu          sync.RWMutex
	cleanup     *sync.WaitGroup
	stop        chan struct{}
}

func newCache(size int, ttl time.Duration) *cache {
	c := &cache{
		size:        size,
		ttl:         ttl,
		data:        make(map[string]*cacheNode),
		expirations: make(map[string]time.Time),
		mu:          sync.RWMutex{},
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
					c.evictNode(c.data[key])
					delete(c.data, key)
					delete(c.expirations, key)
				}
			}
			c.mu.Unlock()
		}
	}
}

// Moves a node to the front of the doubly linked list to mark it as the most recently used
func (c *cache) moveToFront(node *cacheNode) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if node == c.head {
		return
	}
	if node == c.tail {
		c.tail = node.prev
		c.tail.next = nil
	} else {
		node.prev.next = node.next
		node.next.prev = node.prev
	}
	node.next = c.head
	node.prev = nil
	c.head.prev = node
	c.head = node
}

// Evicts the least recently used node from the cache
func (c *cache) evictNode(node *cacheNode) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if node == c.head {
		c.head = node.next
		c.head.prev = nil
	} else if node == c.tail {
		c.tail = node.prev
		c.tail.next = nil
	} else {
		node.prev.next = node.next
		node.next.prev = node.prev
	}
	delete(c.data, node.key)
}

func (c *cache) get(key string) string {
	c.mu.RLock()
	defer c.mu.RUnlock()

	node, ok := c.data[key]
	if !ok {
		fmt.Println("Cache miss for key:", key)
		// Perform a read-through to fetch the data from the underlying data source
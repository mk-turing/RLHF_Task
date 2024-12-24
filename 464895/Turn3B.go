package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"
)

// CacheItem represents a single entry in the cache.
type CacheItem struct {
	Key    string
	Data   interface{}
	Expiry time.Time
}

// FileCache represents the file-based cache.
type FileCache struct {
	filename     string
	cache        map[string]CacheItem
	cacheMutex   sync.Mutex
	writeChannel chan CacheItem
	maxSize      int
	expiryTicker *time.Ticker
}

// NewFileCache creates a new FileCache instance.
func NewFileCache(filename string, maxSize int) *FileCache {
	cache := &FileCache{
		filename:     filename,
		cache:        make(map[string]CacheItem),
		maxSize:      maxSize,
		writeChannel: make(chan CacheItem, 100), // Buffered channel for batching writes
		expiryTicker: time.NewTicker(10 * time.Second),
	}

	// Load cache from file if it exists
	cache.loadCacheFromFile()

	// Start a goroutine to handle cache writes and expiration
	go cache.handleWritesAndExpiry()

	return cache
}

// loadCacheFromFile loads the cache from the file.
func (cache *FileCache) loadCacheFromFile() {
	file, err := os.Open(cache.filename)
	if err != nil {
		if os.IsNotExist(err) {
			return
		}
		fmt.Printf("Error loading cache from file: %v\n", err)
		return
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&cache.cache); err != nil {
		fmt.Printf("Error decoding cache from file: %v\n", err)
		return
	}
}

// saveCacheToFile saves the cache to the file.
func (cache *FileCache) saveCacheToFile() {
	file, err := os.Create(cache.filename)
	if err != nil {
		fmt.Printf("Error creating cache file: %v\n", err)
		return
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	if err := encoder.Encode(cache.cache); err != nil {
		fmt.Printf("Error encoding cache to file: %v\n", err)
		return
	}
}

// Get retrieves a value from the cache.
func (cache *FileCache) Get(key string) (interface{}, bool) {
	cache.cacheMutex.Lock()
	defer cache.cacheMutex.Unlock()

	item, exists := cache.cache[key]
	if !exists {
		return nil, false
	}

	if time.Now().After(item.Expiry) {
		delete(cache.cache, key)
		return nil, false
	}

	return item.Data, true
}

// Set sets a value in the cache with an expiration time.
func (cache *FileCache) Set(key string, data interface{}, ttl time.Duration) {
	expiry := time.Now().Add(ttl)
	cache.writeChannel <- CacheItem{Key: key, Data: data, Expiry: expiry}
}

// Delete deletes a value from the cache.
func (cache *FileCache) Delete(key string) {
	cache.writeChannel <- CacheItem{Key: key}
}

// handleWritesAndExpiry handles cache writes and expiration asynchronously.
func (cache *FileCache) handleWritesAndExpiry() {
	for {
		select {
		case item := <-cache.writeChannel:
			cache.cacheMutex.Lock()
			if item.Key != "" {
				cache.cache[item.Key] = item
			} else {
				delete(cache.cache, item.Key)
			}
			cache.cacheMutex.Unlock()

			// Batch writes to the file every 100 writes or after 5 seconds
			if len(cache.writeChannel) == 0 || time.Now().After(time.Now().Add(5*time.Second)) {
				cache.saveCacheToFile()
			}

		case <-cache.expiryTicker.C:
			cache.cacheMutex.Lock()
			now := time.Now()
			expiredKeys := make([]string, 0)

			// Collect expired keys
			for key, item := range cache.cache {
				if now.After(item.Expiry) {
					expiredKeys = append(expiredKeys, key)
				}
			}

			// Delete expired keys
			for _, key := range expiredKeys {
				delete(cache.cache, key)
			}
			cache.cacheMutex.Unlock()

			cache.saveCacheToFile()
		}
	}
}

func main() {
	cache := NewFileCache("turn3B.json", 10)

	// Example usage with concurrent reads and writes
	go func() {
		for i := 0; i < 100; i++ {
			key := fmt.Sprintf("key%d", i)
			cache.Set(key, fmt.Sprintf("value%d", i), 1*time.Minute)
		}
	}()

	go func() {
		for i := 0; i < 100; i++ {
			key := fmt.Sprintf("key%d", i)
			value, _ := cache.Get(key)
			fmt.Printf("Retrieved value for %s: %v\n", key, value)
		}
	}()

	time.Sleep(5 * time.Second)
}

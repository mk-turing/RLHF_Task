package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"
)

type CacheItem struct {
	Key    string
	Data   interface{}
	Expiry time.Time
}

type FileCache struct {
	filename     string
	cache        map[string]CacheItem
	cacheMutex   sync.Mutex
	maxSize      int
	expiryTicker *time.Ticker
}

func NewFileCache(filename string, maxSize int) *FileCache {
	cache := &FileCache{
		filename:     filename,
		cache:        make(map[string]CacheItem),
		maxSize:      maxSize,
		expiryTicker: time.NewTicker(10 * time.Second),
	}

	// Load cache from file if it exists
	cache.loadCacheFromFile()

	// Start a goroutine to handle cache expiration
	go cache.handleExpiry()

	return cache
}

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

func (cache *FileCache) Get(key string) (interface{}, bool) {
	cache.cacheMutex.Lock()
	defer cache.cacheMutex.Unlock()

	item, exists := cache.cache[key]
	if !exists {
		return nil, false
	}

	if time.Now().After(item.Expiry) {
		cache.deleteItem(key)
		return nil, false
	}

	return item.Data, true
}

func (cache *FileCache) Set(key string, data interface{}, ttl time.Duration) {
	cache.cacheMutex.Lock()
	defer cache.cacheMutex.Unlock()

	expiry := time.Now().Add(ttl)
	cache.cache[key] = CacheItem{
		Key:    key,
		Data:   data,
		Expiry: expiry,
	}

	cache.saveCacheToFile()

	cache.removeOldestEntries()
}

func (cache *FileCache) Delete(key string) {
	cache.cacheMutex.Lock()
	defer cache.cacheMutex.Unlock()

	cache.deleteItem(key)
	cache.saveCacheToFile()
}

func (cache *FileCache) deleteItem(key string) {
	delete(cache.cache, key)
}

func (cache *FileCache) removeOldestEntries() {
	length := len(cache.cache)
	if length > cache.maxSize {
		keys := make([]string, 0, length)
		for k := range cache.cache {
			keys = append(keys, k)
		}
		for i := 0; i < length-cache.maxSize; i++ {
			cache.deleteItem(keys[i])
		}
	}
}

func (cache *FileCache) handleExpiry() {
	for range cache.expiryTicker.C {
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
			cache.deleteItem(key)
		}
		cache.cacheMutex.Unlock()

		cache.saveCacheToFile()
	}
}

func main() {
	cache := NewFileCache("turn2A.json", 10)

	// Example usage
	cache.Set("key1", "value1", 1*time.Minute)
	data, exists := cache.Get("key1")
	if exists {
		fmt.Printf("Retrieved value for key1: %v\n", data)
	} else {
		fmt.Printf("Key1 not found in cache\n")
	}

	cache.Set("key2", "value2", 5*time.Second)
	time.Sleep(6 * time.Second)
	data, exists = cache.Get("key2")
	if exists {
		fmt.Printf("Retrieved value for key2: %v\n", data)
	} else {
		fmt.Printf("Key2 not found in cache\n")
	}

	cache.Set("key3", "value3", 30*time.Second)
	time.Sleep(10 * time.Second)
	data, exists = cache.Get("key3")
	if exists {
		fmt.Printf("Retrieved value for key3: %v\n", data)
	} else {
		fmt.Printf("Key3 not found in cache\n")
	}

	cache.Delete("key1")
	time.Sleep(2 * time.Second)
	data, exists = cache.Get("key1")
	if exists {
		fmt.Printf("Retrieved value for key1: %v\n", data)
	} else {
		fmt.Printf("Key1 not found in cache\n")
	}
}

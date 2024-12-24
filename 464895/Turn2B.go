package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sync"
	"time"
)

type CacheItem struct {
	Data   string    `json:"data"`
	Expiry time.Time `json:"expiry"`
}

type FileCache struct {
	filename   string
	cache      map[string]CacheItem
	cacheMutex sync.Mutex
}

func NewFileCache(filename string) *FileCache {
	cache := &FileCache{
		filename: filename,
		cache:    make(map[string]CacheItem),
	}
	cache.loadCacheFromFile()
	return cache
}

func (cache *FileCache) loadCacheFromFile() {
	data, err := ioutil.ReadFile(cache.filename)
	if err != nil {
		if os.IsNotExist(err) {
			return // File not found, so assume empty cache
		}
		fmt.Printf("Error reading cache file: %v\n", err)
		return
	}

	if err := json.Unmarshal(data, &cache.cache); err != nil {
		fmt.Printf("Error decoding cache from file: %v\n", err)
		return
	}
}

func (cache *FileCache) saveCacheToFile() {
	data, err := json.Marshal(cache.cache)
	if err != nil {
		fmt.Printf("Error encoding cache to file: %v\n", err)
		return
	}

	if err := ioutil.WriteFile(cache.filename, data, 0644); err != nil {
		fmt.Printf("Error writing cache to file: %v\n", err)
	}
}

func (cache *FileCache) Get(key string) (string, bool) {
	cache.cacheMutex.Lock()
	defer cache.cacheMutex.Unlock()

	item, exists := cache.cache[key]
	if exists && time.Now().After(item.Expiry) {
		delete(cache.cache, key)
		cache.saveCacheToFile()
		return "", false
	}

	return item.Data, exists
}

func (cache *FileCache) Set(key string, data string, ttl time.Duration) {
	cache.cacheMutex.Lock()
	defer cache.cacheMutex.Unlock()

	expiry := time.Now().Add(ttl)
	cache.cache[key] = CacheItem{
		Data:   data,
		Expiry: expiry,
	}
	cache.saveCacheToFile()
}

func (cache *FileCache) Delete(key string) {
	cache.cacheMutex.Lock()
	defer cache.cacheMutex.Unlock()

	delete(cache.cache, key)
	cache.saveCacheToFile()
}

func main() {
	cache := NewFileCache("turn2B.json")

	// Set multiple cache entries with different expiration times
	cache.Set("key1", "value1 with TTL 1 minute", 1*time.Minute)
	cache.Set("key2", "value2 with TTL 2 hours", 2*time.Hour)
	cache.Set("key3", "value3 with TTL 10 seconds", 10*time.Second)

	// Retrieve and display cache entries
	value, found := cache.Get("key1")
	if found {
		fmt.Printf("Retrieved value for key1: %s\n", value)
	} else {
		fmt.Printf("Key1 not found in cache\n")
	}

	// ...
	// After 10 seconds, key3 will expire and won't be found in the cache
}

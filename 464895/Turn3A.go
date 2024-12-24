package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"sync/atomic"
	"time"
)

type CacheItem struct {
	Key    string
	Data   interface{}
	Expiry time.Time
}

type FileCache struct {
	filename     string
	cache        sync.Map // Concurrent map
	maxSize      int
	dirtyCounter uint64 // Track if cache needs saving
	expiryTicker *time.Ticker
}

func NewFileCache(filename string, maxSize int) *FileCache {
	cache := &FileCache{
		filename:     filename,
		cache:        sync.Map{},
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
	var cacheItems []CacheItem
	if err := decoder.Decode(&cacheItems); err != nil {
		fmt.Printf("Error decoding cache from file: %v\n", err)
		return
	}

	for _, item := range cacheItems {
		cache.cache.Store(item.Key, item)
	}
}

func (cache *FileCache) get(key string) (CacheItem, bool) {
	var item CacheItem
	if cache.cache.Load(key); item.Expiry.Before(time.Now()) {
		cache.cache.Delete(key)
		return CacheItem{}, false
	}
	return item, true
}

func (cache *FileCache) getWithSync(key string) (CacheItem, bool) {
	for {
		item, ok := cache.get(key)
		if ok {
			return item, ok
		}
		time.Sleep(time.Millisecond) // Expiry could be between load check and expiration
	}
}

func lenSyncMap(m *sync.Map) int {
	var i int
	m.Range(func(k, v interface{}) bool {
		i++
		return true
	})
	return i
}

func (cache *FileCache) set(key string, data interface{}, ttl time.Duration) {
	expiry := time.Now().Add(ttl)
	item := CacheItem{Key: key, Data: data, Expiry: expiry}
	cache.cache.Store(key, item)
	atomic.AddUint64(&cache.dirtyCounter, 1) // Increment dirty counter to indicate need for save

	if lenSyncMap(&cache.cache) > cache.maxSize {
		cache.evictOldestItem()
	}
}

func (cache *FileCache) evictOldestItem() {
	var firstKey string
	cache.cache.Range(func(k, _ interface{}) bool {
		if firstKey == "" {
			firstKey = k.(string)
		} else {
			cache.cache.Delete(k)
		}
		return false
	})
}

func (cache *FileCache) getFileBuffer() ([]byte, error) {
	items := make([]CacheItem, 0, lenSyncMap(&cache.cache))
	cache.cache.Range(func(k, v interface{}) bool {
		items = append(items, v.(CacheItem))
		return true
	})
	return json.Marshal(items)
}

func (cache *FileCache) saveCacheToFile() {
	buffer, err := cache.getFileBuffer()
	if err != nil {
		fmt.Printf("Error encoding cache to buffer: %v\n", err)
		return
	}

	file, err := os.Create(cache.filename)
	if err != nil {
		fmt.Printf("Error creating cache file: %v\n", err)
		return
	}
	defer file.Close()

	_, err = file.Write(buffer)
	if err != nil {
		fmt.Printf("Error writing cache to file: %v\n", err)
		return
	}
	atomic.SwapUint64(&cache.dirtyCounter, 0) // Reset dirty counter
}

func (cache *FileCache) handleExpiry() {
	for range cache.expiryTicker.C {
		// Only write if the cache is dirty
		if atomic.LoadUint64(&cache.dirtyCounter) > 0 {
			cache.saveCacheToFile()
		}
	}
}

func main() {
	cache := NewFileCache("turn3A.json", 10)

	// Example usage with concurrency
	wg := sync.WaitGroup{}

	for i := 0; i < 10; i++ {
		key := fmt.Sprintf("key%d", i)
		value := fmt.Sprintf("value%d", i)

		wg.Add(1)
		go func() {
			cache.set(key, value, 5*time.Second)
			defer wg.Done()
		}()

		wg.Add(1)
		go func() {
			for {
				data, exists := cache.getWithSync(key)
				if exists {
					fmt.Printf("Retrieved value for %s: %v\n", key, data.Data)
				} else {
					break
				}
				time.Sleep(time.Second)
			}
			defer wg.Done()
		}()
	}

	wg.Wait()
}

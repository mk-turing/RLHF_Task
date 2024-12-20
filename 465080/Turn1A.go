package main

import (
	"fmt"
	"sync"
)

type DataSource interface {
	Get(key string) (string, error)
}

type SimpleDataSource struct {
	data map[string]string
}

func (ds *SimpleDataSource) Get(key string) (string, error) {
	return ds.data[key], nil
}

func NewSimpleDataSource() *SimpleDataSource {
	return &SimpleDataSource{data: make(map[string]string)}
}

type ReadThroughCache struct {
	dataSource DataSource
	cache      sync.Map
}

func NewReadThroughCache(dataSource DataSource) *ReadThroughCache {
	return &ReadThroughCache{dataSource: dataSource, cache: sync.Map{}}
}

func (rc *ReadThroughCache) Get(key string) (string, error) {
	if value, loaded := rc.cache.Load(key); loaded {
		return value.(string), nil
	}

	value, err := rc.dataSource.Get(key)
	if err != nil {
		return "", err
	}

	rc.cache.Store(key, value)
	return value, nil
}

func main() {
	ds := NewSimpleDataSource()
	ds.data["key1"] = "value1"
	ds.data["key2"] = "value2"

	cache := NewReadThroughCache(ds)

	value1, _ := cache.Get("key1")
	value2, _ := cache.Get("key2")
	value3, _ := cache.Get("key1") // Read from cache

	fmt.Printf("Value 1: %s\n", value1)
	fmt.Printf("Value 2: %s\n", value2)
	fmt.Printf("Value 3 (from cache): %s\n", value3)
}

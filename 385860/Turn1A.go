package main

import (
	"fmt"
	"reflect"
)

// SampleData is the struct representing your large dataset.
type SampleData struct {
	Nums []int
	Map  map[string]int
}

// createSnapshot creates a deep copy of a struct or a pointer to a struct using reflection.
func createSnapshot(v interface{}) interface{} {
	rv := reflect.ValueOf(v)

	// Ensure v is a struct or a pointer to a struct
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}
	if rv.Kind() != reflect.Struct {
		panic("snapshot can only be created for structs or pointers to structs")
	}

	// Create a copy of the struct
	nv := reflect.New(rv.Type()).Elem()

	// Copy all the fields
	for i := 0; i < rv.NumField(); i++ {
		sf := rv.Field(i)
		nf := nv.Field(i)

		// Handle pointers
		if sf.Kind() == reflect.Ptr {
			if sf.IsNil() {
				nf.Set(reflect.Zero(sf.Type()))
			} else {
				ptrCopy := reflect.New(sf.Elem().Type())
				ptrCopy.Elem().Set(sf.Elem())
				nf.Set(ptrCopy)
			}
		} else if sf.Kind() == reflect.Slice {
			// Handle slices
			if !sf.IsNil() {
				sliceCopy := reflect.MakeSlice(sf.Type(), sf.Len(), sf.Cap())
				reflect.Copy(sliceCopy, sf)
				nf.Set(sliceCopy)
			}
		} else if sf.Kind() == reflect.Map {
			// Handle maps
			if !sf.IsNil() {
				mapCopy := reflect.MakeMap(sf.Type())
				for _, key := range sf.MapKeys() {
					mapCopy.SetMapIndex(key, sf.MapIndex(key))
				}
				nf.Set(mapCopy)
			}
		} else {
			// Handle primitive types
			nf.Set(sf)
		}
	}

	return nv.Addr().Interface() // Return as a pointer to the struct
}

func main() {
	// Original data
	originalData := &SampleData{
		Nums: []int{1, 2, 3, 4, 5},
		Map: map[string]int{
			"a": 1,
			"b": 2,
		},
	}

	// Create a snapshot
	snapshot := createSnapshot(originalData)

	// Modify the original data
	originalData.Nums[2] = 10
	originalData.Map["c"] = 3

	fmt.Println("Original Data:")
	fmt.Printf("Nums: %v\n", originalData.Nums)
	fmt.Printf("Map: %v\n", originalData.Map)

	fmt.Println("\nSnapshot:")
	snapshotData := snapshot.(*SampleData)
	fmt.Printf("Nums: %v\n", snapshotData.Nums)
	fmt.Printf("Map: %v\n", snapshotData.Map)
}

package main

import (
	"fmt"
	"reflect"
)

func customFormat(v interface{}) string {
	rv := reflect.ValueOf(v)
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}

	switch rv.Kind() {
	case reflect.Map:
		return formatMap(rv)
	case reflect.Struct:
		return formatStruct(rv)
	default:
		return fmt.Sprintf("%v", v)
	}
}

func formatMap(rv reflect.Value) string {
	keys := rv.MapKeys()
	pairs := make([]string, len(keys))
	for i, key := range keys {
		value := rv.MapIndex(key)
		pairs[i] = fmt.Sprintf("%v: %v", key.Interface(), customFormat(value.Interface()))
	}
	return fmt.Sprintf("map[%s]", fmt.Sprintf("%s", pairs))
}

func formatStruct(rv reflect.Value) string {
	numField := rv.NumField()
	fields := make([]string, numField)
	for i := 0; i < numField; i++ {
		field := rv.Field(i)
		tag := rv.Type().Field(i).Tag.Get("json")
		if tag == "-" {
			continue
		}
		fields[i] = fmt.Sprintf("%s: %v", tag, customFormat(field.Interface()))
	}
	return fmt.Sprintf("{%s}", fmt.Sprintf("%s", fields))
}

type Person struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

type Address struct {
	Street string `json:"street"`
	City   string `json:"city"`
}

func main() {
	person := Person{Name: "John Doe", Age: 30}
	address := Address{Street: "Main Street", City: "New York"}
	nestedMap := map[string]interface{}{
		"person":  person,
		"address": address,
	}

	fmt.Println("Custom Format:", customFormat(nestedMap))
	// Output: Custom Format: map[address:{street:Main Street city:New York} person:{name:John Doe age:30}]
}

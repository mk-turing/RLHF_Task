package main

import (
	"fmt"
)

type Printer interface {
	String() string
}

type Point struct {
	X, Y int
}

func (p Point) String() string {
	return fmt.Sprintf("Point{%d, %d}", p.X, p.Y)
}

func main() {
	var mapData map[string]Printer = map[string]Printer{
		"origin":      Point{0, 0},
		"destination": Point{10, 20},
	}

	fmt.Printf("Map: %v\n", mapData)
}

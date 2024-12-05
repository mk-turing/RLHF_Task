package main

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
)

func main() {
	data := []byte("Hello, this is a large dataset!")

	// Compress data
	var b bytes.Buffer
	w := gzip.NewWriter(&b)
	_, err := w.Write(data)
	if err != nil {
		panic(err)
	}
	err = w.Close()
	if err != nil {
		panic(err)
	}
	compressedData := b.Bytes()

	// ... Perform operations on compressed data ...

	// Decompress data (if needed)
	var decompressedData bytes.Buffer
	r, err := gzip.NewReader(bytes.NewReader(compressedData))
	if err != nil {
		panic(err)
	}
	_, err = io.Copy(&decompressedData, r)
	if err != nil {
		panic(err)
	}
	err = r.Close()
	if err != nil {
		panic(err)
	}
	fmt.Println(string(decompressedData.Bytes()))
}

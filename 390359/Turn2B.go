package main

import (
	"encoding/json"
	"fmt"
	"google.golang.org/protobuf/proto"
	pb "main/pb/userpb" // Adjust the import path
	"testing"
)

type UserJSON struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

var userProto = &pb.User{Id: 1, Name: "John Doe"}
var userJSON = &UserJSON{ID: 1, Name: "John Doe"}

func BenchmarkJSONDecode(b *testing.B) {
	userBytes, err := json.Marshal(userJSON)
	if err != nil {
		b.Fatalf("error marshalling JSON: %v", err)
	}

	for i := 0; i < b.N; i++ {
		var decodedUser UserJSON
		if err := json.Unmarshal(userBytes, &decodedUser); err != nil {
			b.Fatalf("error unmarshalling JSON: %v", err)
		}
	}
}

func BenchmarkProtobufDecode(b *testing.B) {
	userBytes, err := proto.Marshal(userProto)
	if err != nil {
		b.Fatalf("error marshalling protobuf: %v", err)
	}

	for i := 0; i < b.N; i++ {
		var decodedUser pb.User
		if err := proto.Unmarshal(userBytes, &decodedUser); err != nil {
			b.Fatalf("error unmarshalling protobuf: %v", err)
		}
	}
}

func main() {
	fmt.Println("Running benchmarks...")
	testing.Benchmark(BenchmarkJSONDecode)
	testing.Benchmark(BenchmarkProtobufDecode)
}

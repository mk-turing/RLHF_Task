// benchmark.go
package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"net/http"
	"time"

	"google.golang.org/grpc"
	pb "main/pb/userpb"
)

const (
	restURL  = "http://localhost:8080/user"
	grpcAddr = "localhost:50051"
)

func benchmarkREST(numRequests int) {
	start := time.Now()

	for i := 0; i < numRequests; i++ {
		user := map[string]interface{}{
			"id":   i,
			"name": uuid.New().String(),
		}
		body, _ := json.Marshal(user)
		http.Post(restURL, "application/json", bytes.NewBuffer(body))
	}

	elapsed := time.Since(start)
	fmt.Printf("REST took %s for %d requests\n", elapsed, numRequests)
}

func benchmarkGPRC(numRequests int) {
	conn, err := grpc.Dial(grpcAddr, grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	client := pb.NewUserServiceClient(conn)

	start := time.Now()

	for i := 0; i < numRequests; i++ {
		_, err := client.CreateUser(context.Background(), &pb.User{Id: int32(i), Name: uuid.New().String()})
		if err != nil {
			panic(err)
		}
	}

	elapsed := time.Since(start)
	fmt.Printf("gRPC took %s for %d requests\n", elapsed, numRequests)
}

func main() {
	numRequests := 1000

	go func() {
		benchmarkREST(numRequests)
	}()

	go func() {
		benchmarkGPRC(numRequests)
	}()

	time.Sleep(10 * time.Second) // Wait for both benchmarks to complete
}

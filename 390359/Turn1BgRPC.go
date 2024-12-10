package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
	pb "main/pb/userpb" // Adjust the path to your generated protobuf package
)

type userService struct {
	pb.UnimplementedUserServiceServer
}

func (s *userService) CreateUser(ctx context.Context, req *pb.User) (*pb.User, error) {
	fmt.Printf("User created: %+v\n", req)
	return req, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterUserServiceServer(s, &userService{})

	log.Println("server starting on :50051")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

package main

import (
	"context"
	"log"
	"net"

	"google.golang.org/grpc"
	pb "main/pb/userpb" // Import your generated protobuf code
)

// server is used to implement user.UserServiceServer
type server struct {
	pb.UnimplementedUserServiceServer
}

// CreateUser implements user.UserServiceServer
func (s *server) CreateUser(ctx context.Context, user *pb.User) (*pb.User, error) {
	// Process user (e.g., create or update)
	// ...

	// Echo back the received user
	return user, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterUserServiceServer(s, &server{})
	log.Println("Server listening on port 50051")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

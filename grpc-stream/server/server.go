package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"time"

	"google.golang.org/grpc"

	pb "github.com/nethish/playground/grpc-stream/gen/api"
)

type userServiceServer struct {
	pb.UnimplementedUserServiceServer
}

// Server stream
func (s *userServiceServer) ListUsers(req *pb.ListUsersRequest, stream pb.UserService_ListUsersServer) error {
	log.Printf("ListUsers called with filter: %s", req.Filter)

	// Simulate streaming users
	for i := 1; i <= 5; i++ {
		user := &pb.User{
			Id:    fmt.Sprintf("user-%d", i),
			Name:  fmt.Sprintf("User %d", i),
			Email: fmt.Sprintf("user%d@example.com", i),
		}
		if err := stream.Send(user); err != nil {
			return err
		}
		log.Printf("Sent user-%d", i)
		time.Sleep(500 * time.Millisecond) // Simulate delay
	}
	return nil
}

// Client stream
func (s *userServiceServer) CreateUsers(stream pb.UserService_CreateUsersServer) error {
	count := 0
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&pb.CreateUsersResponse{Count: int32(count)})
		}
		if err != nil {
			return err
		}
		log.Printf("Creating user: %s (%s)", req.Name, req.Email)
		count++
	}
}

// Bidirectional stream
func (s *userServiceServer) ChatWithUsers(stream pb.UserService_ChatWithUsersServer) error {
	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		log.Printf("Received from %s: %s", msg.UserId, msg.Message)

		// Echo back the message
		response := &pb.UserChatMessage{
			UserId:    "server",
			Message:   "Ack: " + msg.Message,
			Timestamp: time.Now().Unix(),
		}

		if err := stream.Send(response); err != nil {
			return err
		}
	}
}

// Test the limitation of number of streams in http2
func (s *userServiceServer) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	log.Printf("Received request: %d", req.GetId())
	time.Sleep(10 * time.Second)
	return nil, nil
}

func LoggingInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	log.Println("Before handling:", info.FullMethod)
	resp, err := handler(ctx, req)
	log.Println("After handling:", info.FullMethod)
	return resp, err
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer(grpc.MaxConcurrentStreams(10), grpc.UnaryInterceptor(LoggingInterceptor))
	pb.RegisterUserServiceServer(grpcServer, &userServiceServer{})

	log.Println("Server listening on :50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

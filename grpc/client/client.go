package main

import (
	"context"
	"io"
	"log"
	"sync"
	"time"

	"google.golang.org/grpc"

	pb "github.com/nethish/playground/grpc-stream/gen/api"
)

func spamGetUser(client pb.UserServiceClient) {
	wg := sync.WaitGroup{}
	for i := range 1000 {
		wg.Add(1)
		go func(i int) {
			a, b := client.GetUser(context.Background(), &pb.GetUserRequest{Id: int64(i)})
			log.Printf("Request %d sent: %v %v\n", i, a, b)
			wg.Done()
		}(i)
		time.Sleep(10 * time.Millisecond)
	}

	wg.Wait()
}

func listUsers(client pb.UserServiceClient) {
	log.Println("Calling ListUsers (Server Streaming)...")
	stream, err := client.ListUsers(context.Background(), &pb.ListUsersRequest{Filter: "all"})
	if err != nil {
		log.Fatalf("ListUsers failed: %v", err)
	}

	for {
		user, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Stream recv failed: %v", err)
		}
		log.Printf("User: %s - %s", user.Name, user.Email)
	}
}

func createUsers(client pb.UserServiceClient) {
	log.Println("Calling CreateUsers (Client Streaming)...")
	stream, err := client.CreateUsers(context.Background())
	if err != nil {
		log.Fatalf("CreateUsers failed: %v", err)
	}

	users := []*pb.CreateUserRequest{
		{Name: "Alice", Email: "alice@example.com"},
		{Name: "Bob", Email: "bob@example.com"},
		{Name: "Charlie", Email: "charlie@example.com"},
	}

	for _, u := range users {
		if err := stream.Send(u); err != nil {
			log.Fatalf("Send failed: %v", err)
		}
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("CloseAndRecv failed: %v", err)
	}

	log.Printf("Created %d users", res.Count)
}

func chatWithUsers(client pb.UserServiceClient) {
	log.Println("Calling ChatWithUsers (Bidirectional Streaming)...")
	stream, err := client.ChatWithUsers(context.Background())
	if err != nil {
		log.Fatalf("ChatWithUsers failed: %v", err)
	}

	done := make(chan struct{})

	// Receiving messages
	go func() {
		for {
			in, err := stream.Recv()
			if err == io.EOF {
				close(done)
				return
			}
			if err != nil {
				log.Fatalf("Recv failed: %v", err)
			}
			log.Printf("Server: %s", in.Message)
		}
	}()

	// Sending messages
	messages := []string{"Hello", "How are you?", "Goodbye"}
	for _, msg := range messages {
		if err := stream.Send(&pb.UserChatMessage{
			UserId:    "client-1",
			Message:   msg,
			Timestamp: time.Now().Unix(),
		}); err != nil {
			log.Fatalf("Send failed: %v", err)
		}
		time.Sleep(500 * time.Millisecond)
	}

	stream.CloseSend()
	<-done
}

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewUserServiceClient(conn)

	listUsers(client)
	createUsers(client)
	chatWithUsers(client)
	spamGetUser(client)
}

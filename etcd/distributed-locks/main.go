package main

import (
	"context"
	"fmt"
	"log"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/concurrency"
)

func main() {
	// Connect to etcd
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"localhost:2379"}, // Update if needed
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		log.Fatalf("Error connecting to etcd: %v", err)
	}
	defer cli.Close()

	// Create a new session for the lock
	session, err := concurrency.NewSession(cli, concurrency.WithTTL(10))
	if err != nil {
		log.Fatalf("Error creating session: %v", err)
	}
	defer session.Close()

	// Create a new mutex under a key prefix
	mutex := concurrency.NewMutex(session, "/my-lock/")

	// Try to acquire the lock
	fmt.Println("Trying to acquire lock...")
	if err := mutex.Lock(context.TODO()); err != nil {
		log.Fatalf("Failed to acquire lock: %v", err)
	}
	fmt.Println("Lock acquired!")

	// Simulate critical section
	time.Sleep(20 * time.Second)

	// Release the lock
	if err := mutex.Unlock(context.TODO()); err != nil {
		log.Fatalf("Failed to release lock: %v", err)
	}
	fmt.Println("Lock released!")
}

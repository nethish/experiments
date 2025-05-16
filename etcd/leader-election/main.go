package main

import (
	"context"
	"crypto/rand"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/concurrency"
)

func main() {
	// Create an etcd client
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"localhost:2379"}, // Replace with your etcd endpoints
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		log.Fatalf("Failed to connect to etcd: %v", err)
	}
	defer cli.Close()

	// Create a new session for election campaign
	session, err := concurrency.NewSession(cli)
	if err != nil {
		log.Fatalf("Failed to create session: %v", err)
	}
	defer session.Close()

	// Create a new election on the given prefix
	election := concurrency.NewElection(session, "/my-election/")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle termination signals to gracefully resign leadership
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigCh
		fmt.Println("Received termination signal, resigning leadership...")
		cancel()
	}()

	// Campaign to become the leader
	hostname, _ := os.Hostname()
	hostname += rand.Text()

	log.Printf("%s is campaigning for leadership...", hostname)
	if err := election.Campaign(ctx, hostname); err != nil {
		log.Fatalf("Failed to campaign for leader: %v", err)
	}

	// Acquired leadership
	log.Printf("%s is the leader", hostname)

	// Keep leadership until context is canceled
	<-ctx.Done()
	log.Printf("%s's leadership relinquished", hostname)
}

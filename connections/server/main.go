package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"sync/atomic"
	"time"
)

const (
	SERVER_HOST = "localhost"
	SERVER_PORT = "8080"
	SERVER_TYPE = "tcp"
)

var (
	RCVD      int64 = 0
	connected int64 = 0
)

func printStats() {
	ticker := time.NewTicker(2 * time.Second)
	for range ticker.C {
		fmt.Printf("Active Connections: %d | RCVD: %d\n", connected, RCVD)
	}
}

func main() {
	go printStats()

	listener, err := net.Listen(SERVER_TYPE, SERVER_HOST+":"+SERVER_PORT)
	if err != nil {
		log.Fatalf("Error listening: %s", err.Error())
	}
	defer listener.Close() // Ensure the listener is closed when main exits

	fmt.Printf("Listening on %s:%s\n", SERVER_HOST, SERVER_PORT)
	fmt.Println("Waiting for clients...")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Error accepting: %s", err.Error())
			continue
		}

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	atomic.AddInt64(&connected, 1)
	defer atomic.AddInt64(&connected, -1)
	defer conn.Close() // Close the connection when the function exits

	reader := bufio.NewReader(conn)

	for {
		message, err := reader.ReadString('\n') // Read until a newline character
		if err != nil {
			// fmt.Printf("Client %s disconnected or error reading: %s\n", conn.RemoteAddr().String(), err.Error())
			return
		}

		atomic.AddInt64(&RCVD, 1)

		_, err = conn.Write([]byte("Server received: " + message))
		if err != nil {
			fmt.Printf("Error writing to client %s: %s\n", conn.RemoteAddr().String(), err.Error())
			return // Exit the goroutine if there's an error writing
		}
	}
}

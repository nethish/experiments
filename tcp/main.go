package main

import (
	"bufio"
	"fmt"
	"net"
	"time"
)

func main() {
	// Start listening on port 8080
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Error starting TCP server:", err)
		return
	}
	defer listener.Close()

	fmt.Println("TCP server listening on port 8080...")

	for {
		// Accept new connection
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}
		fmt.Println("New client connected!")

		// Handle the connection in a new goroutine
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	fmt.Println("Remote Addr", conn.RemoteAddr())
	fmt.Println("Local Addr", conn.LocalAddr())

	reader := bufio.NewReader(conn)
	for {
		// Read data from client
		// NOTE: This has to be set before every read. Otherwise it uses old expired deadline.
		conn.SetReadDeadline(time.Now().Add(5 * time.Second))
		message, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Client disconnected.", err)
			return
		}

		fmt.Printf("Received: %s", message)

		// Echo the message back
		_, err = conn.Write([]byte("Echo: " + message))
		if err != nil {
			fmt.Println("Error writing to client:", err)
			return
		}
	}
}

package main

import (
	"fmt"
	"net"
	"os"
)

func main() {
	socketPath := "/tmp/my_socket"

	os.Remove(socketPath) // remove old socket if exists

	ln, err := net.Listen("unix", socketPath)
	if err != nil {
		panic(err)
	}
	defer ln.Close()

	fmt.Println("Server is listening on", socketPath)

	conn, err := ln.Accept()
	fmt.Println("Remote addr:", conn.RemoteAddr()) // Won't print anything
	fmt.Println("Local addr:", conn.LocalAddr())   // Print's the socket path
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	buf := make([]byte, 1024)
	n, _ := conn.Read(buf)
	fmt.Println("Received:", string(buf[:n]))

	// Write
	conn.Write(buf)
}

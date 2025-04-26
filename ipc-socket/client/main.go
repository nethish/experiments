package main

import (
	"fmt"
	"net"
)

func main() {
	conn, err := net.Dial("unix", "/tmp/my_socket")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	conn.Write([]byte("Hello via UNIX socket!"))

	buf := make([]byte, 1024)
	_, _ = conn.Read(buf[:])
	fmt.Println("Recieved ", string(buf))
}

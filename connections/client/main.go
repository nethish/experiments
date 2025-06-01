package main

import (
	"fmt"
	"net"
	"sync"
	"sync/atomic"
	"time"
)

const (
	SERVER_HOST = "localhost"
	SERVER_PORT = "8080"
	SERVER_TYPE = "tcp"

	TOTAL_CONN = 1030
)

var (
	SENT         int64 = 0
	connected    int64 = 0
	connectedIds       = make(map[int]net.Conn, 0)
)

func printStats() {
	ticker := time.NewTicker(2 * time.Second)
	for range ticker.C {
		fmt.Printf("CONNECTED: %d | SENT: %d \n", connected, SENT)

		// mu.Lock()
		// for i := range TOTAL_CONN {
		// 	if _, ok := connectedIds[i]; !ok {
		// 		fmt.Println("Nope bro", i)
		// 	}
		// 	if conn, ok := connectedIds[i]; ok && conn == nil {
		//
		// 		fmt.Printf("Identified client %d as the first connection not present in the ids\n", i)
		// 		_ = conn
		//
		// 		break
		// 	}
		// }
		// mu.Unlock()
	}
}

func main() {
	go printStats()
	LinearDial()
	// ConcurrentDial()
}

func LinearDial() {
	for i := range TOTAL_CONN {
		conn, err := net.Dial(SERVER_TYPE, SERVER_HOST+":"+SERVER_PORT)
		if err != nil {
			fmt.Printf("Error connecting to server: %s; ClientID: %d\n", err.Error(), i)
			return
		}
		atomic.AddInt64(&connected, 1)
		defer conn.Close()
	}

	time.Sleep(30 * time.Second)
}

var mu = sync.Mutex{}

func ConcurrentDial() {
	wg := sync.WaitGroup{}
	for i := range TOTAL_CONN {
		wg.Add(1)
		go func(clientID int) {
			mu.Lock()
			connectedIds[clientID] = nil
			mu.Unlock()
			conn, err := net.Dial(SERVER_TYPE, SERVER_HOST+":"+SERVER_PORT)
			defer conn.Close()
			if err != nil {
				fmt.Printf("Error connecting to server: %s; ClientID: %d\n", err.Error(), clientID)
				return
			}
			atomic.AddInt64(&connected, 1)
			mu.Lock()
			connectedIds[clientID] = conn
			mu.Unlock()

			_, err = conn.Write([]byte("Hello!\n"))
			if err != nil {
				fmt.Println("Failed to write to server", err)
			} else {
				atomic.AddInt64(&SENT, 1)
			}
			time.Sleep(30 * time.Second)
			wg.Done()
		}(i)
	}

	fmt.Println("Done")
	wg.Wait()
}

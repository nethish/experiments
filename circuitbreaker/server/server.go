package main

import (
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
	"time"
)

var online int32 // 1 means online, 0 means offline

func handler(w http.ResponseWriter, r *http.Request) {
	if atomic.LoadInt32(&online) == 0 {
		http.Error(w, "Server offline", http.StatusServiceUnavailable)
		return
	}
	fmt.Fprintln(w, "Hello, client!")
}

func main() {
	http.HandleFunc("/", handler)

	go func() {
		for {
			atomic.StoreInt32(&online, 1)
			log.Println("Server is ONLINE")
			time.Sleep(5 * time.Second)
			atomic.StoreInt32(&online, 0)
			log.Println("Server is OFFLINE")
			time.Sleep(5 * time.Second)
		}
	}()

	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

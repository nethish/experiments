package main

import (
    "context"
    "fmt"
    "log"
    "time"

    clientv3 "go.etcd.io/etcd/client/v3"
)

func main() {
    // Connect to etcd
    cli, err := clientv3.New(clientv3.Config{
        Endpoints:   []string{"localhost:2379"},
        DialTimeout: 5 * time.Second,
    })
    if err != nil {
        log.Fatalf("Error connecting to etcd: %v", err)
    }
    defer cli.Close()

    ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
    defer cancel()

    // Put key
    _, err = cli.Put(ctx, "foo", "bar")
    if err != nil {
        log.Fatalf("Failed to put key: %v", err)
    }
    fmt.Println("Put key: foo -> bar")

    // Get key
    resp, err := cli.Get(ctx, "foo")
    if err != nil {
        log.Fatalf("Failed to get key: %v", err)
    }

    for _, kv := range resp.Kvs {
        fmt.Printf("Get key: %s -> %s\n", kv.Key, kv.Value)
    }
}

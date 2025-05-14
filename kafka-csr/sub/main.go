package main

import (
    "fmt"
    "github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

func main() {
    c, err := kafka.NewConsumer(&kafka.ConfigMap{
        "bootstrap.servers": "localhost:9092",
        "group.id":          "test-group",
        "auto.offset.reset": "earliest",
    })

    if err != nil {
        panic(err)
    }

    c.SubscribeTopics([]string{"test-topic"}, nil)

    for {
        msg, err := c.ReadMessage(-1)
        if err == nil {
            fmt.Printf("Received: %s\n", string(msg.Value))
        } else {
            fmt.Printf("Error: %v\n", err)
        }
    }
}

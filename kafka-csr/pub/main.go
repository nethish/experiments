package main

import (
        "fmt"

        "github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

func main() {
        p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "localhost:9092"})
        if err != nil {
                panic(err)
        }
        defer p.Close()

        topic := "test-topic"
        for i := 0; i < 10; i++ {
                msg := fmt.Sprintf("Message %d", i)
                p.Produce(&kafka.Message{
                        TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},

                        Value: []byte(msg),
                }, nil)
        }

        // Wait for messages to be delivered
        p.Flush(5000)
}

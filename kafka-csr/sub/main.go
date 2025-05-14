package main

import (
	"encoding/binary"
	"fmt"
	"reflect"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/linkedin/goavro/v2"
	"github.com/riferrei/srclient"
)

func main() {
	topic := "dynamic-user"
	schemaRegistryURL := "http://localhost:8081"

	schemaClient := srclient.CreateSchemaRegistryClient(schemaRegistryURL)

	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost:29092",
		"group.id":          "go-avro-consumer",
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		panic(err)
	}
	defer consumer.Close()

	consumer.SubscribeTopics([]string{topic}, nil)

	for {
		msg, err := consumer.ReadMessage(-1)
		if err != nil {
			fmt.Printf("Consumer error: %v\n", err)
			continue
		}

		value := msg.Value
		if len(value) < 5 || value[0] != 0 {
			fmt.Println("Invalid message format")
			continue
		}

		// Read schema ID
		schemaID := int(binary.BigEndian.Uint32(value[1:5]))
		schema, err := schemaClient.GetSchema(schemaID)
		if err != nil {
			fmt.Printf("Failed to fetch schema: %v\n", err)
			continue
		}

		codec, err := goavro.NewCodec(schema.Schema())
		if err != nil {
			fmt.Printf("Failed to create codec: %v\n", err)
			continue
		}

		// Decode
		native, _, err := codec.NativeFromBinary(value[5:])
		fmt.Println(reflect.TypeOf(native))
		if err != nil {
			fmt.Printf("Decode error: %v\n", err)
			continue
		}

		fmt.Printf("Received message: %+v\n", native)
	}
}

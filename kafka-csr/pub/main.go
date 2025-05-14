package main

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/linkedin/goavro/v2"
	"github.com/riferrei/srclient"
)

func main() {
	topic := "dynamic-user"
	schemaRegistryURL := "http://localhost:8081"

	// Avro schema as string
	avroSchema := `{
		"type": "record",
		"name": "User",
		"namespace": "example",
		"fields": [
			{"name": "name", "type": "string"},
			{"name": "age", "type": "int"}
		]
	}`

	// Connect to Schema Registry
	schemaClient := srclient.CreateSchemaRegistryClient(schemaRegistryURL)

	// Register schema dynamically
	subject := topic + "-value"
	schema, err := schemaClient.CreateSchema(subject, avroSchema, srclient.Avro)
	if err != nil {
		panic(err)
	}
	fmt.Println("Schema registered with ID:", schema.ID())

	// Encode data
	codec, err := goavro.NewCodec(avroSchema)
	if err != nil {
		panic(err)
	}

	native := map[string]interface{}{
		"name": "Alice",
		"age":  30,
	}

	binaryData, err := codec.BinaryFromNative(nil, native)
	if err != nil {
		panic(err)
	}

	// Message format: Magic byte (0) + schema ID (4 bytes) + Avro data
	var buf bytes.Buffer
	buf.WriteByte(0)
	_ = binary.Write(&buf, binary.BigEndian, int32(schema.ID()))
	buf.Write(binaryData)

	fmt.Println("Initializing producer")
	// Kafka producer
	producer, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "localhost:29092"})
	if err != nil {
		panic(err)
	}
	defer producer.Close()

	err = producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          buf.Bytes(),
	}, nil)
	if err != nil {
		panic(err)
	}

	producer.Flush(5000)
	fmt.Println("Produced message")
}

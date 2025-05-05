package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodbstreams"
	streamtypes "github.com/aws/aws-sdk-go-v2/service/dynamodbstreams/types"
)

// Configuration parameters
const (
	tableName         = "nethish-stream-test" // Replace with your table name
	pollingInterval   = 1 * time.Second       // How often to poll for new records
	shardIteratorType = "LATEST"              // Use "TRIM_HORIZON" for reading from the start
)

func main() {
	// Initialize logger
	logger := log.New(os.Stdout, "dynamodb-streams: ", log.LstdFlags)
	logger.Println("Starting DynamoDB Streams consumer")

	// Load AWS configuration
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion("us-east-1"), // Change to your AWS region
	)
	if err != nil {
		logger.Fatalf("Failed to load AWS configuration: %v", err)
	}

	// Create DynamoDB client
	dynamoClient := dynamodb.NewFromConfig(cfg)

	// Create DynamoDB Streams client
	streamsClient := dynamodbstreams.NewFromConfig(cfg)

	// Get the stream ARN for the table
	tableInfo, err := dynamoClient.DescribeTable(context.TODO(), &dynamodb.DescribeTableInput{
		TableName: aws.String(tableName),
	})
	if err != nil {
		logger.Fatalf("Failed to describe table: %v", err)
	}

	if tableInfo.Table.StreamSpecification == nil || !*tableInfo.Table.StreamSpecification.StreamEnabled {
		logger.Fatalf("Streams are not enabled for table %s", tableName)
	}

	streamArn := tableInfo.Table.LatestStreamArn
	logger.Printf("Found stream ARN: %s", *streamArn)

	// Get stream description
	streamInfo, err := streamsClient.DescribeStream(context.TODO(), &dynamodbstreams.DescribeStreamInput{
		StreamArn: streamArn,
	})
	if err != nil {
		logger.Fatalf("Failed to describe stream: %v", err)
	}

	// Set up graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle termination signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigChan
		logger.Println("Shutdown signal received")
		cancel()
	}()

	// Process each shard in the stream
	for _, shard := range streamInfo.StreamDescription.Shards {
		go processShardRecords(ctx, streamsClient, streamArn, shard.ShardId, logger)
	}

	// Wait for shutdown
	<-ctx.Done()
	logger.Println("Shutting down...")
}

func processShardRecords(ctx context.Context, client *dynamodbstreams.Client, streamArn *string, shardID *string, logger *log.Logger) {
	logger.Printf("Starting to process shard: %s", *shardID)

	// Get shard iterator
	iteratorOutput, err := client.GetShardIterator(ctx, &dynamodbstreams.GetShardIteratorInput{
		StreamArn:         streamArn,
		ShardId:           shardID,
		ShardIteratorType: streamtypes.ShardIteratorType(shardIteratorType),
	})
	if err != nil {
		logger.Printf("Failed to get shard iterator for shard %s: %v", *shardID, err)
		return
	}

	iterator := iteratorOutput.ShardIterator

	// Poll for records
	for iterator != nil && ctx.Err() == nil {
		// Get records using the iterator
		recordsOutput, err := client.GetRecords(ctx, &dynamodbstreams.GetRecordsInput{
			ShardIterator: iterator,
			Limit:         aws.Int32(100), // Adjust based on your needs
		})
		if err != nil {
			logger.Printf("Error getting records: %v", err)
			time.Sleep(pollingInterval)
			continue
		}

		// Process records
		for _, record := range recordsOutput.Records {
			processRecord(record, logger)
		}

		// printt(recordsOutput.ResultMetadata)

		// Update iterator for next poll
		iterator = recordsOutput.NextShardIterator

		// If no records were returned, wait before polling again
		if len(recordsOutput.Records) == 0 {
			time.Sleep(pollingInterval)
		}
	}

	if ctx.Err() != nil {
		logger.Printf("Stopping processing of shard %s due to context cancellation", *shardID)
	} else {
		logger.Printf("Shard %s has been closed", *shardID)
	}
}

func printt(a any) {
	jsonBytes, _ := json.MarshalIndent(a, "", "  ")
	fmt.Println(string(jsonBytes))
}

func processRecord(record streamtypes.Record, logger *log.Logger) {
	// Log the event type
	logger.Printf("Received event: %s", record.EventName)
	printt(record)

	// Process based on event type
	switch record.EventName {
	case "INSERT":
		handleInsert(record, logger)
	case "MODIFY":
		handleModify(record, logger)
	case "REMOVE":
		handleRemove(record, logger)
	}
}

func handleInsert(record streamtypes.Record, logger *log.Logger) {
	if record.Dynamodb.NewImage != nil {
		logger.Printf("New item inserted: %v", record.Dynamodb.NewImage)
		// Add your business logic for handling inserts here
	}
}

func handleModify(record streamtypes.Record, logger *log.Logger) {
	if record.Dynamodb.NewImage != nil && record.Dynamodb.OldImage != nil {
		logger.Printf("Item modified from %v to %v", record.Dynamodb.OldImage, record.Dynamodb.NewImage)
		// Add your business logic for handling modifications here
	}
}

func handleRemove(record streamtypes.Record, logger *log.Logger) {
	if record.Dynamodb.OldImage != nil {
		logger.Printf("Item removed: %v", record.Dynamodb.OldImage)
		// Add your business logic for handling removals here
	}
}

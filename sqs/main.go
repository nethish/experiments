package main

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

func main() {
	ctx := context.Background()

	// Load AWS config (uses environment, ~/.aws/config, IAM role, etc.)
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	// Create SQS client
	client := sqs.NewFromConfig(cfg)

	queueURL := "<SQS Queue>"

	go rescueIndefinitely()
	// Receive messages indefinitely but don't delete
	for {
		// Receive a message
		output, err := client.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{
			QueueUrl:            aws.String(queueURL),
			MaxNumberOfMessages: 1,  // Read one message at a time
			WaitTimeSeconds:     5,  // Long poll for 5 seconds
			VisibilityTimeout:   30, // Hide the message for 30 sec after reading
		})
		if err != nil {
			log.Fatalf("failed to receive message, %v", err)
		}

		if len(output.Messages) == 0 {
			fmt.Println("No messages received")
		}

		for _, message := range output.Messages {
			fmt.Printf("Received message ID: %s\n", aws.ToString(message.MessageId))
			fmt.Printf("Body: %s\n", aws.ToString(message.Body))

			// NOTE: You should delete the message after processing
			// (Optional) Delete the message after processing
			// _, err := client.DeleteMessage(ctx, &sqs.DeleteMessageInput{
			// 	QueueUrl:      aws.String(queueURL),
			// 	ReceiptHandle: message.ReceiptHandle,
			// })
			// if err != nil {
			// 	log.Printf("failed to delete message, %v", err)
			// } else {
			// 	fmt.Println("Message deleted successfully")
			// }
		}
	}
}

func rescueIndefinitely() {
	rescue()
}

func rescue() {
	prefix := "[Rescuer] "
	ctx := context.Background()

	// Load AWS config (uses environment, ~/.aws/config, IAM role, etc.)
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	// Create SQS client
	client := sqs.NewFromConfig(cfg)

	queueURL := "<DLQ>"

	// Receive messages indefinitely but don't delete
	for {
		// Receive a message
		output, err := client.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{
			QueueUrl:            aws.String(queueURL),
			MaxNumberOfMessages: 1,  // Read one message at a time
			WaitTimeSeconds:     5,  // Long poll for 5 seconds
			VisibilityTimeout:   30, // Hide the message for 30 sec after reading
		})
		if err != nil {
			log.Fatalf(prefix+"failed to receive message, %v", err)
		}

		if len(output.Messages) == 0 {
			fmt.Println(prefix + "No messages received")
		}

		for _, message := range output.Messages {
			fmt.Printf(prefix+"Received message ID: %s\n", aws.ToString(message.MessageId))
			fmt.Printf(prefix+"Body: %s\n", aws.ToString(message.Body))

		}
	}
}

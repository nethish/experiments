package main

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

func main() {
	ctx := context.Background()

	// Load AWS config
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	client := sqs.NewFromConfig(cfg)

	queueUrl := "https://sqs.us-east-1.amazonaws.com/<account>/nethish.fifo"

	var wg sync.WaitGroup
	wg.Add(2) // one for producer, one for consumer

	// Producer Goroutine
	go func() {
		defer wg.Done()
		// produceMessages(ctx, client, queueUrl)
	}()

	// Consumer Goroutine
	go func() {
		defer wg.Done()
		consumeMessages(ctx, client, queueUrl)
	}()

	wg.Wait()
}

func produceMessages(ctx context.Context, client *sqs.Client, queueUrl string) {
	for i := 0; i < 5; i++ {
		body := fmt.Sprintf("Message number %d", i+1)

		_, err := client.SendMessage(ctx, &sqs.SendMessageInput{
			QueueUrl:               aws.String(queueUrl),
			MessageBody:            aws.String(body),
			MessageGroupId:         aws.String("group-1"),
			MessageDeduplicationId: aws.String(fmt.Sprintf("%d", time.Now().UnixNano())),
		})
		if err != nil {
			log.Printf("failed to send message: %v", err)
		} else {
			log.Printf("Produced: %s", body)
		}

		time.Sleep(2 * time.Second) // small delay between messages
	}
}

func consumeMessages(ctx context.Context, client *sqs.Client, queueUrl string) {
	for {
		output, err := client.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{
			QueueUrl:            aws.String(queueUrl),
			MaxNumberOfMessages: 2,
			WaitTimeSeconds:     5, // long polling
		})
		if err != nil {
			log.Printf("failed to receive messages: %v", err)
			continue
		}

		if len(output.Messages) == 0 {
			log.Println("No messages received. Waiting...")
			time.Sleep(3 * time.Second)
			continue
		}

		for _, message := range output.Messages {
			log.Printf("Consumed: %s; MessageId: %s", aws.ToString(message.Body), aws.ToString(message.MessageId))

			// Delete the message after processing
			// _, err := client.DeleteMessage(ctx, &sqs.DeleteMessageInput{
			// 	QueueUrl:      aws.String(queueUrl),
			// 	ReceiptHandle: message.ReceiptHandle,
			// })
			// if err != nil {
			// 	log.Printf("failed to delete message: %v", err)
			// } else {
			// 	log.Printf("Deleted message ID: %s", aws.ToString(message.MessageId))
			// }
		}
	}
}

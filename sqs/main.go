package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

func main() {
	test()
	ctx := context.Background()

	// Load AWS config
	cfg, err := config.LoadDefaultConfig(ctx)
	cfg.Region = "ap-south-1"
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	client := sqs.NewFromConfig(cfg)

	queueUrl := "https://sqs.ap-south-1.amazonaws.com/<account>/queue"

	var wg sync.WaitGroup
	wg.Add(2) // one for producer, one for consumer

	// Producer Goroutine
	go func() {
		defer wg.Done()
		produceMessages(ctx, client, queueUrl)
	}()

	// Consumer Goroutine
	go func() {
		defer wg.Done()
		consumeMessages(ctx, client, queueUrl)
	}()

	wg.Wait()
}

func test() {
	message := Message{
		Id:   1,
		Time: time.Now(),
	}
	jsonBytes, _ := json.Marshal(message)

	var mess *Message
	err := json.Unmarshal(jsonBytes, mess)
	if err != nil {
		fmt.Printf("error unmarshalling: %w\n", err)
	}

	fmt.Printf("success unmarshalling: %s\n", string(jsonBytes))
}

type Message struct {
	Id   int       `json:"id"`
	Time time.Time `json:"time"`
}

func produceMessages(ctx context.Context, client *sqs.Client, queueUrl string) {
	for i := 0; i < 100; i++ {
		message := Message{
			Id:   i,
			Time: time.Now(),
		}

		jsonBytes, _ := json.Marshal(message)

		_, err := client.SendMessage(ctx, &sqs.SendMessageInput{
			QueueUrl:    aws.String(queueUrl),
			MessageBody: aws.String(string(jsonBytes)),
		})
		if err != nil {
			log.Printf("failed to send message: %v", err)
		} else {
			log.Printf("Produced: %s", string(jsonBytes))
		}

		time.Sleep(1 * time.Second) // small delay between messages
	}
}

func consumeMessages(ctx context.Context, client *sqs.Client, queueUrl string) {
	for {
		output, err := client.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{
			QueueUrl:            aws.String(queueUrl),
			MaxNumberOfMessages: 10,
			WaitTimeSeconds:     10, // long polling
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
			var mess Message
			err = json.Unmarshal([]byte(aws.ToString(message.Body)), &mess)
			if err != nil {
				fmt.Printf("error unmarshalling message: %s; Message: %s\n", err.Error(), aws.ToString(message.Body))
				continue
			}

			fmt.Printf("Message Id: %s; Receive time: %v\n", mess.Id, time.Since(mess.Time))

			// Delete the message after processing
			_, err := client.DeleteMessage(ctx, &sqs.DeleteMessageInput{
				QueueUrl:      aws.String(queueUrl),
				ReceiptHandle: message.ReceiptHandle,
			})
			if err != nil {
				log.Printf("failed to delete message: %v", err)
			} else {
				log.Printf("Deleted message ID: %s", aws.ToString(message.MessageId))
			}
		}
	}
}

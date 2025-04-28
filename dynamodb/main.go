package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/google/uuid"
)

func main() {
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithEndpointResolver(aws.EndpointResolverFunc(
			func(service, region string) (aws.Endpoint, error) {
				return aws.Endpoint{URL: "http://localhost:8000"}, nil
			},
		)),
	)
	if err != nil {
		panic(err)
	}

	client := dynamodb.NewFromConfig(cfg)

	for range 100000 {
		_, _ = client.PutItem(context.TODO(), &dynamodb.PutItemInput{
			Item: map[string]types.AttributeValue{
				"ID": &types.AttributeValueMemberS{Value: string(uuid.New().String())},
			},
			TableName: aws.String("People"),
		})
		// fmt.Println(o, err)
	}

	for range 1 {
		out, err := client.GetItem(context.TODO(), &dynamodb.GetItemInput{
			TableName: aws.String("People"),
			Key: map[string]types.AttributeValue{
				"ID": &types.AttributeValueMemberS{Value: "Fuc"},
			},
		})
		if err != nil {
			panic(err)
		}

		fmt.Printf("Item %v:\n", out.Item)
	}
}

package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/sqs"
)

type MyEvent struct {
	ID      string `json:"id"`
	Message string `json:"message"`
}

func HandleRequest(ctx context.Context, event MyEvent) (string, error) {
	sess := session.Must(session.NewSession(&aws.Config{
		Region:   aws.String("us-east-1"),
		Endpoint: aws.String("http://localhost:4566"),
	}))

	// DynamoDB client
	svc := dynamodb.New(sess)
	item := map[string]*dynamodb.AttributeValue{
		"id": {
			S: aws.String(event.ID),
		},
		"message": {
			S: aws.String(event.Message),
		},
	}

	_, err := svc.PutItem(&dynamodb.PutItemInput{
		TableName: aws.String("my-table"),
		Item:      item,
	})

	if err != nil {
		log.Fatalf("Got error calling PutItem: %s", err)
		return "", err
	}

	// SQS client
	sqsSvc := sqs.New(sess)
	queueURL := os.Getenv("QUEUE_URL")

	msgResult, err := sqsSvc.ReceiveMessage(&sqs.ReceiveMessageInput{
		QueueUrl:            aws.String(queueURL),
		MaxNumberOfMessages: aws.Int64(1),
	})

	if err != nil {
		log.Fatalf("Got error calling ReceiveMessage: %s", err)
		return "", err
	}

	if len(msgResult.Messages) > 0 {
		for _, message := range msgResult.Messages {
			fmt.Printf("Message: %s\n", *message.Body)
		}
	}

	return fmt.Sprintf("Successfully processed %s", event.ID), nil
}

func main() {
	lambda.Start(HandleRequest)
}

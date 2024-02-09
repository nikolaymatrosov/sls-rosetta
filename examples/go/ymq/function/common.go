package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

func sendMessageToQueue(
	ctx context.Context,
	ymqName string,
	message string,
	origin string,
	delay int32,
) (*sqs.SendMessageOutput, error) {
	// Create a custom endpoint resolver to resolve Yandex Message Queue endpoint
	customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		return aws.Endpoint{
			URL:           "https://message-queue.api.cloud.yandex.net",
			SigningRegion: "ru-central1",
		}, nil
	})
	// Load the SDK's configuration from environment and shared config
	// In the serverless environment, the configuration is loaded from the environment variables
	// AWS_ACCESS_KEY_ID and AWS_SECRET_ACCESS_KEY.
	cfg, err := config.LoadDefaultConfig(
		ctx,
		config.WithEndpointResolverWithOptions(customResolver),
	)
	if err != nil {
		panic("configuration error, " + err.Error())
	}
	// Create an Amazon SQS service client
	client := sqs.NewFromConfig(cfg)

	gQInput := &sqs.GetQueueUrlInput{
		QueueName: &ymqName,
	}

	urlRes, err := client.GetQueueUrl(ctx, gQInput)

	if err != nil {
		fmt.Printf("Got an error getting the queue URL: %s", err)
		return nil, err
	}

	sMInput := &sqs.SendMessageInput{
		MessageAttributes: map[string]types.MessageAttributeValue{
			"Origin": {
				DataType:    aws.String("String"),
				StringValue: aws.String(origin),
			},
		},
		MessageBody:  aws.String(message),
		QueueUrl:     urlRes.QueueUrl,
		DelaySeconds: delay,
	}
	resp, err := client.SendMessage(ctx, sMInput)
	return resp, nil
}

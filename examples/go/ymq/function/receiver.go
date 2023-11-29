package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// Структура запроса триггера YMQ

type EventMetadata struct {
	EventId   string    `json:"event_id"`
	EventType string    `json:"event_type"`
	CreatedAt time.Time `json:"created_at"`
	CloudId   string    `json:"cloud_id"`
	FolderId  string    `json:"folder_id"`
}

type MessageAttributeValue struct {
	DataType    string `json:"data_type"`
	BinaryValue []byte `json:"binary_value"`
	StringValue string `json:"string_value"`
}

type YMQMessageDetails struct {
	QueueId string `json:"queue_id"`
	Message struct {
		MessageId              string
		Md5OfBody              string
		Body                   string
		Attributes             map[string]string
		MessageAttributes      map[string]*MessageAttributeValue
		Md5OfMessageAttributes string
	} `json:"message"`
}

type YMQMessage struct {
	EventMetadata EventMetadata     `json:"event_metadata"`
	Details       YMQMessageDetails `json:"details"`
}

type YMQRequest struct {
	Messages []YMQMessage `json:"messages"`
}

type YMQResponse struct {
	StatusCode int
}

type Request struct {
	Name string `json:"name"`
}

func Receiver(ctx context.Context, event *YMQRequest) (*YMQResponse, error) {
	ymqName := os.Getenv("YMQ_NAME")

	var req Request
	for _, message := range event.Messages {
		if err := json.Unmarshal([]byte(message.Details.Message.Body), &req); err != nil {
			return nil, fmt.Errorf("an error has occurred when parsing body: %v", err)
		}

		fmt.Printf("%+v\n", req)
		req.Name = "test"
		resp := fmt.Sprintf(`{"result": "success", "name": "%s"}`, req.Name)
		_, _ = sendMessageToQueue(ctx, ymqName, resp, "From Receiver Function")

	}
	return &YMQResponse{
		StatusCode: 200,
	}, nil
}

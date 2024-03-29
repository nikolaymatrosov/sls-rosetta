package main

import (
	"context"
	"fmt"
	"os"
)

type RequestContext struct {
	Identity struct {
		SourceIp  string `json:"sourceIp"`
		UserAgent string `json:"userAgent"`
	}
	HttpMethod       string      `json:"httpMethod"`
	RequestId        string      `json:"requestId"`
	RequestTime      string      `json:"requestTime"`
	RequestTimeEpoch int         `json:"requestTimeEpoch"`
	Authorizer       interface{} `json:"authorizer"`
	ApiGateway       struct {
		OperationContext interface{} `json:"operationContext"`
	} `json:"apiGateway"`
}

type HttpEvent struct {
	HttpMethod                      string              `json:"httpMethod"`
	Headers                         map[string]string   `json:"headers"`
	MultiValueHeaders               map[string][]string `json:"multiValueHeaders"`
	QueryStringParameters           map[string]string   `json:"queryStringParameters"`
	MultiValueQueryStringParameters map[string][]string `json:"multiValueQueryStringParameters"`
	RequestContext                  RequestContext      `json:"requestContext"`
	Body                            string              `json:"body"`
	IsBase64Encoded                 bool                `json:"isBase64Encoded"`
}

type HttpResult struct {
	StatusCode        int               `json:"statusCode"`
	Headers           map[string]string `json:"headers"`
	MultiValueHeaders map[string]string `json:"multiValueHeaders"`
	Body              string            `json:"body"`
	IsBase64Encoded   bool              `json:"isBase64Encoded"`
}

//goland:noinspection GoUnusedExportedFunction,GoUnusedParameter
func Sender(ctx context.Context, event *HttpEvent) (*HttpResult, error) {
	ymqName := os.Getenv("YMQ_NAME")

	resp, err := sendMessageToQueue(
		ctx,
		ymqName,
		`{"name":"test"}`,
		"From Sender Function",
		30,
	)
	if err != nil {
		fmt.Println("Got an error sending the message:")
		fmt.Println(err)
		return &HttpResult{
			StatusCode: 500,
			Body:       "Got an error sending the message: " + err.Error(),
		}, nil
	}

	fmt.Println("Sent message with ID: " + *resp.MessageId)
	return &HttpResult{
		StatusCode: 200,
		Body:       "Sent message with ID: " + *resp.MessageId,
	}, nil
}

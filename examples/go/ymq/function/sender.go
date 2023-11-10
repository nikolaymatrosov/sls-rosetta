package main

import (
	"context"
	"fmt"
	"os"
)

/*
type HttpMethod =
    | "OPTIONS"
    | "HEAD"
    | "GET"
    | "POST"
    | "PUT"
    | "PATCH"
    | "DELETE";

export namespace Http {
    export interface Result {
        statusCode: number;
        headers?: Record<string, string>;
        multiValueHeaders?: Record<string, string[]>;
        body?: string;
        isBase64Encoded?: boolean;
    }

    export interface Event {
        httpMethod: HttpMethod;
        headers: Record<string, string>;
        multiValueHeaders: Record<string, string[]>;
        queryStringParameters: Record<string, string>;
        multiValueQueryStringParameters: Record<string, string[]>;
        requestContext: RequestContext & {
            authorizer?: unknown; // TODO: describe type
            apiGateway?: {
                operationContext: unknown;
            };
        };
        body: string;
        isBase64Encoded: boolean;
    }

    export type RequestContext = {
        identity: {
            sourceIp: string;
            userAgent: string;
        };
        httpMethod: HttpMethod;
        requestId: string;
        requestTime: string;
        requestTimeEpoch: number;
    };
}
*/

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

func Sender(ctx context.Context, event *HttpEvent) (*HttpResult, error) {
	ymqName := os.Getenv("YMQ_NAME")

	resp, err := sendMessageToQueue(ctx, ymqName, `{"name":"test"}`, "From Sender Function")
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

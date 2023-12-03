package main

import (
	"context"
	"encoding/json"
	"fmt"
)

// APIGatewayRequest is a struct that represents the structure of an API Gateway v1 request.
type APIGatewayRequest struct {
	// OperationID is a unique identifier for the operation defined in the OpenAPI specification.
	OperationID string `json:"operationId"`

	// Resource is the path to the resource that is being accessed.
	Resource string `json:"resource"`

	// HTTPMethod is the HTTP method used for the request (e.g., GET, POST).
	HTTPMethod string `json:"httpMethod"`

	// Path is the path of the request.
	Path string `json:"path"`

	// PathParameters are the parameters in the path of the request.
	PathParameters map[string]string `json:"pathParameters"`

	// Headers are the headers included in the request.
	Headers map[string]string `json:"headers"`
	// MultiValueHeaders are the headers that can have multiple values represented as a string array.
	// For example, "X-Value: a" and "X-Value: b" are two values for the same header.
	// The Headers property on the key "X-Value" would have only "b", while the
	// MultiValueHeaders property on the key "X-Value" would have array ["a", "b"].
	// This may be not true for all headers, some headers like "Accept" or "Cookie" are treated differently.
	// For these headers values are joined with comma of semicolon.
	MultiValueHeaders map[string][]string `json:"multiValueHeaders"`

	// QueryStringParameters are the parameters in the query string of the request.
	QueryStringParameters map[string]string `json:"queryStringParameters"`
	// MultiValueQueryStringParameters are the query string parameters that can have multiple values.
	// For example, "q=1&q=2". The QueryStringParameters property on the key "q" would have only "2",
	// while the MultiValueQueryStringParameters property on the key "q" would have array ["1", "2"].
	MultiValueQueryStringParameters map[string][]string `json:"multiValueQueryStringParameters"`

	// Parameters are the parameters in the request.
	Parameters map[string]string `json:"parameters"`
	// MultiValueParameters are the parameters that can have multiple values.
	MultiValueParameters map[string][]string `json:"multiValueParameters"`

	// Body is the body of the request.
	Body string `json:"body"`
	// IsBase64Encoded indicates whether the body is Base64 encoded.
	IsBase64Encoded bool `json:"isBase64Encoded,omitempty"`

	// RequestContext is the context of the request.
	RequestContext interface{} `json:"requestContext"`
}

// APIGatewayResponse is a struct that represents the structure of an API Gateway v1 response.
type APIGatewayResponse struct {
	// StatusCode is the HTTP status code of the response.
	StatusCode int `json:"statusCode"`
	// Headers are the headers included in the response.
	Headers map[string]string `json:"headers"`
	// MultiValueHeaders are the headers that can have multiple values.
	MultiValueHeaders map[string][]string `json:"multiValueHeaders"`
	// Body is the body of the response.
	Body string `json:"body"`
	// IsBase64Encoded indicates whether the body is Base64 encoded.
	IsBase64Encoded bool `json:"isBase64Encoded,omitempty"`
}

// Request is a struct that represents the structure of a request.
type Request struct {
	// Name is the name of the user making the request.
	Name string `json:"name"`
}

// Response is a struct that represents the structure of a response.
type Response struct {
	// Message is the response message.
	Message string `json:"message"`
}

// Handler is a function that handles API Gateway requests and responses.
func Handler(ctx context.Context, event *APIGatewayRequest) (*APIGatewayResponse, error) {
	req := &Request{}

	// The Body field of the event request is converted into a Request object to get the passed name.
	if err := json.Unmarshal([]byte(event.Body), &req); err != nil {
		// If there is an error in parsing the body, an error is returned.
		return nil, fmt.Errorf("an error has occurred when parsing body: %v", err)
	}

	// The request is printed to the log.
	fmt.Printf("%+v\n", req)

	// Preparing the response.
	response, err := json.Marshal(Response{
		Message: fmt.Sprintf("Hello, %s!", req.Name),
	})
	if err != nil {
		// If there is an error in marshaling the response, an error is returned.
		return nil, fmt.Errorf("an error has occurred when marshaling response: %v", err)
	}

	return &APIGatewayResponse{
		StatusCode:      200,
		Headers:         map[string]string{"content-type": "application/json"},
		Body:            string(response),
		IsBase64Encoded: false,
	}, nil
}

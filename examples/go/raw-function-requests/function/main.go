package main

import (
	"context"
	"encoding/json"
	"fmt"
)

// Request input JSON document will be automatically converted to the object of this type
type Request struct {
	Message string `json:"message"`
	Number  int    `json:"number"`
}

type ResponseBody struct {
	Request interface{} `sjson:"request"`
}

func Handler(ctx context.Context, request *Request) ([]byte, error) {
	// In function logs, the values of the call context and the request body will be printed
	fmt.Println("context", ctx)
	fmt.Println("request", request)

	// The object containing the response body is converted to an array of bytes
	body, err := json.Marshal(&ResponseBody{
		Request: request,
	})

	if err != nil {
		return nil, err
	}

	// The response body must be returned as an array of bytes
	return body, nil
}

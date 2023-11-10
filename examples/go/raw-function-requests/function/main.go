package main

import (
	"context"
	"encoding/json"
	"fmt"
)

// Input JSON document will be automatically converted to the object of this type
type Request struct {
	Message string `json:"message"`
	Number  int    `json:"number"`
}

type ResponseBody struct {
	Context context.Context `json:"context"`
	Request interface{}     `json:"request"`
}

func Handler(ctx context.Context, request *Request) ([]byte, error) {
	// In function logs, the values of the call context and the request body will be printed
	fmt.Println("context", ctx)
	fmt.Println("request", request)

	// The object containing the response body is converted to an array of bytes
	body, err := json.Marshal(&ResponseBody{
		Context: ctx,
		Request: request,
	})

	if err != nil {
		return nil, err
	}

	// The response body must be returned as an array of bytes
	return body, nil
}

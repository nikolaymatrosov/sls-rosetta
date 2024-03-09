package main

import (
	"context"
	"encoding/json"
	"fmt"
)

type Req struct {
	Name string `json:"name"`
}

func Handler(_ context.Context, req Req) ([]byte, error) {
	// get body
	fmt.Printf("Body: %+v\n", req)

	// return response
	resp := map[string]interface{}{
		"result": "success",
		"name":   req.Name,
	}

	respBytes, err := json.Marshal(resp)
	if err != nil {
		panic(err)
	}

	return respBytes, nil

}

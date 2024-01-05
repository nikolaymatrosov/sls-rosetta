package main

import (
	"context"
	"encoding/json"
	"fmt"
)

type Req struct {
	Name string `json:"name"`
}

func Handler(ctx context.Context, req Req) ([]byte, error) {
	// get body
	fmt.Printf("Body: %+v\n", req)

	// unmarshal body
	//var body map[string]interface{}
	//
	//err := json.Unmarshal(req, &body)

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

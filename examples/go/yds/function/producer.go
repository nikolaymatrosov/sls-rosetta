package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/ydb-platform/ydb-go-sdk/v3"
	"github.com/ydb-platform/ydb-go-sdk/v3/topic/topicwriter"
	yc "github.com/ydb-platform/ydb-go-yc"
)

// ProducerHandler handles HTTP requests to write data to YDS stream
func ProducerHandler(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	// Check required environment variables
	ydbEndpoint := os.Getenv("YDB_ENDPOINT")
	topicName := os.Getenv("YDS_TOPIC_ID")
	if ydbEndpoint == "" || topicName == "" {
		http.Error(w, "YDB_ENDPOINT and YDS_TOPIC_ID environment variables must be set", http.StatusInternalServerError)
		return
	}

	// Parse the request
	var req ProducerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Error decoding request: %v", err)
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	// Validate request
	if req.Message == "" {
		http.Error(w, "Message is required", http.StatusBadRequest)
		return
	}

	// Create event data
	eventData := map[string]interface{}{
		"message":   req.Message,
		"user_id":   req.UserID,
		"action":    req.Action,
		"timestamp": time.Now().Unix(),
	}

	// Convert to JSON
	jsonData, err := json.Marshal(eventData)
	if err != nil {
		log.Printf("Error marshaling event data: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Write to YDS topic
	err = writeToTopic(ctx, ydbEndpoint, topicName, string(jsonData))
	if err != nil {
		log.Printf("Error writing to topic: %v", err)
		http.Error(w, "Failed to write to topic", http.StatusInternalServerError)
		return
	}

	// Return success response
	response := ProducerResponse{
		StatusCode: 200,
		Message:    "Data written to topic successfully",
		StreamID:   topicName,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// writeToTopic writes data to the YDB topic using the YDB Go SDK
func writeToTopic(ctx context.Context, ydbEndpoint, topicName, data string) error {
	db, err := ydb.Open(ctx, ydbEndpoint,
		yc.WithMetadataCredentials(),
	)
	if err != nil {
		return fmt.Errorf("failed to connect to YDB: %w", err)
	}
	defer db.Close(ctx)

	writer, err := db.Topic().StartWriter(topicName)
	if err != nil {
		return fmt.Errorf("failed to create topic writer: %w", err)
	}
	defer writer.Close(ctx)

	msg := topicwriter.Message{Data: bytes.NewReader([]byte(data))}
	if err := writer.Write(ctx, msg); err != nil {
		return fmt.Errorf("failed to write message to topic: %w", err)
	}

	return nil
}

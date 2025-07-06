package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"
)

// ConsumerHandler handles YDS trigger events
func ConsumerHandler(ctx context.Context, event *YDSEvent) (*YDSResponse, error) {
	log.Printf("Received YDS event with %d messages", len(event.Messages))

	// Process each message in the batch
	for i, message := range event.Messages {
		log.Printf("Processing message %d: %s", i+1, message.Details.Data)

		// Parse the message data
		var eventData map[string]interface{}
		if err := json.Unmarshal([]byte(message.Details.Data), &eventData); err != nil {
			log.Printf("Error parsing message data: %v", err)
			continue
		}

		// Process the event
		err := processEvent(ctx, eventData, message.EventMetadata)
		if err != nil {
			log.Printf("Error processing event: %v", err)
			continue
		}

		log.Printf("Successfully processed message %d", i+1)
	}

	return &YDSResponse{
		StatusCode: 200,
		Message:    fmt.Sprintf("Processed %d messages successfully", len(event.Messages)),
	}, nil
}

// processEvent processes a single event from the stream
func processEvent(ctx context.Context, eventData map[string]interface{}, metadata EventMetadata) error {
	// Extract event information
	message, _ := eventData["message"].(string)
	userID, _ := eventData["user_id"].(string)
	action, _ := eventData["action"].(string)
	timestamp, _ := eventData["timestamp"].(float64)

	log.Printf("Processing event - User: %s, Action: %s, Message: %s", userID, action, message)

	// Simulate some processing logic
	switch action {
	case "login":
		log.Printf("User %s logged in at %v", userID, time.Unix(int64(timestamp), 0))
	case "logout":
		log.Printf("User %s logged out at %v", userID, time.Unix(int64(timestamp), 0))
	case "purchase":
		log.Printf("User %s made a purchase: %s", userID, message)
	case "view":
		log.Printf("User %s viewed: %s", userID, message)
	default:
		log.Printf("Unknown action '%s' from user %s", action, userID)
	}

	// Simulate processing time
	time.Sleep(50 * time.Millisecond)

	// In a real implementation, you might:
	// - Store data in a database
	// - Send notifications
	// - Update analytics
	// - Trigger other workflows
	// - Send data to external systems

	return nil
}

// BatchProcessor processes multiple events in a batch
func BatchProcessor(ctx context.Context, events []map[string]interface{}) error {
	log.Printf("Processing batch of %d events", len(events))

	// Group events by user for batch processing
	userEvents := make(map[string][]map[string]interface{})
	for _, event := range events {
		if userID, ok := event["user_id"].(string); ok {
			userEvents[userID] = append(userEvents[userID], event)
		}
	}

	// Process events by user
	for userID, events := range userEvents {
		log.Printf("Processing %d events for user %s", len(events), userID)

		// Process user's events
		for _, event := range events {
			if err := processEvent(ctx, event, EventMetadata{}); err != nil {
				log.Printf("Error processing event for user %s: %v", userID, err)
			}
		}
	}

	return nil
}

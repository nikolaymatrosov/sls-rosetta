package main

import "time"

// YDSEvent represents the event structure sent by YDS trigger
type YDSEvent struct {
	Messages []YDSMessage `json:"messages"`
}

// YDSMessage represents a single message from the stream
type YDSMessage struct {
	EventMetadata EventMetadata `json:"event_metadata"`
	Details       YDSDetails    `json:"details"`
}

// EventMetadata contains metadata about the event
type EventMetadata struct {
	EventID   string    `json:"event_id"`
	EventType string    `json:"event_type"`
	CreatedAt time.Time `json:"created_at"`
	CloudID   string    `json:"cloud_id"`
	FolderID  string    `json:"folder_id"`
}

// YDSDetails contains the actual data from the stream
type YDSDetails struct {
	StreamID string `json:"stream_id"`
	Data     string `json:"data"`
	// Additional fields as needed for YDS
}

// YDSResponse represents the response from the consumer function
type YDSResponse struct {
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
}

// ProducerRequest represents the request to the producer function
type ProducerRequest struct {
	Message string `json:"message"`
	UserID  string `json:"user_id"`
	Action  string `json:"action"`
}

// ProducerResponse represents the response from the producer function
type ProducerResponse struct {
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
	StreamID   string `json:"stream_id,omitempty"`
}

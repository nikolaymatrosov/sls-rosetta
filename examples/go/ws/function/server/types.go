package main

import (
	"encoding/json"
	"fmt"
	"time"
)

// WebSocket event types from API Gateway

// WebSocketRequest is a struct that represents the structure of a WebSocket request from API Gateway.
type WebSocketRequest struct {
	// RequestContext contains information about the WebSocket request context.
	RequestContext struct {
		// ConnectionID is the unique identifier for the WebSocket connection.
		ConnectionID string `json:"connectionId"`
		// ConnectedAt is the timestamp when the connection was established (for CONNECT events).
		ConnectedAt int64 `json:"connectedAt,omitempty"`
		// MessageID is the unique identifier for the message (for MESSAGE events).
		MessageID string `json:"messageId,omitempty"`
		// DisconnectReason explains why the connection was closed (for DISCONNECT events).
		DisconnectReason string `json:"disconnectReason,omitempty"`
		// DisconnectStatusCode is the status code for the disconnect (for DISCONNECT events).
		DisconnectStatusCode int `json:"disconnectStatusCode,omitempty"`
		// EventType indicates the type of WebSocket event (CONNECT, MESSAGE, DISCONNECT).
		EventType string `json:"eventType"`
	} `json:"requestContext"`
	// QueryStringParameters are the query parameters from the connection URL (for CONNECT events).
	QueryStringParameters map[string]string `json:"queryStringParameters,omitempty"`
	// Headers are the headers from the connection request (for CONNECT events).
	Headers map[string]string `json:"headers,omitempty"`
	// Body is the message body (for MESSAGE events).
	Body string `json:"body,omitempty"`
	// IsBase64Encoded indicates whether the body is Base64 encoded (for MESSAGE events).
	IsBase64Encoded bool `json:"isBase64Encoded,omitempty"`
}

// ConnectEvent represents the WebSocket CONNECT event from API Gateway
type ConnectEvent struct {
	RequestContext struct {
		ConnectionID string `json:"connectionId"`
		ConnectedAt  int64  `json:"connectedAt"`
		EventType    string `json:"eventType"`
	} `json:"requestContext"`
	QueryStringParameters map[string]string `json:"queryStringParameters"`
	Headers               map[string]string `json:"headers"`
}

// MessageEvent represents the WebSocket MESSAGE event from API Gateway
type MessageEvent struct {
	RequestContext struct {
		ConnectionID string `json:"connectionId"`
		MessageID    string `json:"messageId"`
		EventType    string `json:"eventType"`
	} `json:"requestContext"`
	Body            string `json:"body"`
	IsBase64Encoded bool   `json:"isBase64Encoded"`
}

// DisconnectEvent represents the WebSocket DISCONNECT event from API Gateway
type DisconnectEvent struct {
	RequestContext struct {
		ConnectionID         string `json:"connectionId"`
		DisconnectReason     string `json:"disconnectReason"`
		DisconnectStatusCode int    `json:"disconnectStatusCode"`
		EventType            string `json:"eventType"`
	} `json:"requestContext"`
}

// Message protocol types (client to server)

// ClientMessage represents a message sent from the client
type ClientMessage struct {
	Type      string `json:"type"` // "SEND" or "DISCONNECT"
	Content   string `json:"content,omitempty"`
	Timestamp string `json:"timestamp"`
}

// Message protocol types (server to client)

const (
	MessageTypeConnected  = "CONNECTED"
	MessageTypeBroadcast  = "BROADCAST"
	MessageTypeUserJoined = "USER_JOINED"
	MessageTypeUserLeft   = "USER_LEFT"
	MessageTypeError      = "ERROR"
	MessageTypeAck        = "ACK"
)

// ServerMessage represents a message sent to the client
type ServerMessage struct {
	Type         string `json:"type"`
	UserID       string `json:"userId,omitempty"`
	Content      string `json:"content,omitempty"`
	Message      string `json:"message,omitempty"`      // for ERROR
	Code         string `json:"code,omitempty"`         // for ERROR
	OriginalType string `json:"originalType,omitempty"` // for ACK
	Timestamp    string `json:"timestamp"`
}

// TriggerMessage represents the wrapper format for Data Streams trigger
type TriggerMessage struct {
	Messages []ServerMessage `json:"messages"`
}

// Connection represents a stored connection in YDB
type Connection struct {
	ConnectionID string    `json:"connection_id"`
	UserID       string    `json:"user_id"`
	ConnectedAt  time.Time `json:"connected_at"`
}

// Response types for the Lambda handler

// Response represents the HTTP response from the Lambda function
type Response struct {
	StatusCode int               `json:"statusCode"`
	Body       string            `json:"body,omitempty"`
	Headers    map[string]string `json:"headers,omitempty"`
}

// Helper functions

// ParseConnectEvent parses a CONNECT event from the request body
func ParseConnectEvent(body []byte) (*ConnectEvent, error) {
	var event ConnectEvent
	if err := json.Unmarshal(body, &event); err != nil {
		return nil, fmt.Errorf("failed to parse connect event: %w", err)
	}
	return &event, nil
}

// ParseMessageEvent parses a MESSAGE event from the request body
func ParseMessageEvent(body []byte) (*MessageEvent, error) {
	var event MessageEvent
	if err := json.Unmarshal(body, &event); err != nil {
		return nil, fmt.Errorf("failed to parse message event: %w", err)
	}
	return &event, nil
}

// ParseDisconnectEvent parses a DISCONNECT event from the request body
func ParseDisconnectEvent(body []byte) (*DisconnectEvent, error) {
	var event DisconnectEvent
	if err := json.Unmarshal(body, &event); err != nil {
		return nil, fmt.Errorf("failed to parse disconnect event: %w", err)
	}
	return &event, nil
}

// NewSuccessResponse creates a successful HTTP response
func NewSuccessResponse(sm ServerMessage) *Response {
	jsonBody, err := json.Marshal(sm)
	if err != nil {
		return &Response{
			StatusCode: 500,
			Body:       fmt.Sprintf("Failed to marshal response body: %v", err),
		}
	}
	return &Response{
		StatusCode: 200,
		Body:       string(jsonBody),
	}
}

// NewErrorResponse creates an error HTTP response
func NewErrorResponse(statusCode int, message string) *Response {
	return &Response{
		StatusCode: statusCode,
		Body:       message,
	}
}

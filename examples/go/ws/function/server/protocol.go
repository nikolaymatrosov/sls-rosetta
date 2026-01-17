package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"
)

// Message creation functions

// CreateConnectedMessage creates a CONNECTED message for the client
func CreateConnectedMessage(userID string) ServerMessage {
	return ServerMessage{
		Type:      MessageTypeConnected,
		UserID:    userID,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}
}

// CreateBroadcastMessage creates a BROADCAST message for all clients
func CreateBroadcastMessage(userID, content string) ServerMessage {
	return ServerMessage{
		Type:      MessageTypeBroadcast,
		UserID:    userID,
		Content:   content,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}
}

// CreateUserJoinedMessage creates a USER_JOINED message
func CreateUserJoinedMessage(userID string) ServerMessage {
	return ServerMessage{
		Type:      MessageTypeUserJoined,
		UserID:    userID,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}
}

// CreateUserLeftMessage creates a USER_LEFT message
func CreateUserLeftMessage(userID string) ServerMessage {
	return ServerMessage{
		Type:      MessageTypeUserLeft,
		UserID:    userID,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}
}

// CreateErrorMessage creates an ERROR message
func CreateErrorMessage(message, code string) ServerMessage {
	return ServerMessage{
		Type:      MessageTypeError,
		Message:   message,
		Code:      code,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}
}

// CreateAckMessage creates an ACK message
func CreateAckMessage(originalType string) ServerMessage {
	return ServerMessage{
		Type:         MessageTypeAck,
		OriginalType: originalType,
		Timestamp:    time.Now().UTC().Format(time.RFC3339),
	}
}

// ParseClientMessage parses and validates a client message from JSON
func ParseClientMessage(data []byte) (*ClientMessage, error) {
	var msg ClientMessage
	if err := json.Unmarshal(data, &msg); err != nil {
		return nil, fmt.Errorf("failed to parse client message: %w", err)
	}

	// Validate message type
	if msg.Type != "SEND" && msg.Type != "DISCONNECT" {
		return nil, fmt.Errorf("invalid message type: %s", msg.Type)
	}

	// Validate SEND messages have content
	if msg.Type == "SEND" && msg.Content == "" {
		return nil, fmt.Errorf("SEND message must have content")
	}

	return &msg, nil
}

// DecodeMessageBody decodes the message body, handling base64 encoding if necessary
func DecodeMessageBody(body string, isBase64Encoded bool) ([]byte, error) {
	if isBase64Encoded {
		decoded, err := base64.StdEncoding.DecodeString(body)
		if err != nil {
			return nil, fmt.Errorf("failed to decode base64 body: %w", err)
		}
		return decoded, nil
	}
	return []byte(body), nil
}

// SerializeMessage serializes a ServerMessage to JSON
func SerializeMessage(msg ServerMessage) ([]byte, error) {
	data, err := json.Marshal(msg)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize message: %w", err)
	}
	return data, nil
}

// SerializeTriggerMessage serializes a TriggerMessage (wrapper with messages array)
func SerializeTriggerMessage(messages []ServerMessage) ([]byte, error) {
	wrapper := TriggerMessage{
		Messages: messages,
	}
	data, err := json.Marshal(wrapper)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize trigger message: %w", err)
	}
	return data, nil
}

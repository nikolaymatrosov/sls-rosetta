package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/google/uuid"
	"github.com/ydb-platform/ydb-go-sdk/v3"
	"github.com/ydb-platform/ydb-go-sdk/v3/topic/topicoptions"
	"github.com/ydb-platform/ydb-go-sdk/v3/topic/topicwriter"
	yc "github.com/ydb-platform/ydb-go-yc"
)

var logger *slog.Logger

func init() {
	logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelDebug,
	}))
}

// WebSocketEventHandler is the main entry point for the Lambda function.
// It handles WebSocket events (CONNECT, MESSAGE, DISCONNECT) from API Gateway.
//
//goland:noinspection ALL
func WebSocketEventHandler(ctx context.Context, event *WebSocketRequest) (*Response, error) {
	logger.Info("Received WebSocket event", "eventType", event.RequestContext.EventType)
	logger.Debug("Event details", "connectionId", event.RequestContext.ConnectionID)

	// Route to appropriate handler based on event type
	switch event.RequestContext.EventType {
	case "CONNECT":
		return handleConnectEvent(ctx, event)
	case "MESSAGE":
		return handleMessageEvent(ctx, event)
	case "DISCONNECT":
		return handleDisconnectEvent(ctx, event)
	default:
		logger.Error("Unknown event type", "eventType", event.RequestContext.EventType)
		return NewErrorResponse(400, fmt.Sprintf("Unknown event type: %s", event.RequestContext.EventType)), nil
	}
}

func handleConnectEvent(ctx context.Context, event *WebSocketRequest) (*Response, error) {
	logger.Info("Handling CONNECT event")

	// Get user ID from query parameters or generate one
	userID := event.QueryStringParameters["user_id"]
	if userID == "" {
		userID = uuid.New().String()
		logger.Info("Generated new user ID", "userId", userID)
	}

	connectionID := event.RequestContext.ConnectionID
	logger.Info("Connection details", "connectionId", connectionID, "userId", userID)

	// Initialize YDB connection
	db, err := initYDB(ctx)
	if err != nil {
		logger.Error("Failed to initialize YDB", "error", err)
		return NewErrorResponse(500, "Database connection failed"), nil
	}
	defer func() {
		_ = db.Close(ctx)
	}()

	// Store connection in database
	if err := StoreConnection(ctx, db, connectionID, userID); err != nil {
		logger.Error("Failed to store connection", "error", err)
		return NewErrorResponse(500, "Failed to store connection"), nil
	}

	logger.Info("Connection stored successfully")

	// Publish USER_JOINED message to topic
	if err := publishMessage(ctx, db, CreateUserJoinedMessage(userID)); err != nil {
		logger.Error("Failed to publish USER_JOINED message", "error", err)
		// Don't fail the connection if publish fails
	}

	return NewSuccessResponse(CreateConnectedMessage(userID)), nil
}

func handleMessageEvent(ctx context.Context, event *WebSocketRequest) (*Response, error) {
	logger.Info("Handling MESSAGE event")

	connectionID := event.RequestContext.ConnectionID
	logger.Info("Message details", "connectionId", connectionID)

	// Decode message body
	messageBody, err := DecodeMessageBody(event.Body, event.IsBase64Encoded)
	if err != nil {
		logger.Error("Failed to decode message body", "error", err)
		return NewErrorResponse(400, "Invalid message encoding"), nil
	}

	// Parse client message
	clientMsg, err := ParseClientMessage(messageBody)
	if err != nil {
		logger.Error("Failed to parse client message", "error", err)
		return NewErrorResponse(400, "Invalid client message"), nil
	}

	logger.Info("Received client message", "type", clientMsg.Type, "content", clientMsg.Content)

	// Initialize YDB connection
	db, err := initYDB(ctx)
	if err != nil {
		logger.Error("Failed to initialize YDB", "error", err)
		return NewErrorResponse(500, "Database connection failed"), nil
	}
	defer func() {
		_ = db.Close(ctx)
	}()

	// Get user ID from connection
	userID, err := GetUserIDByConnectionID(ctx, db, connectionID)
	if err != nil {
		logger.Error("Failed to get user ID", "error", err)
		return NewErrorResponse(404, "Connection not found"), nil
	}

	// Handle different message types
	switch clientMsg.Type {
	case "SEND":
		// Publish BROADCAST message to topic
		broadcastMsg := CreateBroadcastMessage(userID, clientMsg.Content)
		if err := publishMessage(ctx, db, broadcastMsg); err != nil {
			logger.Error("Failed to publish broadcast message", "error", err)
			return NewErrorResponse(500, "Failed to broadcast message"), nil
		}
		logger.Info("Broadcast message published successfully")

	case "DISCONNECT":
		// Handle graceful disconnect
		logger.Info("Client requested disconnect")
		if err := RemoveConnectionByID(ctx, db, connectionID); err != nil {
			logger.Error("Failed to remove connection", "error", err)
		}
		// Publish USER_LEFT message
		if err := publishMessage(ctx, db, CreateUserLeftMessage(userID)); err != nil {
			logger.Error("Failed to publish USER_LEFT message", "error", err)
		}
	}

	return NewSuccessResponse(CreateAckMessage(clientMsg.Type)), nil
}

func handleDisconnectEvent(ctx context.Context, event *WebSocketRequest) (*Response, error) {
	logger.Info("Handling DISCONNECT event")

	connectionID := event.RequestContext.ConnectionID
	logger.Info("Disconnect details", "connectionId", connectionID, "reason", event.RequestContext.DisconnectReason)

	// Initialize YDB connection
	db, err := initYDB(ctx)
	if err != nil {
		logger.Error("Failed to initialize YDB", "error", err)
		return NewErrorResponse(500, "Database connection failed"), nil
	}
	defer func() {
		_ = db.Close(ctx)
	}()

	// Get user ID before removing connection
	userID, err := GetUserIDByConnectionID(ctx, db, connectionID)
	if err != nil {
		logger.Error("Failed to get user ID", "error", err)
		// Continue with cleanup even if we can't find the user
	}

	// Remove connection from database
	if err := RemoveConnectionByID(ctx, db, connectionID); err != nil {
		logger.Error("Failed to remove connection", "error", err)
		return NewErrorResponse(500, "Failed to remove connection"), nil
	}

	logger.Info("Connection removed successfully")

	// Publish USER_LEFT message if we have a user ID
	if userID != "" {
		if err := publishMessage(ctx, db, CreateUserLeftMessage(userID)); err != nil {
			logger.Error("Failed to publish USER_LEFT message", "error", err)
			// Don't fail the disconnect if publish fails
		}
	}

	return NewSuccessResponse(CreateAckMessage("DISCONNECT")), nil
}

// initYDB initializes a connection to YDB
func initYDB(ctx context.Context) (*ydb.Driver, error) {
	connectionString := os.Getenv("YDB_CONNECTION_STRING")

	if connectionString == "" {
		return nil, fmt.Errorf("YDB_CONNECTION_STRING and YDB_DATABASE environment variables must be set")
	}

	logger.Debug("Connecting to YDB", "connectionString", connectionString)

	db, err := ydb.Open(ctx, connectionString, yc.WithMetadataCredentials())
	if err != nil {
		return nil, fmt.Errorf("failed to connect to YDB: %w", err)
	}

	return db, nil
}

// publishMessage publishes a message to the YDB topic
func publishMessage(ctx context.Context, db *ydb.Driver, message ServerMessage) error {
	topicPath := os.Getenv("BROADCAST_TOPIC")
	if topicPath == "" {
		return fmt.Errorf("BROADCAST_TOPIC environment variable not set")
	}

	logger.Debug("Publishing message to topic", "topic", topicPath, "messageType", message.Type)

	// Serialize message
	data, err := SerializeMessage(message)
	if err != nil {
		return fmt.Errorf("failed to serialize message: %w", err)
	}

	// Create topic writer
	writer, err := db.Topic().StartWriter(topicPath,
		topicoptions.WithWriterProducerID("ws-handler"),
	)
	if err != nil {
		return fmt.Errorf("failed to create topic writer: %w", err)
	}
	defer func() {
		_ = writer.Close(ctx)
	}()

	// Write message
	err = writer.Write(ctx, topicwriter.Message{
		Data: bytes.NewBuffer(data),
	})
	if err != nil {
		return fmt.Errorf("failed to write message to topic: %w", err)
	}

	logger.Debug("Message published successfully")
	return nil
}

// writeResponse writes an HTTP response
func writeResponse(w http.ResponseWriter, resp Response) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.StatusCode)

	if resp.Body != "" {
		respBody := map[string]interface{}{
			"statusCode": resp.StatusCode,
			"body":       resp.Body,
		}
		if resp.Headers != nil {
			respBody["headers"] = resp.Headers
		}
		_ = json.NewEncoder(w).Encode(respBody)
	}
}

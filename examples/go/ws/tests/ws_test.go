package tests

import (
	"encoding/json"
	"net/url"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Message types
type ClientMessage struct {
	Type      string `json:"type"`
	Content   string `json:"content,omitempty"`
	Timestamp string `json:"timestamp"`
}

type ServerMessage struct {
	Type         string `json:"type"`
	UserID       string `json:"userId,omitempty"`
	Content      string `json:"content,omitempty"`
	Message      string `json:"message,omitempty"`
	Code         string `json:"code,omitempty"`
	OriginalType string `json:"originalType,omitempty"`
	Timestamp    string `json:"timestamp"`
}

type TriggerMessage struct {
	Messages []ServerMessage `json:"messages"`
}

func TestWebSocketE2E(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping E2E test in short mode")
	}

	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "../tf",
		Vars: map[string]interface{}{
			"folder_id": getFolderID(t),
		},
	})

	// Deploy infrastructure
	//defer terraform.Destroy(t, terraformOptions)
	terraform.InitAndApply(t, terraformOptions)

	// Get WebSocket URL from outputs
	wsURL := terraform.Output(t, terraformOptions, "websocket_url")
	require.NotEmpty(t, wsURL, "WebSocket URL should not be empty")

	t.Run("SingleClientConnection", func(t *testing.T) {
		testSingleClientConnection(t, wsURL)
	})

	t.Run("MultipleClientsAndBroadcast", func(t *testing.T) {
		testMultipleClientsAndBroadcast(t, wsURL)
	})

	t.Run("ClientReconnection", func(t *testing.T) {
		testClientReconnection(t, wsURL)
	})
}

func testSingleClientConnection(t *testing.T, baseURL string) {
	userID := uuid.New().String()
	conn := connectClient(t, baseURL, userID)
	defer conn.Close()

	// Should receive CONNECTED message
	msg := readServerMessage(t, conn, 5*time.Second, "CONNECTED")
	assert.Equal(t, userID, msg.UserID)

	// Send a message
	sendClientMessage(t, conn, ClientMessage{
		Type:      "SEND",
		Content:   "Hello, WebSocket!",
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	})

	// Should receive BROADCAST of own message (via topic trigger)
	msg = readServerMessage(t, conn, 10*time.Second, "BROADCAST")
	assert.Equal(t, userID, msg.UserID)
	assert.Equal(t, "Hello, WebSocket!", msg.Content)
}

func testMultipleClientsAndBroadcast(t *testing.T, baseURL string) {
	userID1 := uuid.New().String()
	userID2 := uuid.New().String()

	// Connect first client
	conn1 := connectClient(t, baseURL, userID1)
	defer conn1.Close()

	// Wait for CONNECTED message
	readServerMessage(t, conn1, 5*time.Second, "CONNECTED")

	// Connect second client
	conn2 := connectClient(t, baseURL, userID2)
	defer conn2.Close()

	// Wait for CONNECTED message on second client
	readServerMessage(t, conn2, 5*time.Second, "CONNECTED")

	// First client should receive USER_JOINED for second client
	msg := readServerMessage(t, conn1, 10*time.Second, "USER_JOINED")
	assert.Equal(t, userID2, msg.UserID)

	// Send message from first client
	sendClientMessage(t, conn1, ClientMessage{
		Type:      "SEND",
		Content:   "Hello from client 1",
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	})

	// Both clients should receive the broadcast
	msg1 := readServerMessage(t, conn1, 10*time.Second, "BROADCAST")
	assert.Equal(t, "Hello from client 1", msg1.Content)

	msg2 := readServerMessage(t, conn2, 10*time.Second, "BROADCAST")
	assert.Equal(t, "Hello from client 1", msg2.Content)

	// Close second client
	conn2.Close()

	// First client should receive USER_LEFT
	msg = readServerMessage(t, conn1, 10*time.Second, "USER_LEFT")
	assert.Equal(t, userID2, msg.UserID)
}

func testClientReconnection(t *testing.T, baseURL string) {
	userID := uuid.New().String()

	// First connection
	conn1 := connectClient(t, baseURL, userID)
	readServerMessage(t, conn1, 5*time.Second, "CONNECTED")
	conn1.Close()

	// Wait a bit for cleanup
	time.Sleep(2 * time.Second)

	// Reconnect with same user ID (should replace old connection)
	conn2 := connectClient(t, baseURL, userID)
	defer conn2.Close()

	msg := readServerMessage(t, conn2, 5*time.Second, "CONNECTED")
	assert.Equal(t, userID, msg.UserID)

	// Should be able to send messages
	sendClientMessage(t, conn2, ClientMessage{
		Type:      "SEND",
		Content:   "After reconnection",
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	})

	msg = readServerMessage(t, conn2, 10*time.Second, "BROADCAST")
	assert.Equal(t, "After reconnection", msg.Content)
}

// Helper functions

func connectClient(t *testing.T, baseURL, userID string) *websocket.Conn {
	u, err := url.Parse(baseURL)
	require.NoError(t, err)

	q := u.Query()
	q.Set("user_id", userID)
	u.RawQuery = q.Encode()

	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	require.NoError(t, err)

	return conn
}

func sendClientMessage(t *testing.T, conn *websocket.Conn, msg ClientMessage) {
	data, err := json.Marshal(msg)
	require.NoError(t, err)

	err = conn.WriteMessage(websocket.TextMessage, data)
	require.NoError(t, err)
}

func readServerMessage(t *testing.T, conn *websocket.Conn, timeout time.Duration, expectedType string) ServerMessage {
	deadline := time.Now().Add(timeout)

	for {
		conn.SetReadDeadline(deadline)

		_, message, err := conn.ReadMessage()
		require.NoError(t, err)

		// Try to parse as TriggerMessage first
		var triggerMsg TriggerMessage
		if err := json.Unmarshal(message, &triggerMsg); err == nil && len(triggerMsg.Messages) > 0 {
			msg := triggerMsg.Messages[0]
			if msg.Type != expectedType {
				t.Logf("Ignoring message with type %s (expected %s): %+v", msg.Type, expectedType, msg)
				continue
			}
			return msg
		}

		// Parse as ServerMessage
		var serverMsg ServerMessage
		err = json.Unmarshal(message, &serverMsg)
		require.NoError(t, err)

		if serverMsg.Type != expectedType {
			t.Logf("Ignoring message with type %s (expected %s): %+v", serverMsg.Type, expectedType, serverMsg)
			continue
		}

		return serverMsg
	}
}

func getFolderID(t *testing.T) string {
	// Get folder ID from environment or configuration
	// This should be set in CI/CD or local testing environment
	// For now, we'll require it to be passed as an environment variable
	folderID := terraform.GetVariableAsStringFromVarFile(t, "../tf/.tfvars", "folder_id")
	if folderID == "" {
		t.Skip("folder_id not configured in terraform.tfvars")
	}
	return folderID
}

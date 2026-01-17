package main

import (
	"bufio"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/fatih/color"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

// Message types matching the server protocol

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

var (
	green  = color.New(color.FgGreen).SprintFunc()
	yellow = color.New(color.FgYellow).SprintFunc()
	red    = color.New(color.FgRed).SprintFunc()
	blue   = color.New(color.FgBlue).SprintFunc()
	cyan   = color.New(color.FgCyan).SprintFunc()
)

func main() {
	// Parse command-line flags
	wsURL := flag.String("url", "", "WebSocket URL (e.g., wss://example.com/ws)")
	userID := flag.String("user-id", "", "User ID (optional, will be generated if not provided)")
	flag.Parse()

	if *wsURL == "" {
		log.Fatal("WebSocket URL is required. Use -url flag.")
	}

	// Generate user ID if not provided
	if *userID == "" {
		*userID = uuid.New().String()
		fmt.Printf("%s Generated user ID: %s\n", yellow("⚠"), *userID)
	}

	// Add user_id as query parameter
	u, err := url.Parse(*wsURL)
	if err != nil {
		log.Fatalf("Invalid WebSocket URL: %v", err)
	}
	q := u.Query()
	q.Set("user_id", *userID)
	u.RawQuery = q.Encode()

	fmt.Printf("%s Connecting to %s\n", blue("→"), u.String())

	// Connect to WebSocket
	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	fmt.Printf("%s Connected successfully!\n", green("✓"))
	fmt.Printf("%s Type your messages and press Enter to send. Ctrl+C to exit.\n\n", cyan("ℹ"))

	// Create context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle signals for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigChan
		fmt.Printf("\n%s Shutting down...\n", yellow("⚠"))
		cancel()
	}()

	// Start reading messages in a goroutine
	go readMessages(ctx, conn)

	// Read input from user
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("> ")
	for scanner.Scan() {
		select {
		case <-ctx.Done():
			return
		default:
		}

		text := strings.TrimSpace(scanner.Text())
		if text == "" {
			fmt.Print("> ")
			continue
		}

		// Send message
		msg := ClientMessage{
			Type:      "SEND",
			Content:   text,
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		}

		data, err := json.Marshal(msg)
		if err != nil {
			fmt.Printf("%s Failed to serialize message: %v\n", red("✗"), err)
			fmt.Print("> ")
			continue
		}

		err = conn.WriteMessage(websocket.TextMessage, data)
		if err != nil {
			fmt.Printf("%s Failed to send message: %v\n", red("✗"), err)
			fmt.Print("> ")
			continue
		}

		fmt.Print("> ")
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("%s Error reading input: %v\n", red("✗"), err)
	}
}

func readMessages(ctx context.Context, conn *websocket.Conn) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		_, message, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
				fmt.Printf("\n%s Connection closed\n", yellow("⚠"))
			} else {
				fmt.Printf("\n%s Error reading message: %v\n", red("✗"), err)
			}
			return
		}

		// Try to parse as TriggerMessage (wrapper format)
		var triggerMsg TriggerMessage
		if err := json.Unmarshal(message, &triggerMsg); err == nil && len(triggerMsg.Messages) > 0 {
			// This is a trigger message with multiple server messages
			for _, msg := range triggerMsg.Messages {
				printServerMessage(msg)
			}
			continue
		}

		// Try to parse as ServerMessage
		var serverMsg ServerMessage
		if err := json.Unmarshal(message, &serverMsg); err != nil {
			fmt.Printf("\n%s Failed to parse message: %v\n", red("✗"), err)
			fmt.Printf("%s Raw message: %s\n", yellow("⚠"), string(message))
			continue
		}

		printServerMessage(serverMsg)
	}
}

func printServerMessage(msg ServerMessage) {
	timestamp := formatTimestamp(msg.Timestamp)

	switch msg.Type {
	case "CONNECTED":
		fmt.Printf("\n%s [%s] Connected as %s\n", green("✓"), timestamp, cyan(msg.UserID))

	case "BROADCAST":
		fmt.Printf("\n%s [%s] %s: %s\n", blue("💬"), timestamp, cyan(msg.UserID), msg.Content)

	case "USER_JOINED":
		fmt.Printf("\n%s [%s] User %s joined\n", green("→"), timestamp, cyan(msg.UserID))

	case "USER_LEFT":
		fmt.Printf("\n%s [%s] User %s left\n", yellow("←"), timestamp, cyan(msg.UserID))

	case "ERROR":
		errorMsg := msg.Message
		if msg.Code != "" {
			errorMsg = fmt.Sprintf("%s (%s)", msg.Message, msg.Code)
		}
		fmt.Printf("\n%s [%s] Error: %s\n", red("✗"), timestamp, errorMsg)

	case "ACK":
		fmt.Printf("\n%s [%s] Acknowledged: %s\n", green("✓"), timestamp, msg.OriginalType)

	default:
		fmt.Printf("\n%s [%s] Unknown message type: %s\n", yellow("?"), timestamp, msg.Type)
	}

	fmt.Print("> ")
}

func formatTimestamp(ts string) string {
	t, err := time.Parse(time.RFC3339, ts)
	if err != nil {
		return ts
	}
	return t.Local().Format("15:04:05")
}

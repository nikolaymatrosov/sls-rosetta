package main

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
)

func Handler(rw http.ResponseWriter, req *http.Request) {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelDebug,
	}))
	// Simple usage
	logger.Info("hello", "count", 3)
	logger.Warn("warn message")
	logger.Error("error message", "key", "value")

	// With context
	ctx := context.WithValue(req.Context(), "requestID", "12345")
	logger.DebugContext(ctx, "debug message")

	rw.Header().Set("X-Custom-Header", "Test")
	rw.WriteHeader(200)
	name := req.URL.Query().Get("name")
	_, _ = io.WriteString(rw, fmt.Sprintf("Hello, %s!", name))
}

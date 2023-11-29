package main

// TracingContext holds the tracing information for a request.
type TracingContext struct {
	ParentSpanID string `json:"parent_span_id,omitempty"` // The ID of the parent span.
	SpanID       string `json:"span_id,omitempty"`        // The ID of the current span.
	TraceID      string `json:"trace_id,omitempty"`       // The ID of the trace.
}

// ObjectStorageMessageMetadata holds the metadata for an object storage event.
type ObjectStorageMessageMetadata struct {
	CreatedAt      string         `json:"created_at,omitempty"` // The creation timestamp of the event.
	CloudID        string         `json:"cloud_id,omitempty"`   // The ID of the cloud where the event occurred.
	FolderID       string         `json:"folder_id,omitempty"`  // The ID of the folder where the event occurred.
	TracingContext TracingContext `json:"tracing_context"`      // The tracing context for the event.
}

// ObjectStorageMessageDetails holds the details for an object storage event.
type ObjectStorageMessageDetails struct {
	BucketID string `json:"bucket_id,omitempty"` // The ID of the bucket where the event occurred.
	ObjectID string `json:"object_id,omitempty"` // The ID of the object involved in the event.
}

// ObjectStorageMessage represents an event in object storage.
type ObjectStorageMessage struct {
	Metadata ObjectStorageMessageMetadata `json:"metadata"` // The metadata for the event.
	Details  ObjectStorageMessageDetails  `json:"details"`  // The details of the event.
}

// ObjectStorageResponse represents the response from an object storage operation.
type ObjectStorageResponse struct {
	StatusCode int // The status code of the response.
}

// ObjectStorageEvent represents an event in object storage.
type ObjectStorageEvent struct {
	Messages []ObjectStorageMessage `json:"messages"` // The messages for the event.
}

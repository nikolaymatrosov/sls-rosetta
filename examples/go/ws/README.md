# WebSocket Example in Go

This example demonstrates how to build a real-time WebSocket application using Yandex Cloud services with Go. It implements a chat-like system where messages are broadcast to all connected clients via YDB Topics and Data Streams triggers.

## Architecture

The example follows a topic-based broadcasting pattern:

```
┌──────────┐         ┌─────────────────┐         ┌──────────────┐
│  Client  │◄───────►│  API Gateway    │◄───────►│   Function   │
│          │  WebSocket   (WebSocket)   │         │  (Handler)   │
└──────────┘         └─────────────────┘         └──────┬───────┘
                                                         │
                                                         ▼
                                                  ┌──────────────┐
                                                  │     YDB      │
                                                  │  Database    │
                                                  │ (connections)│
                                                  └──────────────┘
                                                         │
                                                         ▼
                                                  ┌──────────────┐
                                                  │  YDB Topic   │
                                                  │ (broadcast-  │
                                                  │   topic)     │
                                                  └──────┬───────┘
                                                         │
                                                         ▼
                                                  ┌──────────────┐
                                                  │ Data Streams │
                                                  │   Trigger    │
                                                  └──────┬───────┘
                                                         │
                                  ┌──────────────────────┴─────────────────────┐
                                  ▼                                            ▼
                          ┌──────────────┐                            ┌──────────────┐
                          │   Client 1   │                            │   Client 2   │
                          │  (broadcast) │                            │  (broadcast) │
                          └──────────────┘                            └──────────────┘
```

### Flow:

1. **Client connects** to API Gateway WebSocket route
2. **Handler stores** connection in YDB and publishes `USER_JOINED` to topic
3. **Data Streams trigger** reads from topic and broadcasts to all WebSocket connections
4. **Client sends message** → Handler publishes `BROADCAST` to topic → Trigger broadcasts to all clients
5. **Client disconnects** → Handler removes connection and publishes `USER_LEFT` to topic

## Features

- Real-time bidirectional communication via WebSocket
- Topic-based message broadcasting (scalable architecture)
- Connection state management in YDB
- Automatic connection cleanup with TTL
- One connection per user policy
- Interactive CLI client with colored output
- Comprehensive E2E tests

## Prerequisites

- Go 1.23 or later
- Terraform 1.0 or later
- Yandex Cloud CLI (`yc`) configured with valid credentials
- [goose](https://github.com/pressly/goose) for database migrations (`go install github.com/pressly/goose/v3/cmd/goose@latest`)
- Yandex Cloud account with sufficient permissions
- `YC_TOKEN` environment variable set with your Yandex Cloud IAM token

## Project Structure

```
examples/go/ws/
├── function/
│   ├── server/              # WebSocket handler function
│   │   ├── main.go          # Handler entry point
│   │   ├── types.go         # Type definitions
│   │   ├── protocol.go      # Message protocol
│   │   ├── database.go      # YDB operations
│   │   └── go.mod           # Go dependencies
│   └── client/              # WebSocket CLI client
│       ├── main.go          # Client application
│       └── go.mod           # Client dependencies
├── tf/                      # Terraform infrastructure
│   ├── main.tf              # Function deployment
│   ├── ydb.tf               # Database and topic
│   ├── apigateway.tf        # API Gateway setup
│   ├── iam.tf               # IAM roles
│   ├── api-gateway.yaml     # OpenAPI specification
│   ├── variables.tf         # Input variables
│   ├── outputs.tf           # Output values
│   └── terraform.tf         # Provider configuration
├── migrations/              # Database schema
│   └── 001_create_connections.sql
├── tests/                   # E2E tests
│   ├── ws_test.go
│   └── go.mod
├── README.md                # This file
└── .gitignore
```

## Message Protocol

### Client → Server Messages

```go
type ClientMessage struct {
    Type      string // "SEND" or "DISCONNECT"
    Content   string // Message content (for SEND)
    Timestamp string // ISO 8601 timestamp
}
```

### Server → Client Messages

```go
type ServerMessage struct {
    Type         string // "CONNECTED", "BROADCAST", "USER_JOINED", "USER_LEFT", "ERROR", "ACK"
    UserID       string // User ID (optional)
    Content      string // Message content (optional)
    Message      string // Error message (for ERROR)
    Code         string // Error code (for ERROR)
    OriginalType string // Original message type (for ACK)
    Timestamp    string // ISO 8601 timestamp
}
```

### Data Streams Trigger Wrapper

Messages from the Data Streams trigger are wrapped in:

```go
type TriggerMessage struct {
    Messages []ServerMessage
}
```

## Deployment

### 1. Configure Variables

Create a `.tfvars` file in the `tf/` directory:

```hcl
cloud_id  = "your-cloud-id-here"
folder_id = "your-folder-id-here"
```

Optional variables (with defaults):

```hcl
function_name        = "ws-go-handler"
database_name        = "ws-go-database"
gateway_name         = "ws-go-gateway"
service_account_name = "ws-go-function-sa"
topic_name           = "broadcast-topic"
topic_consumer_name  = "broadcast-consumer"
trigger_name         = "ws-go-broadcast-trigger"
```

### 2. Set Environment Variables

Export your Yandex Cloud IAM token (required for goose migrations):

```bash
export YC_TOKEN=$(yc iam create-token)
```

### 3. Deploy Infrastructure

```bash
cd tf/
terraform init
terraform apply
```

This will create:
- YDB serverless database with `connections` table
- YDB topic for message broadcasting
- WebSocket API Gateway
- Lambda function (Go handler)
- Data Streams trigger
- Service account with required IAM roles

### 4. Get WebSocket URL

After deployment completes:

```bash
terraform output websocket_url
```

Example output: `wss://d5d123abc456def.apigw.yandexcloud.net/ws`

## Usage

### Running the Client

```bash
cd function/client/

# Connect with auto-generated user ID
go run main.go -url wss://YOUR-GATEWAY-DOMAIN/ws

# Connect with specific user ID
go run main.go -url wss://YOUR-GATEWAY-DOMAIN/ws -user-id my-user-123
```

You can also use the command from Terraform output:

```bash
eval $(cd tf && terraform output -raw client_command)
```

### Client Commands

Once connected:
- Type a message and press Enter to send
- Press Ctrl+C to disconnect gracefully

### Example Session

```
⚠ Generated user ID: 8f7a3c2b-1d4e-4a9f-b6c8-9e2d1f3a5b7c
→ Connecting to wss://d5d123abc456def.apigw.yandexcloud.net/ws?user_id=8f7a3c2b-1d4e-4a9f-b6c8-9e2d1f3a5b7c
✓ Connected successfully!
ℹ Type your messages and press Enter to send. Ctrl+C to exit.

✓ [14:25:10] Connected as 8f7a3c2b-1d4e-4a9f-b6c8-9e2d1f3a5b7c
> Hello, everyone!

💬 [14:25:12] 8f7a3c2b-1d4e-4a9f-b6c8-9e2d1f3a5b7c: Hello, everyone!

→ [14:25:20] User abc-def-123 joined

💬 [14:25:25] abc-def-123: Hi there!
>
```

## Testing

### Run E2E Tests

```bash
cd tests/
go test -v
```

The tests will:
1. Deploy infrastructure via Terraform
2. Test single client connection
3. Test multiple clients and broadcasting
4. Test client reconnection
5. Clean up infrastructure


## Infrastructure Components

### YDB Database
- **Type**: Serverless
- **Schema**: `connections` table with TTL
- **Purpose**: Store active WebSocket connections

### YDB Topic
- **Name**: `broadcast-topic`
- **Purpose**: Message queue for broadcasting
- **Consumer**: `broadcast-consumer`

### API Gateway
- **Type**: WebSocket
- **Routes**: `/ws`
- **Events**: CONNECT, MESSAGE, DISCONNECT

### Function
- **Runtime**: golang123 (Go 1.23)
- **Memory**: 256 MB
- **Timeout**: 30 seconds
- **Entrypoint**: `index.Handler`

### Data Streams Trigger
- **Source**: `broadcast-topic`
- **Target**: WebSocket API Gateway (broadcast)
- **Batch**: 10 messages, 1 second cutoff

### Service Account Roles
- `ydb.editor`: Database operations
- `api-gateway.websocketWriter`: Send messages to connections
- `api-gateway.websocketBroadcaster`: Broadcast to all connections
- `serverless.functions.invoker`: Invoke functions (for trigger)
- `yds.admin`: Manage topics and triggers

## Environment Variables

The function uses these environment variables (set by Terraform):

- `YDB_CONNECTION_STRING`: YDB database endpoint
- `YDB_DATABASE`: Database path
- `BROADCAST_TOPIC`: Full topic path

## Cleanup

To destroy all resources:

```bash
cd tf/
terraform destroy
```

## Troubleshooting

### Connection Issues

**Problem**: Client can't connect to WebSocket

**Solutions**:
- Verify the WebSocket URL is correct
- Check that API Gateway is deployed: `yc api-gateway get <gateway-id>`
- Check function logs: `yc serverless function logs <function-id>`

### Messages Not Broadcasting

**Problem**: Messages sent but not received by other clients

**Solutions**:
- Verify Data Streams trigger is active: `yc serverless trigger list`
- Check YDB topic has messages: Use YDB console
- Check trigger logs: `yc logging read default --filter 'resource_id=<trigger-id>'`
- Verify service account has required roles

### Database Errors

**Problem**: Function returns 500 errors related to database

**Solutions**:
- Check migrations ran successfully: Query YDB console for `connections` table
- Verify service account has `ydb.editor` role
- Check YDB database is active: `yc ydb database get <database-id>`

## References

- [Yandex Cloud API Gateway WebSocket](https://cloud.yandex.com/docs/api-gateway/concepts/extensions/websocket)
- [YDB Go SDK](https://github.com/ydb-platform/ydb-go-sdk)
- [YDB Topics](https://cloud.yandex.com/docs/ydb/concepts/topic)
- [Data Streams Triggers](https://cloud.yandex.com/docs/functions/concepts/trigger/data-streams-trigger)

## License

This example is part of the sls-rosetta project.

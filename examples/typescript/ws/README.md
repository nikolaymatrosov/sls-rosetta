# TypeScript WebSocket Broadcast Example

This example demonstrates a WebSocket broadcast server using Yandex Cloud API Gateway and Cloud Functions with YDB for connection storage.

## Features

- WebSocket connection handling (connect/message/disconnect)
- Connection storage in serverless YDB
- Message broadcasting to all connected clients
- Single function handling all WebSocket events
- Browser and CLI demo clients

## Architecture

```
┌─────────────┐     ┌──────────────┐     ┌───────────────┐
│   Clients   │────▶│  API Gateway │────▶│   Function    │
│ (WebSocket) │◀────│  (WebSocket) │◀────│  (Node.js)    │
└─────────────┘     └──────────────┘     └───────┬───────┘
                           ▲                     │
                           │                     ▼
                    ┌──────┴──────┐       ┌─────────────┐
                    │  WebSocket  │       │     YDB     │
                    │ Management  │◀──────│ (Serverless)│
                    │     API     │       └─────────────┘
                    └─────────────┘
```

## Prerequisites

- Yandex Cloud CLI configured
- Terraform installed
- Node.js 18+ and npm
- [goose](https://github.com/pressly/goose) for database migrations

## Project Structure

```
examples/typescript/ws/
├── README.md
├── migrations/                # Goose database migrations
│   └── 001_create_connections.sql
├── function/
│   ├── package.json           # Shared dependencies
│   ├── tsconfig.json          # Shared TypeScript config
│   ├── server/
│   │   └── main.ts            # WebSocket handler
│   └── client/
│       ├── index.html         # Browser demo client
│       └── client.ts          # CLI demo client
├── dist/                      # Built output (gitignored)
├── environment/               # Terraform state (gitignored)
└── tf/
    ├── terraform.tf
    ├── variables.tf
    ├── ydb.tf
    ├── iam.tf
    ├── main.tf
    ├── api-gateway.yaml
    ├── apigateway.tf
    └── outputs.tf
```

## Deployment

1. Navigate to the function directory and install dependencies:

   ```bash
   cd examples/typescript/ws/function
   npm install
   ```

2. Initialize Terraform:

   ```bash
   cd ../tf
   terraform init
   ```

3. Set environment variables:

   ```bash
   export TF_VAR_cloud_id="your-cloud-id"
   export TF_VAR_folder_id="your-folder-id"
   export YC_TOKEN=$(yc iam create-token)
   ```

4. Deploy the infrastructure (migrations run automatically):

   ```bash
   terraform apply
   ```

5. Get the WebSocket URL:

   ```bash
   terraform output websocket_url
   ```

## Testing

### Browser Client

Open `function/client/index.html` in a browser, enter the WebSocket URL, and click Connect.

You can also pass the URL as a parameter:

```
file:///path/to/index.html?url=wss://your-gateway.apigw.yandexcloud.net/ws
```

### CLI Client

```bash
cd function
npx ts-node client/client.ts wss://your-gateway.apigw.yandexcloud.net/ws
```

Type messages and press Enter to send. All connected clients will receive the broadcast.

## WebSocket Events

| Event | Handler | Description |
|-------|---------|-------------|
| CONNECT | `handleConnect` | Stores connection ID in YDB |
| MESSAGE | `handleMessage` | Broadcasts message to all connections |
| DISCONNECT | `handleDisconnect` | Removes connection ID from YDB |

## Database Migrations

Migrations run automatically during `terraform apply`. To run manually:

```bash
export IAM_TOKEN=$(yc iam create-token)
export YDB_CONNECTION_STRING=$(terraform -chdir=tf output -raw migrate)
goose -dir migrations ydb "$YDB_CONNECTION_STRING&token=$IAM_TOKEN" up
```

## Environment Variables

The function uses:

- `YDB_CONNECTION_STRING`: YDB connection string (set automatically by Terraform)

## Notes

- Maximum connection lifetime: 60 minutes
- Idle timeout: 10 minutes (send pings to keep alive)
- Message size limit: 128 KB
- Frame size limit: 32 KB

## Cleanup

```bash
cd tf
terraform destroy
```

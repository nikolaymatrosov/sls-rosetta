# OpenTelemetry Tracing with Yandex Cloud Monium

This example demonstrates distributed tracing between two Yandex Cloud Functions using OpenTelemetry, with traces exported to [Yandex Cloud Monium](https://yandex.cloud/ru/docs/monitoring/).

## Architecture

```
HTTP Request → caller function → HTTP fetch → callee function
                  |                                |
                  └──── OTLP/gRPC traces ──────────┘
                              ↓
                  Monium (ingest.monium.yandex.cloud:443)
```

- **Caller** (`caller.ts`): Receives an HTTP request, creates a root span, calls the callee function with trace context propagation, and returns the combined result with trace ID.
- **Callee** (`callee.ts`): Receives the call, extracts trace context from headers, creates a child span, and returns a response.
- **Tracing** (`tracing.ts`): Shared module that initializes OpenTelemetry SDK with OTLP/gRPC exporter targeting Monium.

## Features

- Distributed tracing across two serverless functions
- W3C TraceContext propagation via HTTP headers
- IAM token-based authentication to Monium (from function context)
- Automatic span flush before function response

## Prerequisites

- Yandex Cloud account with Monium enabled
- Terraform >= 1.0
- Node.js >= 18
- Go (for running tests)

## Project Structure

```
otel/
├── README.md
├── function/
│   ├── caller.ts          # Caller function handler
│   ├── callee.ts          # Callee function handler
│   ├── tracing.ts         # Shared OpenTelemetry setup
│   ├── package.json
│   └── tsconfig.json
├── tf/
│   ├── terraform.tf       # Provider configuration
│   ├── variables.tf       # Input variables
│   ├── iam.tf             # Service account with monium.traces.writer
│   ├── main.tf            # Function resources and IAM bindings
│   └── outputs.tf         # Function IDs and URLs
└── environment/
    └── terraform.tfstate
```

## Deployment

### 1. Configure Variables

Create `tf/terraform.tfvars`:

```hcl
cloud_id  = "your-cloud-id"
folder_id = "your-folder-id"
```

### 2. Deploy

```bash
cd tf
terraform init
terraform apply
```

This will:
1. Build the TypeScript code (`npm ci && npm run build`)
2. Package the compiled code with dependencies
3. Deploy two functions with a shared service account
4. Make both functions publicly accessible

## Usage

Invoke the caller function:

```bash
curl https://functions.yandexcloud.net/<caller_function_id>
```

Response:

```json
{
  "message": "Caller completed",
  "traceId": "abc123...",
  "calleeResponse": {
    "message": "Callee processed successfully",
    "traceId": "abc123..."
  }
}
```

Both functions share the same `traceId`, confirming distributed trace propagation.

## How It Works

1. The **caller** receives an HTTP request and initializes OpenTelemetry with the IAM token from `context.token.access_token`.
2. It creates a root span and makes an HTTP request to the **callee**, injecting W3C `traceparent` header.
3. The **callee** extracts the trace context from incoming headers, creating a child span under the same trace.
4. Both functions export spans to Monium via OTLP/gRPC before returning their responses.

### Monium Connection Details

| Parameter | Value |
|-----------|-------|
| Endpoint | `ingest.monium.yandex.cloud:443` |
| Protocol | OTLP/gRPC with TLS |
| Auth | `Authorization: Bearer <IAM_TOKEN>` |
| Project header | `x-monium-project: folder__<folder_id>` |
| Required role | `monium.traces.writer` |

## Testing

```bash
cd tests/typescript/otel
export CLOUD_ID=your-cloud-id
export FOLDER_ID=your-folder-id
export YC_TOKEN=your-token

go test -v -timeout 30m
```

## Cleanup

```bash
cd tf
terraform destroy
```

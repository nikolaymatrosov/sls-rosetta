# Yandex Data Streams (YDS) Example - TypeScript

This example demonstrates how to use Yandex Data Streams (YDS) with Yandex Cloud Functions in TypeScript. It implements a producer-consumer architecture where:

1. **Producer Function**: An HTTP-triggered function that accepts JSON messages and writes them to a YDS topic
2. **Consumer Function**: A YDS-triggered function that processes messages from the topic in batches
3. **YDB Database**: Serverless YDB database hosting the data stream
4. **YDS Topic**: A topic (data stream) for message queueing
5. **Trigger**: Automatically invokes the consumer function when messages arrive

## Architecture

```txt
HTTP Request → Producer Function → YDB Topic (Data Stream) → Trigger → Consumer Function
```

## Prerequisites

- Yandex Cloud account
- Terraform >= 1.0
- Node.js >= 20.19 (for local development)
- npm or yarn

## Project Structure

```txt
examples/typescript/yds/
├── README.md                    # This file
├── function/
│   ├── producer.ts             # Producer handler
│   ├── consumer.ts             # Consumer handler
│   ├── package.json            # Dependencies
│   └── tsconfig.json           # TypeScript config
├── dist/                       # Build output (auto-generated)
│   ├── producer.js
│   ├── consumer.js
│   └── package.json
├── tf/
│   ├── terraform.tf            # Provider and backend configuration
│   ├── variables.tf            # Input variables
│   ├── outputs.tf              # Output values
│   ├── yds.tf                  # YDB database and topic
│   ├── iam.tf                  # Service accounts and IAM
│   └── main.tf                 # Build, functions, and trigger
└── environment/
    └── terraform.tfstate       # Terraform state file
```

## Dependencies

This example uses the YDB JavaScript SDK packages:

- `@ydbjs/core` - Core driver and connection management
- `@ydbjs/auth` - Authentication (MetadataCredentialsProvider for cloud functions)
- `@ydbjs/topic` - Topic/streaming API for YDS

For type definitions:

- `@yandex-cloud/function-types` - TypeScript types for Yandex Cloud Functions events

## Deployment

### 1. Configure Variables

Create a `terraform.tfvars` file in the `tf/` directory:

```hcl
cloud_id  = "your-cloud-id"
folder_id = "your-folder-id"
zone      = "ru-central1-a"
```

### 2. Initialize Terraform

```bash
cd tf
terraform init
```

### 3. Deploy Infrastructure

```bash
terraform apply
```

This will:

1. Install npm dependencies
2. Build TypeScript code to JavaScript
3. Create YDB Serverless Database and Topic
4. Deploy Producer and Consumer functions
5. Create service accounts with appropriate IAM roles
6. Set up YDS trigger

### 4. Get Outputs

After deployment, Terraform will output important information:

```bash
terraform output
```

Key outputs:

- `producer_function_url`: URL to call the producer function
- `curl_example`: Ready-to-use curl command for testing

## Usage

### Send a Message

Use the curl command from the terraform output, or manually:

```bash
export PRODUCER_URL=$(terraform output -raw producer_function_url)
curl -sS -X POST "$PRODUCER_URL" \
  -H "Content-Type: application/json" \
  -d '{
    "message": "User logged in successfully",
    "user_id": "user123",
    "action": "login"
  }'
```

### Message Format

The producer expects the following JSON structure:

```json
{
  "message": "string",     // Message content
  "user_id": "string",     // User identifier
  "action": "string"       // Action type: login, logout, purchase, view, etc.
}
```

### View Logs

Check the producer function logs:

```bash
export PRODUCER_FUNCTION_ID=$(terraform output -raw producer_function_id)
yc serverless function logs $PRODUCER_FUNCTION_ID
```

Check the consumer function logs:

```bash
export CONSUMER_FUNCTION_ID=$(terraform output -raw consumer_function_id)
yc serverless function logs $CONSUMER_FUNCTION_ID --follow
```

## How It Works

### Producer Function

1. Receives HTTP POST requests with JSON payload
2. Validates required fields (`message`, `user_id`, `action`)
3. Adds timestamp to the message
4. Creates YDB driver with MetadataCredentialsProvider (automatic cloud authentication)
5. Uses `createTopicWriter()` from `@ydbjs/topic` to write messages
6. Returns success/error response

**Key Code:**

```typescript
const driver = new Driver(ydbEndpoint, databasePath, {
    credentials: new MetadataCredentialsProvider()
});

const writer = await createTopicWriter(driver, {
    topic: topicName,
    producerId: 'producer-ts',
});

writer.write({ data: Buffer.from(messageData, 'utf-8') });
await writer.flush();
```

### Consumer Function

1. Triggered automatically when messages arrive in the YDS topic
2. Processes messages in batches (up to 10 messages or 5-second cutoff)
3. Decodes base64-encoded message data
4. Groups messages by user
5. Processes different action types with custom logic:
   - `login`: User authentication events
   - `logout`: User session end events
   - `purchase`: Transaction events
   - `view`: Page view events
   - Custom actions as needed
6. Logs processing results and summaries

**Event Structure:**

```typescript
export const handler: Handler.DataStreams = async (event, context) => {
    for (const msg of event.messages) {
        const dataDecoded = Buffer.from(msg.details.message.data, 'base64').toString('utf-8');
        const dataJson = JSON.parse(dataDecoded);
        // Process message
    }
};
```

### YDS Trigger Configuration

- **Batch Size**: 10 messages
- **Batch Cutoff**: 5 seconds
- **Retry Attempts**: 3
- **Retry Interval**: 10 seconds

Messages are delivered to the consumer in batches, either when:

- 10 messages accumulate in the topic, OR
- 5 seconds have passed since the last batch

## Testing

### Test Multiple Messages

```bash
{
  # Populate PRODUCER_URL from terraform outputs (use -raw to avoid quotes)
  export PRODUCER_URL=$(terraform output -raw producer_function_url)
  echo "Producer URL set to: $PRODUCER_URL"

  # Send login event
  echo "Sending login event for user123..."
  curl -sS -X POST "$PRODUCER_URL" \
    -H "Content-Type: application/json" \
    -d '{"message": "User logged in", "user_id": "user123", "action": "login"}' \
    | jq . || echo "Non-JSON or empty response"

  # Send purchase event
  echo "Sending purchase event for user123..."
  curl -sS -X POST "$PRODUCER_URL" \
    -H "Content-Type: application/json" \
    -d '{"message": "Purchased item XYZ", "user_id": "user123", "action": "purchase"}' \
    | jq . || echo "Non-JSON or empty response"

  # Send view event
  echo "Sending view event for user456..."
  curl -sS -X POST "$PRODUCER_URL" \
    -H "Content-Type: application/json" \
    -d '{"message": "Viewed product ABC", "user_id": "user456", "action": "view"}' \
    | jq . || echo "Non-JSON or empty response"

  # Send logout event
  echo "Sending logout event for user123..."
  curl -sS -X POST "$PRODUCER_URL" \
    -H "Content-Type: application/json" \
    -d '{"message": "User logged out", "user_id": "user123", "action": "logout"}' \
    | jq . || echo "Non-JSON or empty response"
}
```

### Monitor Consumer Processing

The consumer function will:

1. Receive messages in batches
2. Group them by user
3. Log processing details for each message
4. Print a summary of all users affected

Check logs to see the processing results:

```bash
# Export consumer function id from terraform outputs (use -raw to avoid quotes)
export CONSUMER_FUNCTION_ID=$(terraform output -raw consumer_function_id)

# Follow consumer logs
yc serverless function logs $CONSUMER_FUNCTION_ID --follow
```

## Local Development

### Build the Project Locally

```bash
cd function
npm install
npm run build
```

This will compile TypeScript to JavaScript in the `../dist/` directory.

### Test TypeScript Code

You can add tests using your preferred testing framework:

```bash
npm install --save-dev jest @types/jest ts-jest
```

## Cost Optimization

The YDB database is configured as serverless with automatic scaling based on usage. The `sleep_after = 5` configuration helps minimize costs during inactivity.

## Customization

### Add Custom Action Types

Edit [function/consumer.ts](function/consumer.ts) and add new cases in the `processAction()` function:

```typescript
function processAction(userId: string, action: string, message: string, timestamp: string): void {
    switch (action) {
        case 'your_custom_action':
            console.log(`Processing custom action for ${userId}`);
            // Your custom logic here
            break;
        // ... other cases
    }
}
```

### Modify Batch Settings

Edit [tf/main.tf](tf/main.tf:86-87) to adjust the trigger batch configuration:

```hcl
data_streams {
  batch_cutoff = 5   # Seconds to wait before processing batch
  batch_size   = 10  # Maximum messages per batch
}
```

### Change Function Resources

Edit [tf/main.tf](tf/main.tf) to adjust memory and timeout:

```hcl
resource "yandex_function" "producer" {
  memory            = 256  # MB
  execution_timeout = "30" # Seconds
  # ...
}
```

### Update SDK Versions

Edit [function/package.json](function/package.json) to use newer versions:

```json
{
  "dependencies": {
    "@ydbjs/auth": "^6.0.0",
    "@ydbjs/core": "^6.0.0",
    "@ydbjs/topic": "^6.0.0"
  }
}
```

## Cleanup

To destroy all created resources:

```bash
cd tf
terraform destroy
```

This will remove:

- All functions
- YDB database and topic
- Service accounts
- IAM bindings
- Triggers

## Troubleshooting

### Producer Returns 500 Error

Check:

1. Service account has `ydb.editor` role
2. YDB endpoint and topic path are correctly configured
3. Producer function logs for detailed error messages
4. Node.js runtime is nodejs20 (required for @ydbjs/topic)

### Consumer Not Processing Messages

Check:

1. Trigger is created and active
2. Trigger service account has `ydb.admin` role
3. Consumer function has `ydb.viewer` role
4. Topic has messages (check YDB console)

### Build Failures

Check:

1. Node.js >= 20.19 is installed
2. TypeScript compilation completes without errors: `cd function && npm run build`
3. All dependencies are installed: `npm install`

### IAM Permission Issues

Wait 5-10 seconds after deployment for IAM bindings to propagate. The Terraform configuration includes `sleep_after = 5` for this reason.

## SDK Documentation

For more details on the YDB JavaScript SDK:

- [Core Documentation](docs/core.md)
- [Authentication Documentation](docs/auth.md)
- [Topic API Documentation](docs/topic.md)

## Related Examples

- [Go YDS Example](../../go/yds/README.md) - Same functionality in Go
- [Python YDS Example](../../python/yds/README.md) - Same functionality in Python

## Resources

- [Yandex Data Streams Documentation](https://cloud.yandex.com/en/docs/data-streams/)
- [YDB Documentation](https://ydb.tech/en/docs/)
- [Yandex Cloud Functions Documentation](https://cloud.yandex.com/en/docs/functions/)
- [YDB JavaScript SDK](https://github.com/ydb-platform/ydb-nodejs-sdk)

# Go Yandex Data Streams (YDS) Example

This example demonstrates how to use Yandex Data Streams (YDS) with Go functions in Yandex Cloud Functions. The example includes two functions: a producer that writes data to a YDS stream and a consumer that is triggered by stream events.

## Architecture

```
HTTP Request → Producer Function → YDB Topic → Trigger → Consumer Function
```

### Components

1. **Producer Function** (`ProducerHandler`): HTTP endpoint that accepts JSON data and writes it to a YDB topic
2. **YDB Topic**: Data topic that stores and forwards messages
3. **Trigger**: Automatically invokes the consumer function when new data arrives in the topic
4. **Consumer Function** (`ConsumerHandler`): Processes batch messages from the topic

## Features

- **Two-Function Architecture**: Separate producer and consumer functions
- **Batch Processing**: Consumer processes multiple messages in batches
- **Event-Driven**: Automatic triggering based on stream events
- **IAM Integration**: Proper service accounts and permissions
- **Error Handling**: Comprehensive error handling and logging
- **Scalable**: Designed to handle high-throughput data processing

## Prerequisites

- Yandex Cloud CLI configured
- Terraform installed
- Go 1.21+ for local development

## Deployment

1. Navigate to the `tf` directory:
   ```bash
   cd examples/go/yds/tf
   ```

2. Initialize Terraform:
   ```bash
   terraform init
   ```

3. Create a `terraform.tfvars` file with your Yandex Cloud credentials:
   ```hcl
   cloud_id  = "your-cloud-id"
   folder_id = "your-folder-id"
   zone      = "ru-central1-a"
   ```

4. Deploy the infrastructure:
   ```bash
   terraform apply
   ```

## Usage

### 1. Send Data to the Topic

Use the producer function to send data to the YDB topic:

```bash
# Get the producer function URL
PRODUCER_URL=$(terraform output -raw producer_function_url)

# Send a test message
curl -X POST $PRODUCER_URL \
  -H "Content-Type: application/json" \
  -d '{
    "message": "Hello from YDS!",
    "user_id": "user123",
    "action": "login"
  }'
```

### 2. Monitor Consumer Function

The consumer function will automatically be triggered when data arrives in the topic. You can monitor the function logs:

```bash
# Get the consumer function ID
CONSUMER_ID=$(terraform output -raw consumer_function_id)

# View function logs (using Yandex Cloud CLI)
yc serverless function logs $CONSUMER_ID
```

### 3. Test Different Event Types

Try different event types to see how the consumer processes them:

```bash
# Login event
curl -X POST $PRODUCER_URL \
  -H "Content-Type: application/json" \
  -d '{"message": "User logged in", "user_id": "user123", "action": "login"}'

# Purchase event
curl -X POST $PRODUCER_URL \
  -H "Content-Type: application/json" \
  -d '{"message": "Product purchased", "user_id": "user456", "action": "purchase"}'

# View event
curl -X POST $PRODUCER_URL \
  -H "Content-Type: application/json" \
  -d '{"message": "Page viewed", "user_id": "user789", "action": "view"}'
```

## Infrastructure

The Terraform configuration creates:

- **YDB Database**: Serverless YDB database for Data Streams
- **YDB Topic**: Data topic for message ingestion
- **Producer Function**: HTTP function that writes to the topic
- **Consumer Function**: Triggered function that processes topic data
- **YDS Trigger**: Links the topic to the consumer function
- **Service Accounts**: With appropriate YDB and Functions permissions
- **IAM Bindings**: Makes the producer function publicly accessible

## Function Details

### Producer Function (`ProducerHandler`)

- **Entry Point**: `main.ProducerHandler`
- **Trigger**: HTTP requests
- **Purpose**: Accepts JSON data and writes it to YDB topic
- **Environment Variables**:
  - `YDB_ENDPOINT`: YDB endpoint URL (e.g., grpcs://ydb.serverless.yandexcloud.net:2135)
  - `YDS_TOPIC_ID`: Name of the YDB topic

### Consumer Function (`ConsumerHandler`)

- **Entry Point**: `main.ConsumerHandler`
- **Trigger**: YDS stream events
- **Purpose**: Processes batch messages from the topic
- **Features**:
  - Batch processing of multiple messages
  - Event type-based processing logic
  - Error handling and logging
  - User-based event grouping

## Event Processing

The consumer function processes different event types:

- **login**: User login events
- **logout**: User logout events
- **purchase**: Purchase events
- **view**: Page/view events
- **default**: Unknown event types

## Configuration

### Topic Configuration

- **Retention Period**: 24 hours
- **Partitions**: 1
- **Database**: Serverless YDB

### Trigger Configuration

- **Batch Size**: 10 messages
- **Batch Cutoff**: 5 seconds
- **Service Account**: Dedicated trigger service account

### Function Configuration

- **Runtime**: Go 1.21
- **Memory**: 128 MB
- **Timeout**: 10 seconds

## Environment Variables

### Producer Function
- `YDB_ENDPOINT`: YDB endpoint URL (e.g., grpcs://ydb.serverless.yandexcloud.net:2135)
- `YDS_TOPIC_ID`: Name of the YDB topic

### Consumer Function
- No environment variables required (receives data via trigger)

## Local Development

To run the functions locally:

1. Set environment variables:
   ```bash
   export YDB_ENDPOINT="your-ydb-endpoint"
   export YDS_TOPIC_ID="your-topic-name"
   ```

2. Run the producer function:
   ```bash
   go run function/producer.go
   ```

3. Test the consumer function with sample data:
   ```bash
   go run function/consumer.go
   ```

## Monitoring and Logging

### Function Logs

Both functions include comprehensive logging:

- Request/response logging
- Error handling and reporting
- Event processing details
- Performance metrics

### Cloud Monitoring

Monitor the functions using Yandex Cloud Monitoring:

- Function invocation metrics
- Error rates
- Execution times
- Topic throughput

## Cleanup

To destroy the infrastructure:

```bash
terraform destroy
```

## Troubleshooting

### Common Issues

1. **Permission Errors**: Ensure service accounts have proper YDS permissions
2. **Topic Not Found**: Verify the topic ID in environment variables
3. **Function Not Triggered**: Check trigger configuration and IAM permissions
4. **Data Not Processing**: Verify consumer function logs for errors

### Debug Mode

Enable verbose logging by setting the log level:

```bash
export LOG_LEVEL=debug
```

## Security Considerations

- Service accounts have minimal required permissions
- Producer function is publicly accessible (can be restricted)
- Consumer function is only accessible via trigger
- All communication uses secure protocols

## Performance Optimization

- Batch processing reduces function invocations
- Configurable batch size and cutoff times
- Efficient JSON parsing and processing
- Minimal memory footprint

## Outputs

After deployment, Terraform will output:

- `producer_function_url`: URL to invoke the producer function
- `consumer_function_id`: ID of the consumer function
- `yds_topic_id`: ID of the YDB topic
- `yds_topic_name`: Name of the YDB topic
- `yds_database_id`: ID of the YDB database
- `yds_database_path`: Database path for YDB connections

## Next Steps

This example can be extended with:

- Database integration (YDB, PostgreSQL)
- Message queue integration (YMQ)
- Object storage integration
- Advanced analytics and monitoring
- Custom event schemas 
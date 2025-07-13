# YDS TypeScript Functions

This folder contains TypeScript implementations for YDS topic producer and consumer functions.

## Producer
- **File:** `producer.ts`
- Accepts POST requests with JSON body `{ message, user_id, action }`
- Writes a JSON event to the YDS topic using @ydbjs/topic
- Environment variables:
  - `YDB_CONNECTION_STRING`: YDB connection string
  - `YDS_TOPIC_ID`: YDS topic name

## Consumer
- **File:** `consumer.ts`
- Accepts POST requests with YDS event format (array of messages)
- Logs and processes each message

## Shared Types
- **File:** `types.ts`
- Contains TypeScript interfaces for request, response, and event types

## Local Testing
You can run the producer locally with:
```
node producer.ts
```

## Deployment
See the `../tf/` folder for Terraform deployment examples. 
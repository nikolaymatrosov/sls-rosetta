import { MetadataCredentialsProvider } from "@ydbjs/auth/metadata";
import { Driver } from "@ydbjs/core";
import { query } from "@ydbjs/query";
import { topic } from "@ydbjs/topic";
import type {
  ConnectEvent,
  MessageEvent,
  DisconnectEvent,
  Context,
  HttpResult,
} from "./types/index.js";
import type { ServerMessage, ClientMessage } from "../shared/protocol.js";
import {
  ClientMessageType,
  createBroadcastMessage,
  createConnectedMessage,
  createUserJoinedMessage,
  createUserLeftMessage,
  createErrorMessage,
  createAckMessage,
  parseClientMessage,
} from "../shared/protocol.js";
import {
  storeConnection,
  removeConnectionById,
  getUserIdByConnectionId,
} from "./database.js";

// YDB Topic configuration for broadcasting
const BROADCAST_TOPIC = process.env.BROADCAST_TOPIC || "/Root/broadcast-topic";
const TOPIC_PRODUCER_ID = "websocket-broadcaster";

// Union type for all WebSocket events
type WebSocketEvent = ConnectEvent | MessageEvent | DisconnectEvent;

// Helper functions for typed HTTP responses
function successResponse(message: string): HttpResult {
  return { statusCode: 200, body: JSON.stringify(createAckMessage(message)) };
}

function errorResponse(statusCode: number, error: string): HttpResult {
  return { statusCode, body: JSON.stringify(createErrorMessage(error)) };
}

// Write any server message to YDB Topic for broadcasting via trigger
async function writeToTopic(
  driver: Driver,
  message: ServerMessage
): Promise<void> {
  const t = topic(driver);

  // Use await using for automatic cleanup
  await using writer = t.createWriter({
    topic: BROADCAST_TOPIC,
    producer: TOPIC_PRODUCER_ID,
  });

  const messageJson = JSON.stringify(message);
  writer.write(new TextEncoder().encode(messageJson));
  await writer.flush();

  console.log(`Written to topic ${BROADCAST_TOPIC}: ${messageJson}`);
}

// Handle CONNECT event - extract user_id from query params and store
async function handleConnect(
  event: ConnectEvent,
  context: Context,
  driver: Driver
): Promise<HttpResult> {
  const sql = query(driver);
  const { connectionId, connectedAt } = event.requestContext;
  const userId = event.queryStringParameters?.user_id;

  if (!userId) {
    console.error("No user_id in query params");
    return errorResponse(400, "Missing user_id query parameter");
  }

  console.log(`User ${userId} connected: ${connectionId}`);
  await storeConnection(sql, userId, connectionId, new Date(connectedAt));

  // Write USER_JOINED message to topic so trigger broadcasts it
  const joinedMsg = createUserJoinedMessage(userId);
  await writeToTopic(driver, joinedMsg);

  // Return CONNECTED message directly to the connecting user
  const message = createConnectedMessage(userId);
  return {
    statusCode: 200,
    body: JSON.stringify(message),
  };
}

// Handle MESSAGE event - dispatch based on protocol message type
async function handleMessage(
  event: MessageEvent,
  context: Context,
  driver: Driver
): Promise<HttpResult> {
  const sql = query(driver);
  const { connectionId } = event.requestContext;
  const body = event.body || "";

  // Get IAM token from context for API calls
  const iamToken = context.token?.access_token;
  if (!iamToken) {
    console.error("No IAM token available");
    return errorResponse(500, "No IAM token");
  }

  // Parse and validate client message
  let message: ClientMessage;
  try {
    message = parseClientMessage(body);
  } catch (error) {
    const errorMessage = error instanceof Error ? error.message : "Parse error";
    console.error(
      `Failed to parse message from ${connectionId}: ${errorMessage}`
    );
    return errorResponse(400, errorMessage);
  }

  // Handle messages based on type
  switch (message.type) {
    case ClientMessageType.SEND: {
      const senderId = await getUserIdByConnectionId(sql, connectionId);
      if (!senderId) {
        console.error(`Unregistered connection: ${connectionId}`);
        return errorResponse(400, "Not registered");
      }

      console.log(
        `Message from ${senderId} (${connectionId}): ${message.message}`
      );

      // Create broadcast message
      const broadcast = createBroadcastMessage(senderId, message.message);

      // Write to YDB Topic instead of broadcasting directly
      // The Data Streams trigger will handle actual broadcasting
      await writeToTopic(driver, broadcast);

      return successResponse("Message sent to topic");
    }

    case ClientMessageType.DISCONNECT: {
      console.log(`Graceful disconnect from ${connectionId}`);
      await removeConnectionById(sql, connectionId);
      return successResponse("Disconnected");
    }
  }
}

// Handle DISCONNECT event
async function handleDisconnect(
  event: DisconnectEvent,
  context: Context,
  driver: Driver
): Promise<HttpResult> {
  const sql = query(driver);
  const { connectionId } = event.requestContext;

  // Get userId before removing connection
  const userId = await getUserIdByConnectionId(sql, connectionId);

  console.log(`Disconnected: ${connectionId} (user: ${userId})`);

  await removeConnectionById(sql, connectionId);

  // Write USER_LEFT message to topic so trigger broadcasts it
  if (userId) {
    const leftMsg = createUserLeftMessage(userId);
    await writeToTopic(driver, leftMsg);
  }

  return successResponse("Disconnected");
}

// Main handler - routes based on event type
export const handler = async (
  event: WebSocketEvent,
  context: Context
): Promise<HttpResult> => {
  // Create driver inside handler (recommended for serverless)
  const credentialsProvider = new MetadataCredentialsProvider();
  const driver = new Driver(process.env.YDB_CONNECTION_STRING!, {
    credentialsProvider,
    "ydb.sdk.enable_discovery": false, // Improves cold start performance
  });

  try {
    await driver.ready();

    // Handle WebSocket events
    const eventType = event.requestContext.eventType;

    switch (eventType) {
      case "CONNECT":
        return await handleConnect(event as ConnectEvent, context, driver);
      case "MESSAGE":
        return await handleMessage(event as MessageEvent, context, driver);
      case "DISCONNECT":
        return await handleDisconnect(event as DisconnectEvent, context, driver);
      default:
        console.error(`Unknown event type: ${eventType}`);
        return errorResponse(400, `Unknown event type: ${eventType}`);
    }
  } finally {
    // Always close the driver
    driver.close();
  }
};

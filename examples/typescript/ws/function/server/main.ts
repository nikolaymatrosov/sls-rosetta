import { MetadataCredentialsProvider } from "@ydbjs/auth/metadata";
import { Driver } from "@ydbjs/core";
import { query } from "@ydbjs/query";
import type {
  ConnectEvent,
  MessageEvent,
  DisconnectEvent,
  Context,
  HttpResult,
} from "./types/index.js";
import type {
  ServerMessage,
  BroadcastMessage,
  ErrorMessage,
  ConnectedMessage,
  UserJoinedMessage,
  UserLeftMessage,
  AckMessage,
  ClientMessage,
} from "../shared/protocol.js";
import {
  isSendMessage,
  ServerMessageType,
  ClientMessageType,
} from "../shared/protocol.js";
import {
  storeConnection,
  removeConnection,
  removeConnectionById,
  getAllConnections,
  getUserIdByConnectionId,
} from "./database.js";

// Union type for all WebSocket events
type WebSocketEvent = ConnectEvent | MessageEvent | DisconnectEvent;

// Helper functions for typed HTTP responses
function successResponse(message: string): HttpResult {
  const body: AckMessage = { type: ServerMessageType.ACK, message };
  return { statusCode: 200, body: JSON.stringify(body) };
}

function errorResponse(statusCode: number, error: string): HttpResult {
  const body: ErrorMessage = { type: ServerMessageType.ERROR, message: error };
  return { statusCode, body: JSON.stringify(body) };
}

// Helper functions for creating protocol messages
function createUserLeftMessage(userId: string): UserLeftMessage {
  return { type: ServerMessageType.USER_LEFT, userId };
}

function createUserJoinedMessage(userId: string): UserJoinedMessage {
  return { type: ServerMessageType.USER_JOINED, userId };
}

function createBroadcastMessage(from: string, message: string): BroadcastMessage {
  return { type: ServerMessageType.BROADCAST, from, message };
}

function createConnectedMessage(userId: string): ConnectedMessage {
  return { type: ServerMessageType.CONNECTED, userId };
}

function createErrorMessage(message: string): ErrorMessage {
  return { type: ServerMessageType.ERROR, message };
}

// WebSocket Management API endpoint
const WS_API_ENDPOINT = "https://apigateway-connections.api.cloud.yandex.net";

// Send message to a specific connection via REST API
async function sendMessage(
  connectionId: string,
  message: string,
  iamToken: string
): Promise<boolean> {
  try {
    const url = `${WS_API_ENDPOINT}/apigateways/websocket/v1/connections/${connectionId}:send`;

    const response = await fetch(url, {
      method: "POST",
      headers: {
        Authorization: `Bearer ${iamToken}`,
        "Content-Type": "application/json",
      },
      body: JSON.stringify({
        data: Buffer.from(message).toString("base64"),
        type: "TEXT",
      }),
    });

    return response.ok;
  } catch (error) {
    console.error(`Failed to send to ${connectionId}:`, error);
    return false;
  }
}

// Broadcast message to all connections
async function broadcastMessage(
  sql: ReturnType<typeof query>,
  message: string,
  iamToken: string,
  excludeConnectionId?: string
): Promise<void> {
  const connections = await getAllConnections(sql);
  console.log(
    `Broadcasting to ${connections.length} connections, excluding: ${excludeConnectionId}`
  );
  console.log(`Connections: ${JSON.stringify(connections)}`);

  const filtered = connections.filter(
    (conn) => conn.connection_id !== excludeConnectionId
  );
  console.log(`After filter: ${filtered.length} connections`);

  const sendPromises = filtered.map(async (conn) => {
    const success = await sendMessage(conn.connection_id, message, iamToken);
    if (!success) {
      // Connection might be stale, remove it
      await removeConnection(sql, conn.user_id);
    }
  });

  await Promise.all(sendPromises);
}

// Handle CONNECT event - extract user_id from query params and store
async function handleConnect(
  event: ConnectEvent,
  context: Context,
  sql: ReturnType<typeof query>
): Promise<HttpResult> {
  const { connectionId, connectedAt } = event.requestContext;
  const userId = event.queryStringParameters?.user_id;

  if (!userId) {
    console.error("No user_id in query params");
    return errorResponse(400, "Missing user_id query parameter");
  }

  // Broadcast user_joined to existing users before storing new connection
  const iamToken = context.token?.access_token;
  if (iamToken) {
    const joinedMsg = createUserJoinedMessage(userId);
    await broadcastMessage(sql, JSON.stringify(joinedMsg), iamToken);
  }

  console.log(`User ${userId} connected: ${connectionId}`);
  await storeConnection(sql, userId, connectionId, new Date(connectedAt));

  const message = createConnectedMessage(userId);
  return {
    statusCode: 200,
    body: JSON.stringify(message),
  };
}

// Send a protocol message to a connection
async function sendProtocolMessage(
  connectionId: string,
  msg: ServerMessage,
  iamToken: string
): Promise<boolean> {
  return sendMessage(connectionId, JSON.stringify(msg), iamToken);
}

// Parse and validate client message from JSON string
function parseClientMessage(body: string): ClientMessage {
  // Parse JSON
  let parsed: unknown;
  try {
    parsed = JSON.parse(body);
  } catch (error) {
    throw new Error("Invalid JSON");
  }

  // Validate that parsed message has a type field
  if (!parsed || typeof parsed !== "object" || !("type" in parsed)) {
    throw new Error("Invalid message format");
  }

  // Validate message type
  const messageType = (parsed as { type: unknown }).type;
  if (
    messageType !== ClientMessageType.SEND &&
    messageType !== ClientMessageType.DISCONNECT
  ) {
    throw new Error("Unknown message type");
  }

  return parsed as ClientMessage;
}

// Handle MESSAGE event - dispatch based on protocol message type
async function handleMessage(
  event: MessageEvent,
  context: Context,
  sql: ReturnType<typeof query>
): Promise<HttpResult> {
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
    console.error(`Failed to parse message from ${connectionId}: ${errorMessage}`);
    return errorResponse(400, errorMessage);
  }

  // Handle messages based on type
  switch (message.type) {
    case ClientMessageType.SEND: {
      const senderId = await getUserIdByConnectionId(sql, connectionId);
      if (!senderId) {
        console.error(`Unregistered connection: ${connectionId}`);
        const errorMsg = createErrorMessage("Not registered. Send connect message first.");
        await sendProtocolMessage(connectionId, errorMsg, iamToken);
        return errorResponse(400, "Not registered");
      }

      console.log(
        `Message from ${senderId} (${connectionId}): ${message.message}`
      );

      const broadcast = createBroadcastMessage(senderId, message.message);
      await broadcastMessage(
        sql,
        JSON.stringify(broadcast),
        iamToken,
        connectionId
      );

      return successResponse("Broadcast sent");
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
  sql: ReturnType<typeof query>
): Promise<HttpResult> {
  const { connectionId } = event.requestContext;

  // Get userId before removing connection
  const userId = await getUserIdByConnectionId(sql, connectionId);

  console.log(`Disconnected: ${connectionId} (user: ${userId})`);

  await removeConnectionById(sql, connectionId);

  // Broadcast user_left to remaining users
  if (userId) {
    const iamToken = context.token?.access_token;
    if (iamToken) {
      const leftMsg = createUserLeftMessage(userId);
      await broadcastMessage(sql, JSON.stringify(leftMsg), iamToken);
    }
  }

  return successResponse("Disconnected");
}

// Main handler - routes based on event type
export const handler = async (
  event: WebSocketEvent,
  context: Context
): Promise<HttpResult> => {
  const eventType = event.requestContext.eventType;

  // Create driver inside handler (recommended for serverless)
  const credentialsProvider = new MetadataCredentialsProvider();
  const driver = new Driver(process.env.YDB_CONNECTION_STRING!, {
    credentialsProvider,
    "ydb.sdk.enable_discovery": false, // Improves cold start performance
  });

  try {
    await driver.ready();
    const sql = query(driver);

    switch (eventType) {
      case "CONNECT":
        return await handleConnect(event as ConnectEvent, context, sql);
      case "MESSAGE":
        return await handleMessage(event as MessageEvent, context, sql);
      case "DISCONNECT":
        return await handleDisconnect(event as DisconnectEvent, context, sql);
      default:
        console.error(`Unknown event type: ${eventType}`);
        return errorResponse(400, `Unknown event type: ${eventType}`);
    }
  } finally {
    // Always close the driver
    driver.close();
  }
};

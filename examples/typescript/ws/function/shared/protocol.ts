// Message type enums
export enum ClientMessageType {
  SEND = "send",
  DISCONNECT = "disconnect",
}

export enum ServerMessageType {
  BROADCAST = "broadcast",
  CONNECTED = "connected",
  USER_JOINED = "user_joined",
  USER_LEFT = "user_left",
  ERROR = "error",
  ACK = "ack",
}

// Client → Server messages
export type SendMessage = { type: ClientMessageType.SEND; message: string };
export type DisconnectMessage = { type: ClientMessageType.DISCONNECT };
export type ClientMessage = SendMessage | DisconnectMessage;

// Helper functions for creating client messages
export function createSendMessage(message: string): SendMessage {
  return { type: ClientMessageType.SEND, message };
}

export function createDisconnectMessage(): DisconnectMessage {
  return { type: ClientMessageType.DISCONNECT };
}

// Server → Client messages
export type BroadcastMessage = {
  type: ServerMessageType.BROADCAST;
  from: string;
  message: string;
};
export type ConnectedMessage = { type: ServerMessageType.CONNECTED; userId: string };
export type UserJoinedMessage = { type: ServerMessageType.USER_JOINED; userId: string };
export type UserLeftMessage = { type: ServerMessageType.USER_LEFT; userId: string };
export type ErrorMessage = { type: ServerMessageType.ERROR; message: string };
export type AckMessage = { type: ServerMessageType.ACK; message: string };
export type ServerMessage =
  | BroadcastMessage
  | ConnectedMessage
  | UserJoinedMessage
  | UserLeftMessage
  | ErrorMessage
  | AckMessage;

// Helper to check if an object has a type property
function hasType(msg: unknown): msg is { type: unknown } {
  return typeof msg === "object" && msg !== null && "type" in msg;
}

export function isSendMessage(msg: unknown): msg is SendMessage {
  return hasType(msg) && msg.type === ClientMessageType.SEND;
}

export function isDisconnectMessage(msg: unknown): msg is DisconnectMessage {
  return hasType(msg) && msg.type === ClientMessageType.DISCONNECT;
}

// Server message type guards
export function isBroadcastMessage(msg: unknown): msg is BroadcastMessage {
  return hasType(msg) && msg.type === ServerMessageType.BROADCAST;
}

export function isConnectedMessage(msg: unknown): msg is ConnectedMessage {
  return hasType(msg) && msg.type === ServerMessageType.CONNECTED;
}

export function isUserJoinedMessage(msg: unknown): msg is UserJoinedMessage {
  return hasType(msg) && msg.type === ServerMessageType.USER_JOINED;
}

export function isUserLeftMessage(msg: unknown): msg is UserLeftMessage {
  return hasType(msg) && msg.type === ServerMessageType.USER_LEFT;
}

export function isErrorMessage(msg: unknown): msg is ErrorMessage {
  return hasType(msg) && msg.type === ServerMessageType.ERROR;
}

export function isAckMessage(msg: unknown): msg is AckMessage {
  return hasType(msg) && msg.type === ServerMessageType.ACK;
}

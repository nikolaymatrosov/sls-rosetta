import type { RequestContext } from "./http.js";

type EventType = "MESSAGE" | "CONNECT" | "DISCONNECT";

interface WebSocketRequestContext<E extends EventType> extends RequestContext {
  connectedAt: number;
  eventType: E;
  connectionId: string;
}

interface Headers {
  "X-Yc-Apigateway-Websocket-Connected-At": string;
  "X-Yc-Apigateway-Websocket-Connection-Id": string;
  "X-Yc-Apigateway-Websocket-Event-Type": EventType;
  [key: string]: string;
}

interface ConnectHeaders extends Headers {
  "X-Yc-Apigateway-Websocket-Event-Type": "CONNECT";
}

interface MessageHeaders extends Headers {
  "X-Yc-Apigateway-Websocket-Event-Type": "MESSAGE";
  "X-Yc-Apigateway-Websocket-Message-Id": string;
}

interface DisconnectHeaders extends Headers {
  "X-Yc-Apigateway-Websocket-Event-Type": "DISCONNECT";
  "X-Yc-Apigateway-Websocket-Disconnect-Reason": string;
  "X-Yc-Apigateway-Websocket-Disconnect-Status-Code": string;
}

interface Common<H extends Headers> {
  resource: string;
  path: string;
  pathParameters: Record<string, string>;
  headers: H;
  multiValueHeaders: { [K in keyof H]: H[K][] };
  queryStringParameters: Record<string, string>;
  multiValueQueryStringParameters: Record<string, string[]>;
  parameters: Record<string, string>;
  multiValueParameters: Record<string, string[]>;
  operationId: string;
}

export interface ConnectEvent extends Common<ConnectHeaders> {
  httpMethod: string;
  requestContext: WebSocketRequestContext<"CONNECT">;
}

export interface MessageEvent extends Common<MessageHeaders> {
  body: string;
  isBase64Encoded: boolean;
  requestContext: WebSocketRequestContext<"MESSAGE"> & {
    messageId: string;
  };
}

export interface DisconnectEvent extends Common<DisconnectHeaders> {
  httpMethod: string;
  requestContext: WebSocketRequestContext<"DISCONNECT"> & {
    disconnectStatusCode: number;
  };
}

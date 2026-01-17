import type { WebSocket } from 'ws';
import { isDataStreamsTriggerMessage } from './protocol.js';

// Type definitions
type TypeGuard<T> = (msg: unknown) => msg is T;
type MessageHandler<T> = (message: T) => void;
type UnknownMessageHandler = (data: unknown) => void;
type ParseErrorHandler = (rawData: string, error: Error) => void;
type HandlerErrorCallback = (error: Error, message: unknown) => void;

interface HandlerEntry<T = unknown> {
  guard: TypeGuard<T>;
  handler: MessageHandler<T>;
}

/**
 * Type-safe message router for WebSocket messages.
 *
 * Provides a fluent API for registering message handlers with type guards.
 * Handles JSON parsing, type checking, and error handling automatically.
 *
 * @example
 * const router = new MessageRouter()
 *   .on(isBroadcastMessage, (msg) => console.log(msg.message))
 *   .on(isErrorMessage, (msg) => console.error(msg.message))
 *   .onUnknown((msg) => console.log('Unknown:', msg))
 *   .onParseError((raw) => console.log('Invalid JSON:', raw));
 *
 * ws.on('message', (data) => router.handle(data));
 */
export class MessageRouter {
  private handlers: HandlerEntry[] = [];
  private unknownHandler?: UnknownMessageHandler;
  private parseErrorHandler?: ParseErrorHandler;
  private handlerErrorCallback?: HandlerErrorCallback;

  /**
   * Register a typed message handler with its type guard.
   *
   * Handlers are checked in registration order. The first matching
   * handler will be executed. This is important for messages without
   * a "type" field (like AckMessage and ErrorResponse) which should
   * be registered after more specific typed messages.
   *
   * @param guard - Type guard function that narrows the message type
   * @param handler - Handler function that receives the typed message
   * @returns this router instance for method chaining
   *
   * @example
   * router.on(isBroadcastMessage, (msg) => {
   *   console.log(`${msg.from}: ${msg.message}`);
   * });
   */
  on<T>(guard: TypeGuard<T>, handler: MessageHandler<T>): this {
    this.handlers.push({ guard, handler } as HandlerEntry);
    return this;
  }

  /**
   * Register a fallback handler for messages that don't match any type guard.
   *
   * This handler receives the parsed JSON object when no registered
   * type guard returns true. Useful for logging unknown message types
   * or implementing fallback behavior.
   *
   * @param handler - Handler function that receives the unknown message
   * @returns this router instance for method chaining
   *
   * @example
   * router.onUnknown((msg) => {
   *   console.log('Unknown message:', msg);
   * });
   */
  onUnknown(handler: UnknownMessageHandler): this {
    this.unknownHandler = handler;
    return this;
  }

  /**
   * Register a handler for JSON parse errors.
   *
   * This handler is called when the incoming data cannot be parsed
   * as valid JSON. Receives the original raw string and the parse error.
   *
   * @param handler - Handler function that receives raw data and error
   * @returns this router instance for method chaining
   *
   * @example
   * router.onParseError((rawData, error) => {
   *   console.log('Invalid JSON:', rawData);
   * });
   */
  onParseError(handler: ParseErrorHandler): this {
    this.parseErrorHandler = handler;
    return this;
  }

  /**
   * Register a callback for errors thrown by message handlers.
   *
   * This callback is invoked when a registered handler throws an error
   * during execution. Prevents one bad handler from crashing the entire
   * message processing pipeline.
   *
   * @param callback - Callback function that receives error and message
   * @returns this router instance for method chaining
   *
   * @example
   * router.onHandlerError((error, msg) => {
   *   console.error('Handler error:', error.message);
   * });
   */
  onHandlerError(callback: HandlerErrorCallback): this {
    this.handlerErrorCallback = callback;
    return this;
  }

  /**
   * Process incoming message data.
   *
   * Handles JSON parsing, type checking, and handler dispatch.
   * Call this method from your WebSocket 'message' event handler.
   *
   * Error handling:
   * - Parse errors trigger parseErrorHandler
   * - Unknown message types trigger unknownHandler
   * - Handler errors trigger handlerErrorCallback
   *
   * @param data - WebSocket message data (string or buffer)
   *
   * @example
   * ws.on('message', (data) => router.handle(data));
   */
  handle(data: WebSocket.Data | string, isBinary: boolean): void {
    const rawData = data.toString();

    let parsed: unknown;
    try {
      parsed = JSON.parse(rawData);
    } catch (error) {
      if (this.parseErrorHandler) {
        this.parseErrorHandler(
          rawData,
          error instanceof Error ? error : new Error(String(error))
        );
      }
      return;
    }

    const handled = this.dispatchMessage(parsed);

    if (!handled && this.unknownHandler) {
      this.unknownHandler(parsed);
    }
  }

  /**
   * Dispatch a parsed message to the appropriate handler.
   *
   * If the message is a DataStreamsTriggerMessage, it unwraps the
   * messages array and dispatches each individual message separately.
   *
   * @param msg - Parsed message object
   * @returns true if a handler was found and executed, false otherwise
   */
  private dispatchMessage(msg: unknown): boolean {
    // Check if this is a Data Streams trigger message
    if (isDataStreamsTriggerMessage(msg)) {
      // Unwrap and dispatch each message in the array
      let handledAny = false;
      for (const serverMsg of msg.messages) {
        const handled = this.dispatchSingleMessage(serverMsg);
        handledAny = handledAny || handled;
      }
      return handledAny;
    }

    // Regular message, dispatch directly
    return this.dispatchSingleMessage(msg);
  }

  /**
   * Dispatch a single message to its handler.
   *
   * @param msg - Parsed message object
   * @returns true if a handler was found and executed, false otherwise
   */
  private dispatchSingleMessage(msg: unknown): boolean {
    for (const entry of this.handlers) {
      if (entry.guard(msg)) {
        try {
          (entry.handler as MessageHandler<unknown>)(msg);
          return true;
        } catch (error) {
          if (this.handlerErrorCallback) {
            this.handlerErrorCallback(
              error instanceof Error ? error : new Error(String(error)),
              msg
            );
          }
          return true;
        }
      }
    }
    return false;
  }
}

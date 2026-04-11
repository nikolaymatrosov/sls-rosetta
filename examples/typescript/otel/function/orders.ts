import {
  withTracing,
  getTracer,
  sqsClient,
  SendMessageCommand,
  injectTraceToSqsMessageAttributes,
  propagation,
  context as otelContext,
  SpanStatusCode,
  SpanKind,
} from "./tracing";

const calleeUrl = process.env.CALLEE_URL || "";
const ordersQueueUrl = process.env.ORDERS_QUEUE_URL || "";

export const handler = withTracing(async (span) => {
  const orderId = `order-${Date.now()}`;

  // Call payment service
  const calleeBody = await getTracer().startActiveSpan(
    "invoke-payment-service",
    { kind: SpanKind.CLIENT },
    async (clientSpan) => {
      try {
        const headers: Record<string, string> = {};
        propagation.inject(otelContext.active(), headers);

        const response = await fetch(calleeUrl, { method: "GET", headers });
        const body = await response.text();
        clientSpan.setAttribute("http.status_code", response.status);
        clientSpan.setStatus({ code: SpanStatusCode.OK });
        clientSpan.end();
        return body;
      } catch (err) {
        clientSpan.setStatus({
          code: SpanStatusCode.ERROR,
          message: String(err),
        });
        clientSpan.end();
        throw err;
      }
    },
  );

  // Publish order event to YMQ for inventory service — AWS SDK call is auto-instrumented
  await sqsClient.send(
    new SendMessageCommand({
      QueueUrl: ordersQueueUrl,
      MessageBody: JSON.stringify({ orderId, status: "confirmed" }),
      // MessageAttributes: injectTraceToSqsMessageAttributes(),
    }),
  );

  span.setAttribute("order.id", orderId);

  return {
    statusCode: 200,
    body: JSON.stringify({
      message: "Order placed",
      orderId,
      traceId: span.spanContext().traceId,
      paymentResponse: JSON.parse(calleeBody),
    }),
    headers: { "Content-Type": "application/json" },
  };
});

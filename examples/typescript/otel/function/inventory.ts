import {
  withMqTracing,
  getTracer,
  extractContextFromSqsMessageAttributes,
  context as otelContext,
  SpanStatusCode,
  SpanKind,
} from "./tracing";

export const handler = withMqTracing(async (_span, messages) => {
  const tracer = getTracer();
  for (const message of messages) {
    console.log("Received message:", JSON.stringify(message));
    const body = JSON.parse(message.details.message.body) as {
      orderId: string;
      status: string;
    };

    const parentContext = extractContextFromSqsMessageAttributes(
      message.details.message.message_attributes,
    );

    await otelContext.with(parentContext, () =>
      tracer.startActiveSpan(
        "reserve-inventory",
        { kind: SpanKind.CONSUMER },
        async (span) => {
          span.setAttribute("order.id", body.orderId);
          span.setAttribute("order.status", body.status);
          try {
            // Simulate inventory reservation
            await new Promise((resolve) => setTimeout(resolve, 20));
            span.addEvent("inventory-reserved");
            span.setStatus({ code: SpanStatusCode.OK });
          } catch (err) {
            span.setStatus({
              code: SpanStatusCode.ERROR,
              message: String(err),
            });
            throw err;
          } finally {
            span.end();
          }
        },
      ),
    );
  }
});

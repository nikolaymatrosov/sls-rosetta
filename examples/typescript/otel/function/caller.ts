import {
  withTracing,
  propagation,
  context as otelContext,
  SpanStatusCode,
  SpanKind,
} from "./tracing";

const calleeUrl = process.env.CALLEE_URL || "";

export const handler = withTracing("order-service", async (tracer, span) => {
  const calleeBody = await tracer.startActiveSpan(
    "invoke-callee",
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
        clientSpan.setStatus({ code: SpanStatusCode.ERROR, message: String(err) });
        clientSpan.end();
        throw err;
      }
    },
  );

  span.setAttribute("callee.status", 200);

  return {
    statusCode: 200,
    body: JSON.stringify({
      message: "Caller completed",
      traceId: span.spanContext().traceId,
      calleeResponse: JSON.parse(calleeBody),
    }),
    headers: { "Content-Type": "application/json" },
  };
});

import { NodeTracerProvider } from "@opentelemetry/sdk-trace-node";
import { BatchSpanProcessor } from "@opentelemetry/sdk-trace-base";
import { OTLPTraceExporter } from "@opentelemetry/exporter-trace-otlp-grpc";
import { Resource } from "@opentelemetry/resources";
import { ATTR_SERVICE_NAME } from "@opentelemetry/semantic-conventions";
import {
  trace,
  context,
  propagation,
  Tracer,
  Span,
  SpanStatusCode,
  SpanKind,
} from "@opentelemetry/api";
import { W3CTraceContextPropagator } from "@opentelemetry/core";
import { registerInstrumentations } from "@opentelemetry/instrumentation";
import { AwsInstrumentation } from "@opentelemetry/instrumentation-aws-sdk";
import * as grpc from "@grpc/grpc-js";
import { Http } from "@yandex-cloud/function-types/dist/src/http";
import Context from "@yandex-cloud/function-types/dist/src/context";

export { trace, context, propagation, SpanStatusCode, SpanKind };

// Construct the provider and register instrumentations at module load time,
// before any AWS SDK imports, so monkey-patching takes effect in time.
const serviceName = process.env.OTEL_SERVICE_NAME || "unknown-service";
const folderId = process.env.FOLDER_ID || "";
const apiKey = process.env.MONIUM_API_KEY || "";

const metadata = new grpc.Metadata();
metadata.set("authorization", `Api-Key ${apiKey}`);
metadata.set("x-monium-project", `folder__${folderId}`);

const exporter = new OTLPTraceExporter({
  url: "https://ingest.monium.yandex.cloud:443",
  metadata,
});

const provider = new NodeTracerProvider({
  resource: new Resource({
    [ATTR_SERVICE_NAME]: serviceName,
    cluster: "default",
  }),
  spanProcessors: [new BatchSpanProcessor(exporter)],
});
provider.register({ propagator: new W3CTraceContextPropagator() });

registerInstrumentations({
  instrumentations: [new AwsInstrumentation()],
});

// Require AWS SDK *after* instrumentations are registered so the
// monkey-patching is in place before any client is created.
// Using require() instead of import because TypeScript hoists imports
// to the top of the compiled file, which would defeat the ordering.
// eslint-disable-next-line @typescript-eslint/no-var-requires
const {
  SQSClient,
  SendMessageCommand: SendMessageCommandCtor,
} = require("@aws-sdk/client-sqs") as typeof import("@aws-sdk/client-sqs");

export const sqsClient = new SQSClient({
  region: "ru-central1",
  endpoint: "https://message-queue.api.cloud.yandex.net",
});

export { SendMessageCommandCtor as SendMessageCommand };

export function getTracer(): Tracer {
  return trace.getTracer(serviceName);
}

export async function flushTracing(): Promise<void> {
  await provider.forceFlush();
}

type SqsMessageAttribute = { DataType: string; StringValue: string };

export function injectTraceToSqsMessageAttributes(
  messageAttributes: Record<string, SqsMessageAttribute> = {},
): Record<string, SqsMessageAttribute> {
  const carrier: Record<string, string> = {};
  propagation.inject(context.active(), carrier);

  for (const [key, value] of Object.entries(carrier)) {
    messageAttributes[key] = { DataType: "String", StringValue: value };
  }

  return messageAttributes;
}

export function extractContextFromSqsMessageAttributes(
  messageAttributes: Record<string, { string_value?: string }> = {},
) {
  const carrier: Record<string, string> = {};
  for (const key of ["traceparent", "tracestate", "baggage"]) {
    const attr = messageAttributes[key];
    if (attr && typeof attr.string_value === "string") {
      carrier[key] = attr.string_value;
    }
  }
  return propagation.extract(context.active(), carrier);
}

type TracedHttpHandler = (
  event: Http.Event,
  context: Context,
) => Promise<Http.Result>;

export function withTracing(
  fn: (span: Span, event: Http.Event) => Promise<Http.Result>,
): TracedHttpHandler {
  return async (event) => {
    const tracer = getTracer();

    const carrier: Record<string, string> = {};
    if (event.headers) {
      for (const [key, value] of Object.entries(event.headers)) {
        carrier[key.toLowerCase()] = value;
      }
    }
    const extractedContext = propagation.extract(context.active(), carrier);

    try {
      return await context.with(extractedContext, () =>
        tracer.startActiveSpan(
          "handle-request",
          { kind: SpanKind.SERVER },
          async (span: Span) => {
            span.setAttribute("http.method", event.httpMethod || "GET");
            try {
              const result = await fn(span, event);
              span.setStatus({ code: SpanStatusCode.OK });
              span.end();
              return result;
            } catch (err) {
              span.setStatus({
                code: SpanStatusCode.ERROR,
                message: String(err),
              });
              span.end();
              throw err;
            }
          },
        ),
      );
    } catch (err) {
      return {
        statusCode: 500,
        body: JSON.stringify({ error: String(err) }),
        headers: { "Content-Type": "application/json" },
      };
    } finally {
      await flushTracing();
    }
  };
}

type Metadata = {
  created_at: string;
  event_id: string;
  event_type: string;
  tracing_context: unknown;
  cloud_id: string;
  folder_id: string;
};
type MessageAttributeValue = {
  data_type: string;
  string_value: string;
};
type Details = {
  queue_id: string;
  message: {
    message_id: string;
    md5_of_body: string;
    body: string;
    attributes: {
      ApproximateFirstReceiveTimestamp: string;
      ApproximateReceiveCount: string;
      SenderId: string;
      SentTimestamp: string;
    };
    message_attributes: Record<string, MessageAttributeValue>;
    md5_of_message_attributes: string;
  };
};
type Message = {
  details: Details;
  event_metadata: Metadata;
};

type TracedMqHandler = (
  event: any,
  context: Context,
) => Promise<{ statusCode: number }>;

export function withMqTracing(
  fn: (span: Span, messages: Message[]) => Promise<void>,
): TracedMqHandler {
  return async (event) => {
    const tracer = getTracer();

    try {
      return await tracer.startActiveSpan(
        "process-messages",
        { kind: SpanKind.CONSUMER },
        async (span: Span) => {
          span.setAttribute("messaging.system", "aws_sqs");
          span.setAttribute(
            "messaging.batch.message_count",
            event.messages.length,
          );
          try {
            await fn(span, event.messages);
            span.setStatus({ code: SpanStatusCode.OK });
            span.end();
            return { statusCode: 200 };
          } catch (err) {
            span.setStatus({
              code: SpanStatusCode.ERROR,
              message: String(err),
            });
            span.end();
            throw err;
          }
        },
      );
    } catch (err) {
      return { statusCode: 500 };
    } finally {
      await flushTracing();
    }
  };
}

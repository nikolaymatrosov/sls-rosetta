import { NodeSDK } from '@opentelemetry/sdk-node';
import { OTLPTraceExporter } from '@opentelemetry/exporter-trace-otlp-grpc';
import { Resource } from '@opentelemetry/resources';
import { ATTR_SERVICE_NAME } from '@opentelemetry/semantic-conventions';
import { trace, context, propagation, Tracer, Span, SpanStatusCode, SpanKind } from '@opentelemetry/api';
import { W3CTraceContextPropagator } from '@opentelemetry/core';
import * as grpc from '@grpc/grpc-js';
import { Http } from '@yandex-cloud/function-types/dist/src/http';
import Context from '@yandex-cloud/function-types/dist/src/context';

export { trace, context, propagation, SpanStatusCode, SpanKind };

let sdk: NodeSDK | null = null;

export function initTracing(serviceName: string, folderId: string): Tracer {
    const apiKey = process.env.MONIUM_API_KEY || '';

    propagation.setGlobalPropagator(new W3CTraceContextPropagator());

    const metadata = new grpc.Metadata();
    metadata.set('authorization', `Api-Key ${apiKey}`);
    metadata.set('x-monium-project', `folder__${folderId}`);

    const exporter = new OTLPTraceExporter({
        url: 'https://ingest.monium.yandex.cloud:443',
        metadata,
    });

    sdk = new NodeSDK({
        resource: new Resource({
            [ATTR_SERVICE_NAME]: serviceName,
            'cluster': 'default',
        }),
        traceExporter: exporter,
    });

    sdk.start();

    return trace.getTracer(serviceName);
}

export async function shutdownTracing(): Promise<void> {
    if (sdk) {
        await sdk.shutdown();
        sdk = null;
    }
}

type TracedHandler = (
    event: Http.Event,
    context: Context & { functionFolderId: string },
) => Promise<Http.Result>;

export function withTracing(serviceName: string, fn: (tracer: Tracer, span: Span, event: Http.Event) => Promise<Http.Result>): TracedHandler {
    return async (event, fnContext) => {
        const tracer = initTracing(
            process.env.OTEL_SERVICE_NAME || serviceName,
            fnContext.functionFolderId,
        );

        const carrier: Record<string, string> = {};
        if (event.headers) {
            for (const [key, value] of Object.entries(event.headers)) {
                carrier[key.toLowerCase()] = value;
            }
        }
        const extractedContext = propagation.extract(context.active(), carrier);

        try {
            return await context.with(extractedContext, () =>
                tracer.startActiveSpan('handle-request', { kind: SpanKind.SERVER }, async (span) => {
                    span.setAttribute('http.method', event.httpMethod || 'GET');
                    try {
                        const result = await fn(tracer, span, event);
                        span.setStatus({ code: SpanStatusCode.OK });
                        span.end();
                        return result;
                    } catch (err) {
                        span.setStatus({ code: SpanStatusCode.ERROR, message: String(err) });
                        span.end();
                        throw err;
                    }
                }),
            );
        } catch (err) {
            return {
                statusCode: 500,
                body: JSON.stringify({ error: String(err) }),
                headers: { 'Content-Type': 'application/json' },
            };
        } finally {
            await shutdownTracing();
        }
    };
}

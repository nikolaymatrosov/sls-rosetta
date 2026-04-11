import { getTracer, flushTracing, propagation, context as otelContext, SpanStatusCode, SpanKind } from './tracing';

const orderUrl = process.env.ORDER_URL || '';

async function main() {
    const tracer = getTracer();

    await tracer.startActiveSpan('call-order-function', { kind: SpanKind.CLIENT }, async (span) => {
        try {
            const headers: Record<string, string> = {};
            propagation.inject(otelContext.active(), headers);

            const response = await fetch(orderUrl, { method: 'GET', headers });
            const body = await response.text();

            span.setAttribute('http.status_code', response.status);
            span.setStatus({ code: SpanStatusCode.OK });
            span.end();

            console.log('Response:', JSON.parse(body));
        } catch (err) {
            span.setStatus({ code: SpanStatusCode.ERROR, message: String(err) });
            span.end();
            throw err;
        }
    });

    await flushTracing();
}

main().catch((err) => {
    console.error(err);
    process.exit(1);
});

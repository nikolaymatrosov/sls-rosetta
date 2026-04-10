import { initTracing, shutdownTracing, propagation, context as otelContext, SpanStatusCode, SpanKind } from './tracing';

const callerUrl = process.env.CALLER_URL || '';
const folderId = process.env.FOLDER_ID || '';

async function main() {
    const tracer = initTracing('order-client', folderId);

    await tracer.startActiveSpan('call-caller-function', { kind: SpanKind.CLIENT }, async (span) => {
        try {
            const headers: Record<string, string> = {};
            propagation.inject(otelContext.active(), headers);

            const response = await fetch(callerUrl, { method: 'GET', headers });
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

    await shutdownTracing();
}

main().catch((err) => {
    console.error(err);
    process.exit(1);
});

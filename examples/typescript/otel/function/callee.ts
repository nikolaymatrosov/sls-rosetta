import { withTracing } from './tracing';

export const handler = withTracing('payment-service', async (_tracer, span) => {
    // Simulate some work
    await new Promise((resolve) => setTimeout(resolve, 50));
    span.addEvent('processing-complete');

    return {
        statusCode: 200,
        body: JSON.stringify({
            message: 'Callee processed successfully',
            traceId: span.spanContext().traceId,
        }),
        headers: { 'Content-Type': 'application/json' },
    };
});

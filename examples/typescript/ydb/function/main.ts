import {handler} from './handler';

try {
    const result = await handler({
        queryStringParameters: { userId: "e061e6a1-65eb-47f2-b71e-0839ec3d3c72" },
        httpMethod: 'GET',
        headers: {},
        multiValueHeaders: {},
        multiValueQueryStringParameters: {},
        requestContext: {
            identity: {
                sourceIp: '',
                userAgent: ''
            },
            httpMethod: 'GET',
            requestId: 'test-request-id',
            requestTime: '01/Jan/2024:00:00:00 +0000',
            requestTimeEpoch: Date.now()
        },
        body: '',
        isBase64Encoded: false
    });
    console.log("Function result:", result);
} catch (err) {
    console.error("Error executing function:", err);
}
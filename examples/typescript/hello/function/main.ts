import {Http} from '@yandex-cloud/function-types/dist/src/http';

// The handler function is an asynchronous function that handles HTTP events.
// It takes a Http.Event object as a parameter and returns a Promise that resolves to a Http.Result object.
export async function handler(event: Http.Event): Promise<Http.Result> {

    // Return a Http.Result object with a status code of 200, custom headers, and a body.
    return {
        statusCode: 200,
        headers: {
            'X-Custom-Header': 'Test',
        },
        body: `Hello, ${event.queryStringParameters.name}!`
    };

}
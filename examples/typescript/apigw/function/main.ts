import {Http} from '@yandex-cloud/function-types/dist/src/http';

// RequestBody interface represents the structure of the request body.
interface RequestBody {
    name: string; // The name property of the request body.
}

// The handler function is an asynchronous function that handles HTTP events.
// It takes a Http.Event object as a parameter and returns a Promise that resolves to a Http.Result object.
export async function handler(event: Http.Event): Promise<Http.Result> {

    // Parse the body of the event into a RequestBody object.
    const body = JSON.parse(event.body || '{}') as RequestBody;

    // Create a response object with a message property.
    const response = {
        message: `Hello, ${body.name || 'World'}!`
    }

    // Return a Http.Result object with a status code of 200, a 'Content-Type' header set to 'application/json',
    // and the body set to the JSON stringified version of the response object.
    return {
        statusCode: 200,
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(response)
    };

}
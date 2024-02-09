import {Http} from '@yandex-cloud/function-types/dist/src/http';
import {DynamoDBClient} from "@aws-sdk/client-dynamodb";
import {DynamoDBDocumentClient, PutCommand, QueryCommand} from "@aws-sdk/lib-dynamodb";

const ENDPOINT = process.env.ENDPOINT;

// The handler function is an asynchronous function that handles HTTP events.
// It takes a Http.Event object as a parameter and returns a Promise that resolves to a Http.Result object.
export async function handler(event: Http.Event): Promise<Http.Result> {

    const client = new DynamoDBClient({
        endpoint: ENDPOINT,
        region: 'ru-central1',
    });
    const docClient = DynamoDBDocumentClient.from(client);

    switch (event.httpMethod) {
        case 'GET': {
            const command = new QueryCommand({
                TableName: 'demo',
                KeyConditionExpression: 'id = :id',
                ExpressionAttributeValues: {
                    ':id': parseInt(event.queryStringParameters.id, 10)
                }
            });

            // Return a Http.Result object with a status code of 200, custom headers, and a body.
            const response = await docClient.send(command);
            return {
                statusCode: 200,
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify(response.Items)
            };
        }
        case 'POST': {
            const body = JSON.parse(event.body);

            const command = new PutCommand({
                TableName: 'demo',
                Item: {
                    id: body.id,
                    key: body.key,
                    value: body.value
                }
            });

            await docClient.send(command);
            return {
                statusCode: 200,
                body: 'OK'
            };
        }
        default:
            return {
                statusCode: 405,
                headers: {},
                body: 'Method Not Allowed'
            };
    }
}
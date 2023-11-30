import {Http} from '@yandex-cloud/function-types/dist/src/http';

interface RequestBody {
    name: string;
}

export async function handler(event: Http.Event): Promise<Http.Result> {

    const body = JSON.parse(event.body || '{}') as RequestBody;

    const response = {
        message: `Hello, ${body.name || 'World'}!`
    }

    return {
        statusCode: 200,
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(response)
    };

}
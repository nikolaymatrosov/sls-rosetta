import {Http} from '@yandex-cloud/function-types/dist/src/http';
import {pino} from 'pino';

const logger = pino(
    {
        level: 'debug',
        formatters: {
            level: (label) => {
                return {
                    level: label
                }
            }
        }
    }
);

export async function handler(event: Http.Event): Promise<Http.Result> {
    logger.info('Received event', event);
    logger.debug('Debugging event', event);
    logger.error('Error event', event);
    // Return a Http.Result object with a status code of 200, custom headers, and a body.
    return {
        statusCode: 200,
        body: `Hello, ${event.queryStringParameters.name}!`
    };

}
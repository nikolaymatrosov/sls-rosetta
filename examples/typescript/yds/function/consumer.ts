// Consumer function for YDS topic (TypeScript)
import type { Handler } from '@yandex-cloud/function-types';
import type { YDSEvent, YDSResponse } from './types';

export const consumerHandler: Handler.DataStreams = async (event) => {
  // event.messages is always an array according to DataStreams contract
  if (!event || !Array.isArray(event.messages)) {
    return {
      status_code: 400,
      message: 'Invalid event format',
    } satisfies YDSResponse;
  }

  let processed = 0;
  for (const [i, message] of event.messages.entries()) {
    try {
      const eventData = JSON.parse(message.details.data);
      console.log(`Processing message ${i + 1}:`, eventData);
      processed++;
    } catch (err) {
      console.error(`Error parsing message ${i + 1}:`, err);
    }
  }

  return {
    status_code: 200,
    message: `Processed ${processed} messages successfully`,
  } satisfies YDSResponse;
}; 
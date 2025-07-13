import { Driver } from '@ydbjs/core';
import { createTopicWriter } from '@ydbjs/topic/writer';
import type { Handler } from '@yandex-cloud/function-types';
import type { ProducerRequest, ProducerResponse } from './types';

export const producerHandler: Handler.Http = async (event) => {
  const connectionString = process.env['YDB_CONNECTION_STRING'];
  const topicName = process.env['YDS_TOPIC_ID'];
  if (!connectionString || !topicName) {
    return {
      statusCode: 500,
      body: JSON.stringify({ status_code: 500, message: 'YDB_CONNECTION_STRING and YDS_TOPIC_ID environment variables must be set' } satisfies ProducerResponse),
    };
  }

  let body: ProducerRequest;
  try {
    body = typeof event.body === 'string' ? JSON.parse(event.body) : event.body;
  } catch (err) {
    return {
      statusCode: 400,
      body: JSON.stringify({ status_code: 400, message: 'Invalid request format' } satisfies ProducerResponse),
    };
  }

  if (!body.message) {
    return {
      statusCode: 400,
      body: JSON.stringify({ status_code: 400, message: 'Message is required' } satisfies ProducerResponse),
    };
  }

  // Create event data
  const eventData = {
    message: body.message,
    user_id: body.user_id,
    action: body.action,
    timestamp: Math.floor(Date.now() / 1000),
  };

  try {
    const driver = new Driver(connectionString);
    await driver.ready();
    await writeToTopic(driver, topicName, eventData);
    if (typeof (driver as any).close === 'function') {
      await (driver as any).close();
    }
    const response: ProducerResponse = {
      status_code: 200,
      message: 'Data written to topic successfully',
      stream_id: topicName,
    };
    return {
      statusCode: 200,
      body: JSON.stringify(response),
    };
  } catch (err: any) {
    return {
      statusCode: 500,
      body: JSON.stringify({ status_code: 500, message: 'Failed to write to topic: ' + err.message } satisfies ProducerResponse),
    };
  }
};

async function writeToTopic(driver: Driver, topic: string, data: any) {
  const writer = createTopicWriter(driver, { topic });
  try {
    const payload = Buffer.from(JSON.stringify(data), 'utf8');
    await writer.write(payload);
  } finally {
    await writer.close();
  }
} 
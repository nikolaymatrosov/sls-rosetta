import { Driver } from '@ydbjs/core';
import { TopicWriter } from '@ydbjs/topic/writer';
import type { Request, Response } from 'express';

// Types matching Go example
interface ProducerRequest {
  message: string;
  user_id: string;
  action: string;
}

interface ProducerResponse {
  status_code: number;
  message: string;
  stream_id?: string;
}

// Main handler function (Express-style)
export async function producerHandler(req: Request, res: Response) {
  const connectionString = process.env['YDB_CONNECTION_STRING'];
  const topicName = process.env['YDS_TOPIC_ID'];
  if (!connectionString || !topicName) {
    res.status(500).json({ status_code: 500, message: 'YDB_CONNECTION_STRING and YDS_TOPIC_ID environment variables must be set' });
    return;
  }

  let body: ProducerRequest;
  try {
    body = req.body;
  } catch (err) {
    res.status(400).json({ status_code: 400, message: 'Invalid request format' });
    return;
  }

  if (!body.message) {
    res.status(400).json({ status_code: 400, message: 'Message is required' });
    return;
  }

  // Create event data
  const eventData = {
    message: body.message,
    user_id: body.user_id,
    action: body.action,
    timestamp: Math.floor(Date.now() / 1000),
  };

  // Write to YDB topic
  try {
    const driver = new Driver(connectionString);
    await driver.ready();
    await writeToTopic(driver, topicName, eventData);
    await driver.destroy();
    const response: ProducerResponse = {
      status_code: 200,
      message: 'Data written to topic successfully',
      stream_id: topicName,
    };
    res.status(200).json(response);
  } catch (err: any) {
    res.status(500).json({ status_code: 500, message: 'Failed to write to topic: ' + err.message });
  }
}

// Helper to write a message to the topic
async function writeToTopic(driver: Driver, topic: string, data: any) {
  const writer = await TopicWriter.create(driver, { topic });
  try {
    const payload = Buffer.from(JSON.stringify(data), 'utf8');
    await writer.write(payload);
  } finally {
    await writer.close();
  }
} 
import {GetQueueUrlCommand, SendMessageCommand, SendMessageCommandInput, SQSClient} from "@aws-sdk/client-sqs";
import {GetQueueUrlCommandInput} from "@aws-sdk/client-sqs/dist-types/commands/GetQueueUrlCommand";

export async function sendMessageToQueue(ymqName: string, s: string, fromSenderFunction: string) {
    // Create an Amazon SQS service client
    const client = new SQSClient({
        region: 'ru-central1',
        signingRegion: 'ru-central1',
        endpoint: 'https://message-queue.api.cloud.yandex.net'

    });

    const getQueueCmd: GetQueueUrlCommandInput = {
        QueueName: ymqName,
    };

    const urlRes = await client.send(new GetQueueUrlCommand(getQueueCmd));

    // Set the parameters
    const sendMessageInput: SendMessageCommandInput = {
        MessageAttributes: {
            "Origin": {
                DataType: "String",
                StringValue: fromSenderFunction,
            },
        },
        MessageBody: s,
        QueueUrl: urlRes.QueueUrl,
    };

    return client.send(new SendMessageCommand(sendMessageInput));
}
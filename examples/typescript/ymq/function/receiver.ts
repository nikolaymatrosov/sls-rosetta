import {sendMessageToQueue} from "./common";
import {Handler} from "@yandex-cloud/function-types";

// It takes a MessageQueue.Event object as a parameter and returns a Promise that resolves to a Http.Result object.
// noinspection JSUnusedGlobalSymbols
export const receiver: Handler.MessageQueue = async (event) => {

    const ymqName = process.env.YMQ_NAME;

    const messagePromises = event.messages.map(async (message) => {
        const messageBody = JSON.parse(message.details.message.body);

        const respMessage = {
            "result": "success",
            "name": messageBody.name
        }
        return sendMessageToQueue(ymqName, JSON.stringify(respMessage), "From Receiver Function");

    });
    await Promise.all(messagePromises);
    console.log(`Sent ${messagePromises.length} messages`)
    // Return a Http.Result object with a status code of 200, custom headers, and a body.
    return {
        statusCode: 200,
        body: `success`,
    };

}
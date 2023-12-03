import {Http} from '@yandex-cloud/function-types/dist/src/http';
import {sendMessageToQueue} from "./common";

// It takes a Http.Event object as a parameter and returns a Promise that resolves to a Http.Result object.
// noinspection JSUnusedGlobalSymbols
export async function sender(event: Http.Event): Promise<Http.Result> {

    const ymqName = process.env.YMQ_NAME;
    const resp = await sendMessageToQueue(ymqName, `{"name":"test"}`, "From Sender Function");

    console.log("Sent message with ID: " + resp.MessageId)

    // Return a Http.Result object with a status code of 200, custom headers, and a body.
    return {
        statusCode: 200,
        body: `Sent message with ID: ${resp.MessageId}`
    };

}
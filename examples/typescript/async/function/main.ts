export async function handler(event: any, context: any): Promise<object> {

    console.log('event', event);
    console.log('context', context);
    return {
        name: event.name
    };
}
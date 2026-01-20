import { Handler } from '@yandex-cloud/function-types';

/**
 * Consumer function triggered by YDS (Yandex Data Streams).
 *
 * Processes batches of messages from the YDS topic.
 * Messages are automatically decoded and parsed.
 */

interface ProcessedMessage {
    user_id: string;
    action: string;
    status: string;
}

interface UserAction {
    action: string;
    message: string;
    timestamp: string;
}

export const handler: Handler.DataStreams = async (event, context) => {
    try {
        const messages = event.messages || [];
        console.log(`Received ${messages.length} messages from YDS`);

        // Debug: Log the full event structure
        console.log(`Full event: ${JSON.stringify(event, null, 2)}`);

        if (messages.length === 0) {
            return {
                statusCode: 200,
                body: JSON.stringify({ status: 'success', processed: 0 })
            };
        }

        // Process messages
        const processedMessages: ProcessedMessage[] = [];
        const userActions: Record<string, UserAction[]> = {};

        for (const msg of messages) {
            try {
                // Debug: Log the message structure
                console.log(`Message structure: ${JSON.stringify(msg, null, 2)}`);

                // Extract fields
                const userId = msg.user_id;
                const action = msg.action;
                const message = msg.message;
                const timestamp = msg.timestamp;

                console.log(`Processing message: user=${userId}, action=${action}, message=${message}`);

                // Group actions by user
                if (!userActions[userId]) {
                    userActions[userId] = [];
                }
                userActions[userId].push({
                    action,
                    message,
                    timestamp
                });

                // Process based on action type
                processAction(userId, action, message, timestamp);

                processedMessages.push({
                    user_id: userId,
                    action,
                    status: 'processed'
                });
            } catch (error) {
                console.error('Error processing individual message:', error);
                console.error('Error stack:', error instanceof Error ? error.stack : String(error));
                continue;
            }
        }

        // Log batch summary
        console.log(`Batch processing complete. Processed ${processedMessages.length} messages`);
        console.log(`User summary: ${JSON.stringify(userActions)}`);

        // Print summary by user
        for (const [userId, actions] of Object.entries(userActions)) {
            const actionList = actions.map(a => a.action).join(', ');
            console.log(`User ${userId}: ${actions.length} actions - [${actionList}]`);
        }

        return {
            statusCode: 200,
            body: JSON.stringify({
                status: 'success',
                processed: processedMessages.length,
                users_affected: Object.keys(userActions).length
            })
        };

    } catch (error) {
        console.error('Error in consumer handler:', error);
        return {
            statusCode: 500,
            body: JSON.stringify({
                error: error instanceof Error ? error.message : String(error)
            })
        };
    }
};

/**
 * Process different action types with custom business logic.
 *
 * @param userId - User identifier
 * @param action - Action type (login, logout, purchase, view, etc.)
 * @param message - Message content
 * @param timestamp - Message timestamp
 */
function processAction(userId: string, action: string, message: string, timestamp: string): void {
    switch (action) {
        case 'login':
            console.log(`User ${userId} logged in at ${timestamp}`);
            // Add custom logic for login events
            // e.g., update user session, send analytics, etc.
            break;

        case 'logout':
            console.log(`User ${userId} logged out at ${timestamp}`);
            // Add custom logic for logout events
            break;

        case 'purchase':
            console.log(`User ${userId} made a purchase: ${message}`);
            // Add custom logic for purchase events
            // e.g., update inventory, send confirmation email, etc.
            break;

        case 'view':
            console.log(`User ${userId} viewed: ${message}`);
            // Add custom logic for view events
            // e.g., track page views, update recommendations, etc.
            break;

        default:
            console.log(`User ${userId} performed action '${action}': ${message}`);
            // Handle other action types
            break;
    }
}

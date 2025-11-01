import json
import logging
import base64
from collections import defaultdict

# Configure logging
logger = logging.getLogger()
logger.setLevel(logging.INFO)


def consumer_handler(event, context):
    """
    Consumer function triggered by YDS (Yandex Data Streams).

    Processes batches of messages from the YDS topic.
    Messages are automatically decoded and parsed.

    Event structure:
    {
        "messages": [
            {
                "event_metadata": {
                    "event_id": "string",
                    "event_type": "string",
                    "created_at": "timestamp"
                },
                "details": {
                    "topic": "string",
                    "partition": "string",
                    "database": "string",
                    "message": {
                        "data": "base64_encoded_string",
                        "seq_no": "string",
                        "created_at": "timestamp",
                        "message_group_id": "string"
                    }
                }
            }
        ]
    }
    """
    try:
        messages = event.get('messages', [])
        logger.info(f"Received {len(messages)} messages from YDS")

        if not messages:
            return {
                'statusCode': 200,
                'body': json.dumps({'status': 'success', 'processed': 0})
            }

        # Process messages
        processed_messages = []
        user_actions = defaultdict(list)

        for msg in messages:
            try:
                # Extract message data
                details = msg.get('details', {})
                message_data = details.get('message', {})
                data_encoded = message_data.get('data', '')

                # Decode base64 data
                if data_encoded:
                    data_decoded = base64.b64decode(data_encoded).decode('utf-8')
                    data_json = json.loads(data_decoded)

                    # Extract fields
                    user_id = data_json.get('user_id')
                    action = data_json.get('action')
                    message = data_json.get('message')
                    timestamp = data_json.get('timestamp')

                    logger.info(f"Processing message: user={user_id}, action={action}, message={message}")

                    # Group actions by user
                    user_actions[user_id].append({
                        'action': action,
                        'message': message,
                        'timestamp': timestamp
                    })

                    # Process based on action type
                    process_action(user_id, action, message, timestamp)

                    processed_messages.append({
                        'user_id': user_id,
                        'action': action,
                        'status': 'processed'
                    })

            except Exception as e:
                logger.error(f"Error processing individual message: {e}")
                continue

        # Log batch summary
        logger.info(f"Batch processing complete. Processed {len(processed_messages)} messages")
        logger.info(f"User summary: {dict(user_actions)}")

        # Print summary by user
        for user_id, actions in user_actions.items():
            logger.info(f"User {user_id}: {len(actions)} actions - {[a['action'] for a in actions]}")

        return {
            'statusCode': 200,
            'body': json.dumps({
                'status': 'success',
                'processed': len(processed_messages),
                'users_affected': len(user_actions)
            })
        }

    except Exception as e:
        logger.error(f"Error in consumer handler: {e}")
        return {
            'statusCode': 500,
            'body': json.dumps({'error': str(e)})
        }


def process_action(user_id, action, message, timestamp):
    """
    Process different action types with custom business logic.

    Args:
        user_id: User identifier
        action: Action type (login, logout, purchase, view, etc.)
        message: Message content
        timestamp: Message timestamp
    """
    if action == 'login':
        logger.info(f"User {user_id} logged in at {timestamp}")
        # Add custom logic for login events
        # e.g., update user session, send analytics, etc.

    elif action == 'logout':
        logger.info(f"User {user_id} logged out at {timestamp}")
        # Add custom logic for logout events

    elif action == 'purchase':
        logger.info(f"User {user_id} made a purchase: {message}")
        # Add custom logic for purchase events
        # e.g., update inventory, send confirmation email, etc.

    elif action == 'view':
        logger.info(f"User {user_id} viewed: {message}")
        # Add custom logic for view events
        # e.g., track page views, update recommendations, etc.

    else:
        logger.info(f"User {user_id} performed action '{action}': {message}")
        # Handle other action types

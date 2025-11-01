import json
import os
import logging
from datetime import datetime
import ydb
import ydb.iam

# Configure logging
logger = logging.getLogger()
logger.setLevel(logging.INFO)


def producer_handler(event, context):
    """
    Producer function that receives HTTP requests and writes messages to YDS topic.

    Expected request body:
    {
        "message": "string",
        "user_id": "string",
        "action": "string"
    }
    """
    try:
        # Parse request body
        body = event.get('body', '{}')
        if event.get('isBase64Encoded', False):
            import base64
            body = base64.b64decode(body).decode('utf-8')

        data = json.loads(body)

        # Validate required fields
        required_fields = ['message', 'user_id', 'action']
        for field in required_fields:
            if field not in data:
                return {
                    'statusCode': 400,
                    'body': json.dumps({'error': f'Missing required field: {field}'})
                }

        # Get environment variables
        ydb_endpoint = os.environ.get('YDB_ENDPOINT')
        yds_topic_path = os.environ.get('YDS_TOPIC_PATH')

        if not ydb_endpoint or not yds_topic_path:
            logger.error("Missing environment variables: YDB_ENDPOINT or YDS_TOPIC_PATH")
            return {
                'statusCode': 500,
                'body': json.dumps({'error': 'Configuration error'})
            }

        # Prepare message with timestamp
        message = {
            'message': data['message'],
            'user_id': data['user_id'],
            'action': data['action'],
            'timestamp': datetime.utcnow().isoformat()
        }

        # Parse YDB endpoint to extract database path
        # Format: grpcs://host:port/?database=/path/to/db
        database_path = yds_topic_path.rsplit('/', 1)[0]  # Remove topic name to get database path

        # Connect to YDB and write to topic
        driver_config = ydb.DriverConfig(
            ydb_endpoint,
            database=database_path,
            credentials=ydb.iam.MetadataUrlCredentials(),
        )

        with ydb.Driver(driver_config) as driver:
            driver.wait(timeout=5, fail_fast=True)

            # Write message to topic
            topic_path = yds_topic_path.split(database_path)[-1].lstrip('/')  # Get relative topic path
            with driver.topic_client.writer(topic_path) as writer:
                writer.write(
                    ydb.TopicWriterMessage(
                        data=json.dumps(message).encode('utf-8')
                    )
                )
                writer.flush()

        logger.info(f"Message sent to YDS topic: {message}")

        return {
            'statusCode': 200,
            'body': json.dumps({
                'status': 'success',
                'message': 'Message sent to YDS topic',
                'data': message
            })
        }

    except json.JSONDecodeError as e:
        logger.error(f"Invalid JSON in request body: {e}")
        return {
            'statusCode': 400,
            'body': json.dumps({'error': 'Invalid JSON format'})
        }
    except Exception as e:
        logger.error(f"Error processing request: {e}")
        return {
            'statusCode': 500,
            'body': json.dumps({'error': str(e)})
        }

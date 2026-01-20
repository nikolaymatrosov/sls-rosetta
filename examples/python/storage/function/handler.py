import json
import logging
import os
from io import BytesIO
from typing import Dict, List, Any

import boto3
from PIL import Image

# Configure logging
logger = logging.getLogger()
logger.setLevel(logging.INFO)

# Constants
THUMBNAIL_SIZE = (100, 100)
THUMBNAIL_PREFIX = "thumbnails/"
STORAGE_ENDPOINT = "https://storage.yandexcloud.net"
REGION = "ru-central1"


def handler(event: Dict[str, Any], context: Any) -> Dict[str, int]:
    """
    Main handler for Object Storage events.

    Processes uploaded images and creates thumbnails.

    Args:
        event: Object Storage event with structure:
            {
                "messages": [
                    {
                        "event_metadata": {...},
                        "details": {
                            "bucket_id": "bucket-name",
                            "object_id": "uploads/image.png"
                        }
                    }
                ]
            }
        context: Yandex Cloud Function context

    Returns:
        Response dict with statusCode
    """
    try:
        # Initialize S3 client
        s3_client = get_s3_client()

        # Extract messages
        messages = event.get('messages', [])
        logger.info(f"Processing {len(messages)} messages")

        # Process each message
        for message in messages:
            try:
                process_message(s3_client, message)
            except Exception as e:
                logger.error(f"Error processing message: {e}", exc_info=True)
                # Continue processing other messages

        return {"statusCode": 200}

    except Exception as e:
        logger.error(f"Error in handler: {e}", exc_info=True)
        return {"statusCode": 500}


def get_s3_client():
    """Initialize and return S3 client for Yandex Object Storage."""
    return boto3.client(
        's3',
        endpoint_url=STORAGE_ENDPOINT,
        aws_access_key_id=os.environ['AWS_ACCESS_KEY_ID'],
        aws_secret_access_key=os.environ['AWS_SECRET_ACCESS_KEY'],
        region_name=REGION
    )


def process_message(s3_client, message: Dict[str, Any]) -> None:
    """
    Process single Object Storage message.

    Args:
        s3_client: boto3 S3 client
        message: Message dict with details
    """
    details = message.get('details', {})
    bucket_id = details.get('bucket_id')
    object_id = details.get('object_id')

    if not bucket_id or not object_id:
        logger.warning("Missing bucket_id or object_id in message")
        return

    logger.info(f"Processing: {bucket_id}/{object_id}")

    # Download image
    image_bytes = download_image(s3_client, bucket_id, object_id)

    # Resize image
    thumbnail_bytes = resize_image(image_bytes)

    # Upload thumbnail
    thumbnail_key = get_thumbnail_key(object_id)
    upload_thumbnail(s3_client, bucket_id, thumbnail_key, thumbnail_bytes)

    logger.info(f"Created thumbnail: {bucket_id}/{thumbnail_key}")


def download_image(s3_client, bucket: str, key: str) -> bytes:
    """Download image from S3."""
    response = s3_client.get_object(Bucket=bucket, Key=key)
    return response['Body'].read()


def resize_image(image_bytes: bytes) -> bytes:
    """
    Resize image to thumbnail size.

    Args:
        image_bytes: Original image bytes

    Returns:
        Resized image bytes
    """
    # Open image
    image = Image.open(BytesIO(image_bytes))

    # Get original format
    original_format = image.format or 'PNG'

    # Resize using high-quality resampling
    image.thumbnail(THUMBNAIL_SIZE, Image.LANCZOS)

    # Save to bytes
    output = BytesIO()
    image.save(output, format=original_format)
    output.seek(0)

    return output.read()


def upload_thumbnail(s3_client, bucket: str, key: str, data: bytes) -> None:
    """Upload thumbnail to S3."""
    s3_client.put_object(
        Bucket=bucket,
        Key=key,
        Body=data,
        ContentType=get_content_type(key)
    )


def get_thumbnail_key(original_key: str) -> str:
    """
    Generate thumbnail key from original object key.

    uploads/image.png -> thumbnails/image.png
    """
    filename = original_key.split('/')[-1]
    return f"{THUMBNAIL_PREFIX}{filename}"


def get_content_type(key: str) -> str:
    """Determine content type from file extension."""
    ext = key.lower().split('.')[-1]
    content_types = {
        'jpg': 'image/jpeg',
        'jpeg': 'image/jpeg',
        'png': 'image/png',
        'gif': 'image/gif',
        'webp': 'image/webp'
    }
    return content_types.get(ext, 'application/octet-stream')

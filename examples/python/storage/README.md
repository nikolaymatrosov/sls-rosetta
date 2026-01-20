# Yandex Object Storage Trigger Example - Python

This example demonstrates how to use Yandex Object Storage triggers with Yandex Cloud Functions in Python. It implements an automatic image thumbnail generation service:

1. **Object Storage Bucket**: Stores uploaded images and generated thumbnails
2. **Storage Trigger**: Automatically invokes the function when images are uploaded to `uploads/` folder
3. **Handler Function**: Downloads images, resizes them to 100x100 pixels, and saves thumbnails to `thumbnails/` folder
4. **Image Processing**: Uses Pillow (pure Python library) for high-quality image resizing

## Architecture

```txt
Image Upload → S3 Bucket (uploads/) → Storage Trigger → Handler Function → S3 Bucket (thumbnails/)
```

## Key Features

- **Object Storage trigger** on `uploads/` prefix
- **Image resizing** using Pillow (pure Python, no binary dependencies)
- **boto3** for S3-compatible operations with Yandex Object Storage
- **Automatic thumbnail generation** at 100x100 pixels
- **Error handling** and comprehensive logging
- **Multiple format support**: PNG, JPEG, GIF, WebP

## Prerequisites

- Yandex Cloud account
- Terraform >= 1.0
- AWS CLI (for testing S3 operations)
- curl (for testing)

## Project Structure

```txt
examples/python/storage/
├── README.md                    # This file
├── function/
│   ├── handler.py              # Main handler function
│   └── requirements.txt        # Python dependencies (boto3, Pillow)
├── tf/
│   ├── terraform.tf            # Provider and backend configuration
│   ├── variables.tf            # Input variables
│   ├── outputs.tf              # Output values
│   ├── iam.tf                  # Service accounts and IAM roles
│   ├── storage.tf              # S3 bucket configuration
│   └── main.tf                 # Function and trigger resources
└── environment/
    └── terraform.tfstate       # Terraform state file
```

## Deployment

### 1. Configure Variables

Create a `terraform.tfvars` file in the `tf/` directory:

```hcl
cloud_id  = "your-cloud-id"
folder_id = "your-folder-id"
zone      = "ru-central1-a"  # Optional, default value
```

### 2. Initialize Terraform

```bash
cd tf
terraform init
```

### 3. Deploy Infrastructure

```bash
terraform apply
```

This will create:

- S3-compatible storage bucket with random UUID name
- Storage handler function (Python 3.12 runtime)
- Storage trigger (monitors `uploads/` prefix)
- Service accounts with appropriate IAM roles
- Static access keys for S3 operations

### 4. Get Outputs

After deployment, get important information:

```bash
terraform output
```

Key outputs:

- `bucket_name`: Name of the storage bucket
- `function_id`: ID of the handler function
- `test_upload_command`: Command to upload test images
- `test_list_command`: Command to list generated thumbnails

## Usage

### Upload a Test Image

Use the command from terraform output, or manually:

```bash
BUCKET=$(terraform output -raw bucket_name)

aws s3 cp --endpoint-url=https://storage.yandexcloud.net \
  ./your-image.png \
  s3://$BUCKET/uploads/your-image.png
```

### Verify Thumbnail Creation

Wait a few seconds for the function to process, then check:

```bash
BUCKET=$(terraform output -raw bucket_name)

# List thumbnails
aws s3 ls --endpoint-url=https://storage.yandexcloud.net \
  s3://$BUCKET/thumbnails/

# Download thumbnail
aws s3 cp --endpoint-url=https://storage.yandexcloud.net \
  s3://$BUCKET/thumbnails/your-image.png \
  ./thumbnail.png
```

### View Function Logs

```bash
FUNCTION_ID=$(terraform output -raw function_id)

yc serverless function logs $FUNCTION_ID
```

## How It Works

### Event Structure

When an object is created or updated in the bucket with `uploads/` prefix, the storage trigger sends an event to the function:

```python
{
    "messages": [
        {
            "event_metadata": {
                "created_at": "2024-01-20T12:00:00Z",
                "cloud_id": "b1g...",
                "folder_id": "b1g..."
            },
            "details": {
                "bucket_id": "bucket-name",
                "object_id": "uploads/image.png"
            }
        }
    ]
}
```

### Handler Logic

The [handler.py](function/handler.py) function processes each message:

1. **Extract details**: Parses `bucket_id` and `object_id` from event
2. **Download image**: Uses boto3 to download from `uploads/` folder
3. **Resize image**: Opens with Pillow and resizes to 100x100 pixels using LANCZOS resampling
4. **Upload thumbnail**: Saves to `thumbnails/` folder with same filename
5. **Log results**: Records success or errors for monitoring

Key code snippet:

```python
def resize_image(image_bytes: bytes) -> bytes:
    """Resize image to 100x100 thumbnail."""
    image = Image.open(BytesIO(image_bytes))
    original_format = image.format or 'PNG'

    # High-quality resize
    image.thumbnail((100, 100), Image.LANCZOS)

    output = BytesIO()
    image.save(output, format=original_format)
    output.seek(0)

    return output.read()
```

### Storage Trigger Configuration

The trigger is configured in [main.tf](tf/main.tf):

```hcl
object_storage {
  bucket_id    = yandex_storage_bucket.for_uploads.bucket
  prefix       = "uploads/"      # Only trigger on uploads/ folder
  create       = true            # Trigger on new objects
  update       = true            # Trigger on updated objects
  batch_cutoff = 1               # Process immediately
}

function {
  id                 = yandex_function.storage_handler.id
  service_account_id = yandex_iam_service_account.trigger_sa.id
  retry_attempts     = 3         # Retry on failures
  retry_interval     = 10        # 10 seconds between retries
}
```

## Testing

### Basic Test

1. Upload a PNG or JPEG image to `uploads/`:

```bash
BUCKET=$(terraform output -raw bucket_name)
aws s3 cp --endpoint-url=https://storage.yandexcloud.net \
  ./test.png s3://$BUCKET/uploads/test.png
```

2. Wait 5-10 seconds for processing

3. Verify thumbnail exists:

```bash
aws s3 ls --endpoint-url=https://storage.yandexcloud.net \
  s3://$BUCKET/thumbnails/test.png
```

4. Download and verify dimensions:

```bash
aws s3 cp --endpoint-url=https://storage.yandexcloud.net \
  s3://$BUCKET/thumbnails/test.png ./thumbnail.png

# Check dimensions (requires ImageMagick)
identify thumbnail.png
# Should show: thumbnail.png PNG 100x100 ...
```

### Multiple Images

```bash
BUCKET=$(terraform output -raw bucket_name)

# Upload multiple images
for img in image1.png image2.jpg image3.png; do
  aws s3 cp --endpoint-url=https://storage.yandexcloud.net \
    ./$img s3://$BUCKET/uploads/$img
done

# Wait a moment
sleep 10

# List all thumbnails
aws s3 ls --endpoint-url=https://storage.yandexcloud.net \
  s3://$BUCKET/thumbnails/
```

## Comparison with Go Example

This Python implementation is functionally equivalent to the [Go storage example](../../go/storage/) but with some differences:

| Aspect | Go | Python |
|--------|-----|--------|
| **Image Library** | libvips (C bindings) | Pillow (pure Python) |
| **Binary Dependencies** | Yes (shared libraries) | No |
| **Build Process** | Docker + plugin compilation | Simple zip archive |
| **Memory** | 128 MB | 256 MB |
| **Timeout** | 10 seconds | 30 seconds |
| **Performance** | Faster (compiled) | Slower (interpreted) |
| **Deployment Complexity** | Higher (Docker, shared libs) | Lower (direct zip) |
| **Concurrency** | Goroutines | Sequential processing |
| **Deployment Bucket** | Required | Not needed |

**Why Python?**

- **Simpler deployment**: No Docker build, no shared library management
- **Easier maintenance**: Pure Python dependencies
- **Lower complexity**: Fewer moving parts
- **Good enough performance**: For thumbnail generation, the extra seconds don't matter

**When to use Go instead?**

- High-throughput scenarios (thousands of images per minute)
- Strict latency requirements (< 1 second processing)
- Advanced image processing (Go's libvips has more features)

## Customization

### Change Thumbnail Size

Edit [handler.py](function/handler.py):

```python
# Change from 100x100 to 200x200
THUMBNAIL_SIZE = (200, 200)
```

### Filter by File Extension

Edit [main.tf](tf/main.tf) trigger configuration:

```hcl
object_storage {
  bucket_id = yandex_storage_bucket.for_uploads.bucket
  prefix    = "uploads/"
  suffix    = ".jpg"  # Only process .jpg files
  # ...
}
```

### Change Thumbnail Folder

Edit [handler.py](function/handler.py):

```python
# Change from thumbnails/ to resized/
THUMBNAIL_PREFIX = "resized/"
```

### Adjust Resampling Quality

Edit [handler.py](function/handler.py):

```python
# Change resampling algorithm
# Options: Image.NEAREST, Image.BILINEAR, Image.BICUBIC, Image.LANCZOS
image.thumbnail(THUMBNAIL_SIZE, Image.BICUBIC)
```

## Monitoring

### View Function Logs

```bash
yc serverless function logs $(terraform output -raw function_id) --follow
```

### Check Function Metrics

```bash
yc serverless function version list --function-name storage-handler-python
```

### Debug Issues

Common log patterns:

- `Processing: bucket/uploads/image.png` - Image download started
- `Created thumbnail: bucket/thumbnails/image.png` - Success
- `Error processing message: ...` - Individual message failed (continues with others)
- `Missing bucket_id or object_id` - Malformed event

## Troubleshooting

### Thumbnail Not Created

**Check function logs**:

```bash
yc serverless function logs $(terraform output -raw function_id) --limit 50
```

**Common issues**:

1. **Permission denied**: Wait 15 seconds after deployment for IAM propagation
2. **Invalid image format**: Handler only supports PNG, JPEG, GIF, WebP
3. **File too large**: Large images may exceed 30s timeout, increase in [main.tf](tf/main.tf)
4. **Wrong prefix**: Ensure file is in `uploads/` folder, not root

### Function Timeout

For very large images, increase timeout in [main.tf](tf/main.tf):

```hcl
resource "yandex_function" "storage_handler" {
  # ...
  execution_timeout = "60"  # Increase to 60 seconds
}
```

### Out of Memory

For high-resolution images, increase memory in [main.tf](tf/main.tf):

```hcl
resource "yandex_function" "storage_handler" {
  # ...
  memory = 512  # Increase to 512 MB
}
```

### IAM Permission Errors

Wait 15-30 seconds after `terraform apply` for IAM roles to propagate. Then try uploading again.

## Cleanup

### Delete All Objects First

```bash
BUCKET=$(terraform output -raw bucket_name)

# Delete all objects (uploads and thumbnails)
aws s3 rm --endpoint-url=https://storage.yandexcloud.net \
  s3://$BUCKET --recursive
```

### Destroy Infrastructure

```bash
terraform destroy
```

## Related Examples

- [Go Storage Example](../../go/storage/) - Same functionality with Go and libvips
- [Python YDS Example](../yds/) - Demonstrates Data Streams triggers
- [Python GigaChat Example](../gigachat/) - Shows HTTP-triggered functions

## Resources

- [Yandex Object Storage Documentation](https://cloud.yandex.com/docs/storage/)
- [Storage Triggers](https://cloud.yandex.com/docs/functions/concepts/trigger/os-trigger)
- [Yandex Cloud Functions - Python Runtime](https://cloud.yandex.com/docs/functions/lang/python/)
- [boto3 Documentation](https://boto3.amazonaws.com/v1/documentation/api/latest/index.html)
- [Pillow Documentation](https://pillow.readthedocs.io/)

# Developer Guide: Adding Examples with Terratest

This guide explains how to add new examples to the sls-rosetta repository with corresponding Terratest integration tests, based on established patterns and conventions.

## Table of Contents

- [Introduction](#introduction)
- [Repository Structure](#repository-structure)
- [Quick Start](#quick-start)
- [Example Structure](#example-structure)
- [Terraform Patterns](#terraform-patterns)
- [Function Handler Patterns](#function-handler-patterns)
- [Terratest Patterns](#terratest-patterns)
- [Python Storage Example Walkthrough](#python-storage-example-walkthrough)
- [Step-by-Step Guide](#step-by-step-guide)
- [Common Patterns and Conventions](#common-patterns-and-conventions)
- [Testing Best Practices](#testing-best-practices)
- [Troubleshooting](#troubleshooting)
- [References](#references)

## Introduction

The `sls-rosetta` repository contains examples demonstrating Yandex Cloud Serverless Functions across multiple programming languages (Go, TypeScript, Python, C#, Java). Each example includes:

- **Function code** in the target language
- **Terraform configuration** for infrastructure as code
- **Terratest integration tests** written in Go
- **Comprehensive documentation** (README.md)

This guide helps you contribute new examples following the established patterns.

## Repository Structure

```txt
sls-rosetta/
├── examples/
│   ├── go/           # Go examples (hello, storage, yds, etc.)
│   ├── typescript/   # TypeScript examples (apigw, async, ws, etc.)
│   ├── python/       # Python examples (storage, yds, gigachat)
│   ├── csharp/       # C# examples
│   └── java/         # Java examples
├── tests/
│   ├── go/           # Go tests matching examples/go/*
│   ├── typescript/   # TypeScript tests matching examples/typescript/*
│   └── python/       # Python tests matching examples/python/*
├── go.mod            # Go module for tests
└── go.sum
```

**Key Principles:**

- All tests are written in Go using Terratest, regardless of the example language
- Directory structure mirrors between `examples/` and `tests/`
- Each example is self-contained with its own Terraform state

## Quick Start

To add a new example:

1. **Create directory structure**: `examples/{language}/{example-name}/`
2. **Implement function handler**: `function/handler.{ext}` or `function/main.{ext}`
3. **Configure Terraform**: Modular files in `tf/` directory
4. **Write Terratest test**: `tests/{language}/{example-name}/default_test.go`
5. **Document**: Comprehensive `README.md` in the example directory
6. **Test**: Run `terraform apply` and `go test` to verify

## Example Structure

Every example follows this consistent structure:

```txt
examples/{language}/{example-name}/
├── README.md                    # Comprehensive documentation
├── function/                    # Source code directory
│   ├── handler.py              # Python: handler.py
│   ├── main.go                 # Go: main.go
│   ├── main.ts                 # TypeScript: main.ts
│   ├── requirements.txt        # Python: pip dependencies
│   ├── package.json            # TypeScript/Node.js: npm dependencies
│   └── go.mod                  # Go: module definition
├── tf/                          # Terraform configuration
│   ├── terraform.tf            # Provider and backend config
│   ├── variables.tf            # Input variables
│   ├── outputs.tf              # Output values
│   ├── main.tf                 # Function and trigger resources
│   ├── iam.tf                  # Service accounts and IAM
│   ├── {service}.tf            # Service-specific resources (storage.tf, yds.tf)
│   └── .terraform.lock.hcl     # Provider version lock
└── environment/                 # Local Terraform state
    └── terraform.tfstate
```

## Terraform Patterns

### terraform.tf - Provider Configuration

**Template:**

```hcl
terraform {
  required_providers {
    yandex = {
      source  = "yandex-cloud/yandex"
      version = ">= 0.100"
    }
    archive = {
      source  = "hashicorp/archive"
      version = ">= 2.0"
    }
  }
  required_version = ">= 1.0"

  backend "local" {
    path = "../environment/terraform.tfstate"
  }
}

provider "yandex" {
  cloud_id  = var.cloud_id
  folder_id = var.folder_id
  zone      = var.zone
}
```

**Key Points:**

- Always use local backend pointing to `../environment/terraform.tfstate`
- Include archive provider for packaging function code
- Specify minimum versions for reproducibility

### variables.tf - Standard Inputs

**Template:**

```hcl
variable "cloud_id" {
  description = "Yandex Cloud ID"
  type        = string
}

variable "folder_id" {
  description = "Yandex Cloud Folder ID"
  type        = string
}

variable "zone" {
  description = "Yandex Cloud availability zone"
  type        = string
  default     = "ru-central1-a"
}
```

**Key Points:**

- These three variables are standard across all examples
- `cloud_id` and `folder_id` are required
- `zone` has a sensible default

### outputs.tf - Test Integration

**Template:**

```hcl
output "function_id" {
  description = "ID of the deployed function"
  value       = yandex_function.example.id
}

output "function_url" {
  description = "URL to invoke the function"
  value       = "https://functions.yandexcloud.net/${yandex_function.example.id}"
}

# Service-specific outputs
output "bucket_name" {
  description = "Name of the storage bucket"
  value       = yandex_storage_bucket.example.bucket
}

output "access_key" {
  description = "S3 access key for testing"
  value       = yandex_iam_service_account_static_access_key.example.access_key
  sensitive   = true
}

output "secret_key" {
  description = "S3 secret key for testing"
  value       = yandex_iam_service_account_static_access_key.example.secret_key
  sensitive   = true
}
```

**Key Points:**

- Always output function IDs for test verification
- Include service-specific outputs (bucket names, topic IDs, etc.)
- Mark sensitive outputs (credentials) as `sensitive = true`
- Terratest retrieves these outputs using `terraform.Output()`

### main.tf - Function and Trigger Resources

**Template:**

```hcl
# Archive function code
data "archive_file" "function_code" {
  type        = "zip"
  source_dir  = "../function"
  output_path = "./function.zip"
  excludes    = ["__pycache__", "*.pyc", ".pytest_cache", "node_modules"]
}

# Serverless Function
resource "yandex_function" "example" {
  name               = "example-function-{language}"
  description        = "Example function demonstrating {feature}"
  user_hash          = data.archive_file.function_code.output_sha256
  runtime            = "python312"  # or golang123, nodejs20, dotnet8
  entrypoint         = "handler.handler"  # or main.Handler
  memory             = 256
  execution_timeout  = "30"
  service_account_id = yandex_iam_service_account.function_sa.id

  environment = {
    # Service-specific configuration
    AWS_ACCESS_KEY_ID     = yandex_iam_service_account_static_access_key.sa.access_key
    AWS_SECRET_ACCESS_KEY = yandex_iam_service_account_static_access_key.sa.secret_key
  }

  content {
    zip_filename = data.archive_file.function_code.output_path
  }

  depends_on = [
    yandex_resourcemanager_folder_iam_member.function_sa
  ]
}

# Storage Trigger (example)
resource "yandex_function_trigger" "storage_trigger" {
  name        = "storage-trigger-{language}"
  description = "Trigger on object storage events"
  folder_id   = var.folder_id

  object_storage {
    bucket_id    = yandex_storage_bucket.example.bucket
    prefix       = "uploads/"
    create       = true
    update       = true
    batch_cutoff = 1
  }

  function {
    id                 = yandex_function.example.id
    service_account_id = yandex_iam_service_account.trigger_sa.id
    retry_attempts     = 3
    retry_interval     = 10
  }

  depends_on = [
    yandex_resourcemanager_folder_iam_member.trigger_sa
  ]
}
```

**Key Points:**

- Use `user_hash = archive_file.output_sha256` for change detection
- Specify appropriate runtime: `python312`, `golang123`, `nodejs20`, `dotnet8`
- Entrypoint format varies by language:
  - Python: `handler.handler` (module.function)
  - Go: `main.Handler` (package.function)
  - TypeScript: `main.handler` (file.function)
- Include `depends_on` for IAM resources to ensure proper ordering
- Configure triggers with appropriate batch settings

### iam.tf - Service Accounts and Permissions

**Template:**

```hcl
# Service Account for Function Execution
resource "yandex_iam_service_account" "function_sa" {
  name        = "{example}-function-sa"
  description = "Service account for {example} function execution"
  folder_id   = var.folder_id
}

# IAM Role Binding
resource "yandex_resourcemanager_folder_iam_member" "function_sa" {
  folder_id = var.folder_id
  role      = "storage.editor"  # or ydb.editor, etc.
  member    = "serviceAccount:${yandex_iam_service_account.function_sa.id}"

  sleep_after = 5  # Wait for IAM propagation
}

# Static Access Key (for S3-compatible services)
resource "yandex_iam_service_account_static_access_key" "function_sa" {
  service_account_id = yandex_iam_service_account.function_sa.id
  description        = "Static access key for S3 operations"
}

# Trigger Service Account
resource "yandex_iam_service_account" "trigger_sa" {
  name        = "{example}-trigger-sa"
  description = "Service account for trigger invocation"
  folder_id   = var.folder_id
}

resource "yandex_resourcemanager_folder_iam_member" "trigger_sa" {
  folder_id = var.folder_id
  role      = "functions.functionInvoker"
  member    = "serviceAccount:${yandex_iam_service_account.trigger_sa.id}"

  sleep_after = 5
}
```

**Common IAM Roles:**

- `functions.functionInvoker` - Invoke functions (for triggers)
- `storage.editor` - Read/write Object Storage
- `storage.viewer` - Read-only Object Storage
- `ydb.editor` - Write to YDB/YDS topics
- `ydb.viewer` - Read from YDB/YDS topics
- `ydb.admin` - Admin operations (for triggers)

**Key Points:**

- Create separate service accounts for different roles (function, trigger)
- Always include `sleep_after = 5` for IAM propagation
- Use static access keys for SDK interactions (S3, etc.)
- Include explicit `depends_on` in resources using these service accounts

## Function Handler Patterns

### Python

**File:** `function/handler.py`

**Template:**

```python
import logging
from typing import Dict, Any

logger = logging.getLogger()
logger.setLevel(logging.INFO)

def handler(event: Dict[str, Any], context: Any) -> Dict[str, int]:
    """
    Main handler function.

    Args:
        event: Event data (structure depends on trigger type)
        context: Runtime context

    Returns:
        Response with statusCode
    """
    logger.info(f"Received event: {event}")

    try:
        # Process event
        result = process_event(event)
        logger.info(f"Processing successful: {result}")
        return {"statusCode": 200}
    except Exception as e:
        logger.error(f"Error processing event: {e}")
        return {"statusCode": 500}

def process_event(event: Dict[str, Any]) -> Any:
    """Process the incoming event."""
    # Implementation
    pass
```

**Dependencies:** `function/requirements.txt`

```txt
boto3>=1.34.0
# Add other dependencies
```

### Go

**File:** `function/main.go`

**Template:**

```go
package main

import (
    "context"
    "fmt"
)

// Event represents the incoming event structure
type Event struct {
    // Define event fields based on trigger type
}

// Response represents the function response
type Response struct {
    StatusCode int    `json:"statusCode"`
    Body       string `json:"body,omitempty"`
}

// Handler is the main entry point
func Handler(ctx context.Context, event Event) (Response, error) {
    fmt.Printf("Received event: %+v\n", event)

    // Process event
    if err := processEvent(event); err != nil {
        return Response{StatusCode: 500}, err
    }

    return Response{StatusCode: 200, Body: "Success"}, nil
}

func processEvent(event Event) error {
    // Implementation
    return nil
}
```

**Dependencies:** `function/go.mod`

```go
module example

go 1.23

require (
    // Add dependencies
)
```

### TypeScript

**File:** `function/main.ts`

**Template:**

```typescript
import { Handler } from '@yandex-cloud/function-types';

interface Event {
    // Define event structure
}

interface Response {
    statusCode: number;
    body?: string;
}

export const handler: Handler<Event, Response> = async (event, context) => {
    console.log('Received event:', event);

    try {
        // Process event
        await processEvent(event);
        return { statusCode: 200, body: 'Success' };
    } catch (error) {
        console.error('Error:', error);
        return { statusCode: 500 };
    }
};

async function processEvent(event: Event): Promise<void> {
    // Implementation
}
```

**Dependencies:** `function/package.json`

```json
{
  "name": "example-function",
  "version": "1.0.0",
  "main": "dist/main.js",
  "scripts": {
    "build": "tsc"
  },
  "dependencies": {
    "@yandex-cloud/function-types": "^2.0.0"
  },
  "devDependencies": {
    "typescript": "^5.0.0",
    "@types/node": "^20.0.0"
  }
}
```

**Build Configuration:** `function/tsconfig.json`

```json
{
  "compilerOptions": {
    "target": "ES2020",
    "module": "commonjs",
    "outDir": "./dist",
    "rootDir": "./",
    "strict": true,
    "esModuleInterop": true
  }
}
```

## Terratest Patterns

### Test Structure

**File:** `tests/{language}/{example}/default_test.go`

**Template:**

```go
package examplename

import (
    "context"
    "os"
    "testing"
    "time"

    "github.com/gruntwork-io/terratest/modules/terraform"
    "github.com/stretchr/testify/assert"
)

func TestExampleName(t *testing.T) {
    // Configure Terratest options
    terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
        TerraformDir: "../../../examples/{language}/{example}/tf",
        Vars: map[string]interface{}{
            "cloud_id":  os.Getenv("CLOUD_ID"),
            "folder_id": os.Getenv("FOLDER_ID"),
        },
        EnvVars: map[string]string{
            "YC_TOKEN": os.Getenv("YC_TOKEN"),
        },
    })

    // Deploy infrastructure
    terraform.InitAndApply(t, terraformOptions)

    // Setup cleanup
    defer terraform.Destroy(t, terraformOptions)

    // Retrieve outputs
    functionID := terraform.Output(t, terraformOptions, "function_id")

    // Wait for IAM propagation
    time.Sleep(5 * time.Second)

    // Run tests
    assert.NotEmpty(t, functionID, "Function ID should not be empty")

    // Additional test logic here
}
```

### Custom AWS Endpoint Resolver (S3)

For testing S3-compatible Object Storage:

```go
import (
    "context"
    "net/url"

    "github.com/aws/aws-sdk-go-v2/aws"
    "github.com/aws/aws-sdk-go-v2/config"
    "github.com/aws/aws-sdk-go-v2/credentials"
    "github.com/aws/aws-sdk-go-v2/service/s3"
    smithyendpoints "github.com/aws/smithy-go/endpoints"
)

type resolverV2 struct{}

func (*resolverV2) ResolveEndpoint(ctx context.Context, params s3.EndpointParameters) (
    smithyendpoints.Endpoint, error,
) {
    u, err := url.Parse("https://storage.yandexcloud.net")
    if err != nil {
        return smithyendpoints.Endpoint{}, err
    }
    u.Path += "/" + *params.Bucket
    return smithyendpoints.Endpoint{
        URI: *u,
    }, nil
}

func createS3Client(ctx context.Context, accessKey, secretKey string) (*s3.Client, error) {
    cfg, err := config.LoadDefaultConfig(ctx,
        config.WithDefaultRegion("ru-central1"),
    )
    if err != nil {
        return nil, err
    }

    return s3.NewFromConfig(cfg, func(o *s3.Options) {
        o.Region = "ru-central1"
        o.EndpointResolverV2 = &resolverV2{}
        o.Credentials = aws.NewCredentialsCache(
            credentials.NewStaticCredentialsProvider(accessKey, secretKey, ""),
        )
    }), nil
}
```

### Test Workflow

1. **Initialize Terratest**

   ```go
   terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
       TerraformDir: "../../../examples/{language}/{example}/tf",
       // ...
   })
   ```

2. **Deploy Infrastructure**

   ```go
   terraform.InitAndApply(t, terraformOptions)
   ```

3. **Setup Cleanup (deferred)**

   ```go
   defer func() {
       // Clean up test resources (S3 objects, etc.)
       terraform.Destroy(t, terraformOptions)
   }()
   ```

4. **Wait for IAM Propagation**

   ```go
   time.Sleep(5 * time.Second)
   ```

5. **Retrieve Outputs**

   ```go
   bucket := terraform.Output(t, terraformOptions, "bucket_name")
   ```

6. **Run Test Logic**

   ```go
   // Upload file, invoke function, verify results
   ```

7. **Wait for Processing**

   ```go
   time.Sleep(5 * time.Second)  // Trigger processing time
   ```

8. **Assert Results**

   ```go
   assert.NoError(t, err, "Failed to get result")
   ```

## Python Storage Example Walkthrough

The [python/storage](examples/python/storage) example demonstrates image thumbnail generation using Object Storage triggers. It's an excellent reference for understanding patterns.

### Directory Structure

```txt
examples/python/storage/
├── README.md
├── function/
│   ├── handler.py          # Main handler (downloads, resizes, uploads)
│   └── requirements.txt    # boto3, Pillow
├── tf/
│   ├── terraform.tf        # Providers: yandex, archive
│   ├── variables.tf        # cloud_id, folder_id, zone
│   ├── outputs.tf          # bucket_name, function_id, credentials
│   ├── main.tf             # Function + Storage Trigger
│   ├── storage.tf          # S3 bucket with random UUID
│   └── iam.tf              # 2 service accounts (handler, trigger)
└── environment/
    └── terraform.tfstate
```

### Handler Implementation

[handler.py](examples/python/storage/function/handler.py:1) processes storage events:

```python
def handler(event: Dict[str, Any], context: Any) -> Dict[str, int]:
    """Process batch of Object Storage events."""
    s3_client = get_s3_client()

    for message in event.get('messages', []):
        try:
            process_message(s3_client, message)
        except Exception as e:
            logger.error(f"Error processing message: {e}")
            # Continue with other messages

    return {"statusCode": 200}

def process_message(s3_client, message: Dict[str, Any]) -> None:
    """Process single storage event."""
    details = message.get('details', {})
    bucket = details.get('bucket_id')
    object_key = details.get('object_id')

    # Download original image
    image_bytes = download_image(s3_client, bucket, object_key)

    # Resize to thumbnail
    thumbnail_bytes = resize_image(image_bytes)

    # Upload thumbnail
    thumbnail_key = get_thumbnail_key(object_key)  # uploads/x.png -> thumbnails/x.png
    upload_thumbnail(s3_client, bucket, thumbnail_key, thumbnail_bytes)
```

**Key Patterns:**

- Batch processing with error handling for each message
- boto3 for S3-compatible operations
- Pillow for pure-Python image processing (no binary dependencies)
- Environment variables for AWS credentials

### Terraform Configuration

**[main.tf](examples/python/storage/tf/main.tf:1) - Function:**

```hcl
resource "yandex_function" "storage_handler" {
  name               = "storage-handler-python"
  runtime            = "python312"
  entrypoint         = "handler.handler"
  memory             = 256
  execution_timeout  = "30"

  environment = {
    AWS_ACCESS_KEY_ID     = yandex_iam_service_account_static_access_key.sa_storage_editor.access_key
    AWS_SECRET_ACCESS_KEY = yandex_iam_service_account_static_access_key.sa_storage_editor.secret_key
  }
}
```

**[main.tf](examples/python/storage/tf/main.tf:38) - Trigger:**

```hcl
resource "yandex_function_trigger" "storage_trigger" {
  object_storage {
    bucket_id    = yandex_storage_bucket.for_uploads.bucket
    prefix       = "uploads/"      # Only trigger on uploads/ folder
    create       = true
    update       = true
    batch_cutoff = 1               # Process immediately
  }

  function {
    id                 = yandex_function.storage_handler.id
    service_account_id = yandex_iam_service_account.trigger_sa.id
    retry_attempts     = 3
    retry_interval     = 10
  }
}
```

### Terratest Implementation

**[default_test.go](tests/python/storage/default_test.go:20):**

```go
func TestPythonStorageExample(t *testing.T) {
    // 1. Deploy infrastructure
    terraform.InitAndApply(t, terraformOptions)

    // 2. Setup S3 client with custom endpoint
    s3Client := s3.NewFromConfig(cfg, func(o *s3.Options) {
        o.EndpointResolverV2 = &resolverV2{}
        o.Credentials = credentials.NewStaticCredentialsProvider(accessKey, secretKey, "")
    })

    // 3. Cleanup (deferred)
    defer func() {
        s3Client.DeleteObject(ctx, &s3.DeleteObjectInput{
            Bucket: aws.String(bucket),
            Key:    aws.String("uploads/star.png"),
        })
        s3Client.DeleteObject(ctx, &s3.DeleteObjectInput{
            Bucket: aws.String(bucket),
            Key:    aws.String("thumbnails/star.png"),
        })
        terraform.Destroy(t, terraformOptions)
    }()

    // 4. Wait for IAM propagation
    time.Sleep(5 * time.Second)

    // 5. Upload test image
    file, _ := os.Open("star.png")
    s3Client.PutObject(ctx, &s3.PutObjectInput{
        Bucket:      aws.String(bucket),
        Key:         aws.String("uploads/star.png"),
        Body:        file,
        ContentType: aws.String("image/png"),
    })

    // 6. Wait for function processing
    time.Sleep(5 * time.Second)

    // 7. Verify thumbnail created
    _, err := s3Client.GetObject(ctx, &s3.GetObjectInput{
        Bucket: aws.String(bucket),
        Key:    aws.String("thumbnails/star.png"),
    })
    assert.NoError(t, err, "Failed to get resized image from bucket %s", bucket)
}
```

**Key Patterns:**

- Custom endpoint resolver for Yandex Cloud Storage
- Test image ([star.png](tests/python/storage/star.png)) in test directory
- Two 5-second waits (IAM propagation, function processing)
- Cleanup of both test objects before destroying infrastructure

## Step-by-Step Guide

### 1. Create Example Directory Structure

```bash
# Choose your language and example name
LANG=python
EXAMPLE=myexample

# Create directories
mkdir -p examples/$LANG/$EXAMPLE/{function,tf,environment}
mkdir -p tests/$LANG/$EXAMPLE

# Create placeholder files
touch examples/$LANG/$EXAMPLE/README.md
touch examples/$LANG/$EXAMPLE/function/handler.py  # or main.go, main.ts
touch examples/$LANG/$EXAMPLE/tf/{terraform.tf,variables.tf,outputs.tf,main.tf,iam.tf}
touch tests/$LANG/$EXAMPLE/default_test.go
```

### 2. Implement Function Handler

Choose the appropriate template from [Function Handler Patterns](#function-handler-patterns) and implement your business logic.

**Python Example:**

```python
# examples/python/myexample/function/handler.py
import logging
from typing import Dict, Any

logger = logging.getLogger()
logger.setLevel(logging.INFO)

def handler(event: Dict[str, Any], context: Any) -> Dict[str, int]:
    logger.info(f"Processing event: {event}")

    # Your implementation here

    return {"statusCode": 200}
```

**Add Dependencies:**

```bash
# For Python
echo "boto3>=1.34.0" > examples/python/myexample/function/requirements.txt

# For TypeScript
cd examples/typescript/myexample/function
npm init -y
npm install @yandex-cloud/function-types

# For Go
cd examples/go/myexample/function
go mod init example
```

### 3. Configure Terraform Infrastructure

**a. terraform.tf** - Use the template from [Terraform Patterns](#terraform-patterns)

**b. variables.tf** - Standard variables

```hcl
variable "cloud_id" {
  type = string
}

variable "folder_id" {
  type = string
}

variable "zone" {
  type    = string
  default = "ru-central1-a"
}
```

**c. main.tf** - Function and trigger

```hcl
data "archive_file" "function_code" {
  type        = "zip"
  source_dir  = "../function"
  output_path = "./function.zip"
}

resource "yandex_function" "example" {
  name               = "myexample-${var.lang}"
  runtime            = "python312"
  entrypoint         = "handler.handler"
  memory             = 128
  execution_timeout  = "10"
  user_hash          = data.archive_file.function_code.output_sha256
  service_account_id = yandex_iam_service_account.function_sa.id

  content {
    zip_filename = data.archive_file.function_code.output_path
  }

  depends_on = [yandex_resourcemanager_folder_iam_member.function_sa]
}
```

**d. iam.tf** - Service accounts

```hcl
resource "yandex_iam_service_account" "function_sa" {
  name      = "myexample-function-sa"
  folder_id = var.folder_id
}

resource "yandex_resourcemanager_folder_iam_member" "function_sa" {
  folder_id = var.folder_id
  role      = "editor"  # Adjust based on permissions needed
  member    = "serviceAccount:${yandex_iam_service_account.function_sa.id}"

  sleep_after = 5
}
```

**e. outputs.tf** - Test outputs

```hcl
output "function_id" {
  value = yandex_function.example.id
}

output "function_url" {
  value = "https://functions.yandexcloud.net/${yandex_function.example.id}"
}
```

**f. Create terraform.tfvars** (for local testing)

```hcl
cloud_id  = "your-cloud-id"
folder_id = "your-folder-id"
```

### 4. Create Terratest Test

**tests/{language}/{example}/default_test.go:**

```go
package myexample

import (
    "os"
    "testing"
    "time"

    "github.com/gruntwork-io/terratest/modules/terraform"
    "github.com/stretchr/testify/assert"
)

func TestMyExample(t *testing.T) {
    terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
        TerraformDir: "../../../examples/python/myexample/tf",
        Vars: map[string]interface{}{
            "cloud_id":  os.Getenv("CLOUD_ID"),
            "folder_id": os.Getenv("FOLDER_ID"),
        },
        EnvVars: map[string]string{
            "YC_TOKEN": os.Getenv("YC_TOKEN"),
        },
    })

    terraform.InitAndApply(t, terraformOptions)
    defer terraform.Destroy(t, terraformOptions)

    // Wait for IAM propagation
    time.Sleep(5 * time.Second)

    // Get outputs
    functionID := terraform.Output(t, terraformOptions, "function_id")

    // Assertions
    assert.NotEmpty(t, functionID, "Function ID should not be empty")

    // Additional test logic...
}
```

### 5. Document the Example

Create a comprehensive README.md following this structure:

```markdown
# {Example Title}

Brief description of what this example demonstrates.

## Architecture

```txt
Component A → Component B → Function → Component C
```

## Features

- Feature 1
- Feature 2
- Feature 3

## Prerequisites

- Yandex Cloud account
- Terraform >= 1.0
- {Language} runtime

## Project Structure

```txt
{example-name}/
├── README.md
├── function/
│   └── handler.{ext}
└── tf/
    ├── terraform.tf
    └── ...
```

## Deployment

### 1. Configure Variables

Create `terraform.tfvars`:

```hcl
cloud_id  = "your-cloud-id"
folder_id = "your-folder-id"
```

### 2. Deploy

```bash
cd tf
terraform init
terraform apply
```

## Usage

{Step-by-step usage instructions}

## How It Works

{Detailed explanation of the implementation}

## Testing

{How to verify the example works}

## Troubleshooting

{Common issues and solutions}

## Cleanup

```bash
terraform destroy
```

```txt

### 6. Test and Verify

**a. Local Terraform Test:**
```bash
cd examples/{language}/{example}/tf
terraform init
terraform plan
terraform apply
# ... test manually ...
terraform destroy
```

**b. Run Terratest:**

```bash
cd tests/{language}/{example}
export CLOUD_ID=your-cloud-id
export FOLDER_ID=your-folder-id
export YC_TOKEN=your-token

go test -v -timeout 30m
```

**c. Verify Test Output:**

- Infrastructure deployed successfully
- Function invoked correctly
- Expected outputs generated
- Resources cleaned up

## Common Patterns and Conventions

### Naming Conventions

**Resources:**

- Use kebab-case: `storage-handler-python`, `yds-producer-go`
- Include language suffix: `-python`, `-go`, `-ts`
- Descriptive names: `{service}-{role}-{language}`

**Prefixes/Paths:**

- Lowercase with trailing slash: `uploads/`, `thumbnails/`
- Consistent across Terraform and code

**Test Files:**

- Always `default_test.go` for main integration test
- Optional `infrastructure_test.go` for infrastructure-only validation

**Package Names:**

- Match example name: `package storage`, `package apigw`

### IAM Patterns

**Service Account Roles:**

- **Function execution**: Specific roles based on services accessed (`storage.editor`, `ydb.viewer`)
- **Trigger invocation**: Always `functions.functionInvoker`
- **Separate accounts**: Different service accounts for different roles

**IAM Propagation:**

- Always include `sleep_after = 5` in IAM bindings
- Wait 5 seconds in tests after `terraform apply` before testing
- Critical for avoiding permission errors

### Environment Variables

**Standard Pattern:**

```hcl
environment = {
  AWS_ACCESS_KEY_ID     = yandex_iam_service_account_static_access_key.sa.access_key
  AWS_SECRET_ACCESS_KEY = yandex_iam_service_account_static_access_key.sa.secret_key
  # Service-specific variables
  YDS_TOPIC_ID         = yandex_ydb_topic.example.id
  DATABASE_PATH        = yandex_ydb_database_serverless.example.database_path
}
```

**Test Environment:**

```go
EnvVars: map[string]string{
    "YC_TOKEN": os.Getenv("YC_TOKEN"),
},
```

### Error Handling

**Python:**

```python
try:
    process_message(s3_client, message)
except Exception as e:
    logger.error(f"Error processing message: {e}")
    # Continue processing other messages
```

**Go:**

```go
if err := processEvent(event); err != nil {
    return Response{StatusCode: 500}, fmt.Errorf("process event: %w", err)
}
```

**TypeScript:**

```typescript
try {
    await processEvent(event);
} catch (error) {
    console.error('Error:', error);
    return { statusCode: 500 };
}
```

### Trigger Batch Configuration

**Immediate Processing (Storage):**

```hcl
object_storage {
  batch_cutoff = 1  # Process immediately
}
```

**Batch Processing (YDS/Topics):**

```hcl
data_streams {
  batch_size   = 10     # Process up to 10 messages
  batch_cutoff = "5"    # Or wait 5 seconds
}
```

## Testing Best Practices

### Test Organization

**File Structure:**

```txt
tests/{language}/{example}/
├── default_test.go          # Main integration test
├── infrastructure_test.go   # Infrastructure validation (optional)
└── test_data/              # Test files (images, JSON, etc.)
    └── star.png
```

### Terratest Options

**Standard Configuration:**

```go
terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
    TerraformDir: "../../../examples/{language}/{example}/tf",
    Vars: map[string]interface{}{
        "cloud_id":  os.Getenv("CLOUD_ID"),
        "folder_id": os.Getenv("FOLDER_ID"),
    },
    EnvVars: map[string]string{
        "YC_TOKEN": os.Getenv("YC_TOKEN"),
    },
})
```

**With Timeout:**

```go
terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
    // ... other options ...
    MaxRetries:         3,
    TimeBetweenRetries: 5 * time.Second,
})
```

### Assertions

**Standard Patterns:**

```go
// Not empty
assert.NotEmpty(t, functionID, "Function ID should not be empty")

// No error
assert.NoError(t, err, "Failed to upload file %s to bucket %s", filename, bucket)

// Equal
assert.Equal(t, expected, actual, "Values should match")

// Contains
assert.Contains(t, response, "expected-substring", "Response should contain substring")
```

### Resource Cleanup

**Defer Pattern:**

```go
defer func() {
    // 1. Clean up test-created resources (S3 objects, database records)
    cleanupTestResources(ctx, client)

    // 2. Destroy infrastructure
    terraform.Destroy(t, terraformOptions)
}()
```

**S3 Cleanup Example:**

```go
defer func() {
    _, _ = s3Client.DeleteObject(ctx, &s3.DeleteObjectInput{
        Bucket: aws.String(bucket),
        Key:    aws.String("uploads/test.png"),
    })
    _, _ = s3Client.DeleteObject(ctx, &s3.DeleteObjectInput{
        Bucket: aws.String(bucket),
        Key:    aws.String("thumbnails/test.png"),
    })

    terraform.Destroy(t, terraformOptions)
}()
```

### Wait Times and Synchronization

**IAM Propagation:**

```go
// After terraform.InitAndApply
time.Sleep(5 * time.Second)
```

**Event Processing:**

```go
// After uploading file or triggering event
time.Sleep(5 * time.Second)  // Wait for trigger and function execution
```

**Polling Alternative:**

```go
// For more robust testing
err := retry.DoWithRetry(t, "Check for result", 10, 3*time.Second, func() (string, error) {
    result, err := checkResult()
    if err != nil {
        return "", err
    }
    return "Success", nil
})
```

## Troubleshooting

### IAM Permission Errors

**Symptom:**

```txt
Error: AccessDenied: Access Denied
```

**Solutions:**

1. **Wait for propagation**: Add 5-15 seconds after `terraform apply`
2. **Check IAM roles**: Verify service account has correct roles
3. **Verify sleep_after**: Ensure `sleep_after = 5` in IAM bindings
4. **Check depends_on**: Function should depend on IAM member resources

### Terraform State Issues

**Symptom:**

```txt
Error: state file locked
```

**Solutions:**

```bash
# Remove lock file
rm examples/{language}/{example}/environment/.terraform.tfstate.lock.info

# Force unlock
terraform force-unlock <lock-id>
```

### Missing Environment Variables

**Symptom:**

```txt
Error: cloud_id is required
```

**Solutions:**

```bash
# Set environment variables
export CLOUD_ID=your-cloud-id
export FOLDER_ID=your-folder-id
export YC_TOKEN=your-token

# Or create terraform.tfvars
cat > examples/{language}/{example}/tf/terraform.tfvars <<EOF
cloud_id  = "your-cloud-id"
folder_id = "your-folder-id"
EOF
```

### S3 Endpoint Configuration

**Symptom:**

```txt
Error: could not resolve endpoint
```

**Solution:**
Ensure custom endpoint resolver is configured:

```go
s3Client := s3.NewFromConfig(cfg, func(o *s3.Options) {
    o.EndpointResolverV2 = &resolverV2{}  // Custom resolver
    o.Region = "ru-central1"
})
```

### Function Timeout

**Symptom:**

```txt
Error: Function execution timed out
```

**Solutions:**

1. **Increase timeout** in main.tf:

   ```hcl
   execution_timeout = "60"  # Increase from 10 to 60 seconds
   ```

2. **Increase memory** (more memory = more CPU):

   ```hcl
   memory = 512  # Increase from 128 to 512 MB
   ```

### Test Timeouts

**Symptom:**

```txt
panic: test timed out after 10m0s
```

**Solution:**

```bash
# Increase test timeout
go test -v -timeout 30m
```

### Archive Provider Excludes

**Symptom:**
Function package is too large or includes unnecessary files.

**Solution:**

```hcl
data "archive_file" "function_code" {
  type        = "zip"
  source_dir  = "../function"
  output_path = "./function.zip"
  excludes = [
    "__pycache__",
    "*.pyc",
    ".pytest_cache",
    "node_modules",
    ".git",
    "*.test",
  ]
}
```

## References

### Yandex Cloud Documentation

- [Cloud Functions](https://cloud.yandex.com/docs/functions/)
- [Object Storage](https://cloud.yandex.com/docs/storage/)
- [YDB (Data Streams)](https://cloud.yandex.com/docs/ydb/)
- [Triggers](https://cloud.yandex.com/docs/functions/concepts/trigger/)
- [IAM](https://cloud.yandex.com/docs/iam/)

### Runtime Documentation

- [Python Runtime](https://cloud.yandex.com/docs/functions/lang/python/)
- [Go Runtime](https://cloud.yandex.com/docs/functions/lang/golang/)
- [Node.js Runtime](https://cloud.yandex.com/docs/functions/lang/nodejs/)

### Testing Tools

- [Terratest](https://terratest.gruntwork.io/)
- [Testify](https://github.com/stretchr/testify)
- [AWS SDK for Go v2](https://aws.github.io/aws-sdk-go-v2/docs/)

### Example References

- [Python Storage Example](examples/python/storage/) - Image thumbnail generation
- [Go YDS Example](examples/go/yds/) - Data Streams producer/consumer
- [TypeScript API Gateway Example](examples/typescript/apigw/) - HTTP API endpoint

---

**Questions or Issues?**

If you encounter problems or have questions about adding examples:

1. Review existing examples in the repository for patterns
2. Check the [python/storage](examples/python/storage) example as a reference
3. Ensure all prerequisites are met (environment variables, Terraform version)
4. Review error logs carefully for IAM propagation or configuration issues

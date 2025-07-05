# Go YDB Example

This example demonstrates how to use Yandex Database (YDB) with Go functions in Yandex Cloud Functions.

## Features

- Serverless YDB database
- Go function with YDB SDK integration
- Automatic IAM permissions setup
- Environment-based configuration

## Prerequisites

- Yandex Cloud CLI configured
- Terraform installed
- Go 1.23+ for local development

## Deployment

1. Navigate to the `tf` directory:
   ```bash
   cd examples/go/ydb/tf
   ```

2. Initialize Terraform:
   ```bash
   terraform init
   ```

3. Create a `terraform.tfvars` file with your Yandex Cloud credentials:
   ```hcl
   cloud_id  = "your-cloud-id"
   folder_id = "your-folder-id"
   zone      = "ru-central1-a"
   ```

4. Deploy the infrastructure:
   ```bash
   terraform apply
   ```

## Usage

After deployment, you need to set up the database schema:

1. Connect to your YDB database using the YDB CLI or web console
2. Run the setup script:
   ```bash
   ydb -e $(terraform output -raw ydb_endpoint) -d $(terraform output -raw ydb_database_path) -f ddl.sql
   ```

3. Invoke the function using the provided URL:
   ```bash
   curl $(terraform output -raw function_url)
   ```

The function will return user data for ID 3 (Charlie).

## Infrastructure

The Terraform configuration creates:

- **YDB Serverless Database**: A serverless YDB instance
- **Service Account**: With YDB viewer and editor permissions
- **Cloud Function**: Go function with YDB SDK integration
- **IAM Binding**: Makes the function publicly accessible

## Environment Variables

The function uses the following environment variables:

- `YDB_ENDPOINT`: YDB server endpoint
- `YDB_DATABASE`: YDB database path

These are automatically set by Terraform during deployment.

## Local Development

To run the function locally:

1. Set environment variables:
   ```bash
   export YDB_ENDPOINT="your-ydb-endpoint"
   export YDB_DATABASE="your-database-path"
   ```

2. Run the function:
   ```bash
   go run function/main.go
   ```

## Cleanup

To destroy the infrastructure:

```bash
terraform destroy
```

## Notes

- The function expects a `users` table with `id` and `name` columns
- The example queries for user with ID 3
- Make sure to create the required table structure in your YDB database 
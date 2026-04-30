# TypeScript YDB Query with Transaction Statistics

This example demonstrates how to query a YDB Serverless database using the `@ydbjs/query` package and retrieve execution statistics to estimate read unit consumption.

## Architecture

```
HTTP Request → Cloud Function → YDB Serverless Database
                    ↓
              Query with StatsMode.FULL
                    ↓
              Response: { users, stats }
```

## Features

- Query YDB using `@ydbjs/query` tagged template literals
- Collect transaction statistics via `StatsMode.FULL`
- Return read rows, bytes, and CPU time metrics
- Metadata-based authentication for Cloud Functions

## Prerequisites

- Yandex Cloud account
- Terraform >= 1.0
- Node.js >= 22
- YDB CLI (for database setup)

## Project Structure

```
ydb/
├── README.md
├── ddl.sql                  # Table schema
├── dml.sql                  # Sample data
├── function/
│   ├── main.ts              # HTTP handler with query stats
│   ├── package.json
│   └── tsconfig.json
└── tf/
    ├── terraform.tf
    ├── variables.tf
    ├── ydb.tf               # YDB Serverless database
    ├── iam.tf               # Service account and roles
    ├── main.tf              # Function and build
    └── outputs.tf
```

## Deployment

### 1. Configure Variables

Create `tf/terraform.tfvars`:

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

### 3. Set Up Database

```bash
ydb -e <ydb-endpoint> -d <database-path> sql -f ../ddl.sql
ydb -e <ydb-endpoint> -d <database-path> sql -f ../dml.sql
```

## Usage

Invoke the function via HTTP:

```bash
curl https://functions.yandexcloud.net/<function-id>
```

Response:

```json
{
  "users": [
    { "id": 1, "name": "Alice" },
    { "id": 2, "name": "Bob" }
  ],
  "stats": {
    "queryPhases": [
      {
        "tableAccess": [
          {
            "name": "users",
            "reads": { "rows": 5, "bytes": 120 }
          }
        ],
        "cpuTimeUs": 1234
      }
    ],
    "processCpuTimeUs": 5678
  }
}
```

The `stats` object contains detailed execution metrics from YDB's query service, including per-phase table access statistics with read row/byte counts useful for estimating read unit consumption.

## How It Works

1. The function receives an HTTP request
2. Connects to YDB using `@ydbjs/core` Driver with metadata credentials
3. Creates a query client via `query(driver)` from `@ydbjs/query`
4. Executes `SELECT id, name FROM users` with `.withStats(StatsMode.FULL)`
5. Retrieves statistics via `.stats()` on the query object
6. Returns both the query results and execution statistics as JSON

## Cleanup

```bash
cd tf
terraform destroy
```

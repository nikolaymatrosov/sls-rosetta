# YDB Example Tests

This directory contains comprehensive tests for the Go/YDB example.

## Test Structure

- **`default_test.go`**: Main integration test that deploys infrastructure and tests the function
- **`infrastructure_test.go`**: Tests infrastructure outputs and validation
- **`function_test.go`**: Unit tests for function behavior and error handling

## Prerequisites

- Go 1.21+
- Terraform
- Yandex Cloud CLI configured
- YDB CLI (optional, for database setup)

## Environment Variables

Set the following environment variables before running tests:

```bash
export CLOUD_ID="your-cloud-id"
export FOLDER_ID="your-folder-id"
export YC_TOKEN="your-yandex-cloud-token"
```

## Running Tests

### Run all tests
```bash
go test -v
```

### Run specific test
```bash
go test -v -run TestGoYdbExample
go test -v -run TestYdbInfrastructure
go test -v -run TestFunctionHandler
```

### Run tests with timeout
```bash
go test -v -timeout 30m
```

### Run tests in parallel (for unit tests only)
```bash
go test -v -parallel 4
```

## Test Categories

### Integration Tests
- **TestGoYdbExample**: Deploys full infrastructure, sets up database, and tests function
- **TestYdbInfrastructure**: Validates Terraform outputs and infrastructure configuration

### Unit Tests
- **TestFunctionHandler**: Tests function error handling with different environment configurations
- **TestUserStruct**: Tests JSON marshaling/unmarshaling of the User struct
- **TestFunctionResponseHeaders**: Tests HTTP response headers

### Benchmark Tests
- **BenchmarkFunctionHandler**: Performance benchmark for the function

## Test Behavior

### Integration Test Flow
1. Deploys YDB serverless database
2. Creates service account with YDB permissions
3. Deploys Go function with YDB SDK
4. Sets up database schema (if YDB CLI available)
5. Tests function by calling it and verifying response
6. Cleans up infrastructure

### Error Handling Tests
- Missing environment variables
- Invalid YDB connection
- JSON parsing errors
- HTTP response validation

## Notes

- Tests will skip YDB CLI setup if the CLI is not available
- Integration tests require valid Yandex Cloud credentials
- Tests automatically clean up infrastructure after completion
- Some tests may take several minutes to complete due to infrastructure deployment

## Troubleshooting

### Common Issues

1. **Authentication errors**: Ensure YC_TOKEN is valid and has sufficient permissions
2. **Resource limits**: Check your Yandex Cloud quota for YDB and Functions
3. **Network issues**: Ensure your environment can reach Yandex Cloud APIs
4. **YDB CLI not found**: Tests will continue without database setup

### Debug Mode

Run tests with verbose output:
```bash
go test -v -debug
```

### Cleanup

If tests fail and leave infrastructure behind, manually clean up:
```bash
cd ../../../examples/go/ydb/tf
terraform destroy
``` 
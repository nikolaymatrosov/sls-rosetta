#!/bin/bash

# YDB Example Test Runner
# This script runs the YDB example tests with proper setup and cleanup

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check prerequisites
check_prerequisites() {
    print_status "Checking prerequisites..."
    
    # Check if Go is installed
    if ! command -v go &> /dev/null; then
        print_error "Go is not installed or not in PATH"
        exit 1
    fi
    
    # Check if Terraform is installed
    if ! command -v terraform &> /dev/null; then
        print_error "Terraform is not installed or not in PATH"
        exit 1
    fi
    
    # Check environment variables
    if [ -z "$CLOUD_ID" ]; then
        print_error "CLOUD_ID environment variable is not set"
        exit 1
    fi
    
    if [ -z "$FOLDER_ID" ]; then
        print_error "FOLDER_ID environment variable is not set"
        exit 1
    fi
    
    if [ -z "$YC_TOKEN" ]; then
        print_error "YC_TOKEN environment variable is not set"
        exit 1
    fi
    
    # Check if YDB CLI is available (optional)
    if ! command -v ydb &> /dev/null; then
        print_warning "YDB CLI not found. Database setup will be skipped."
    fi
    
    print_status "Prerequisites check passed"
}

# Run tests
run_tests() {
    local test_type=$1
    
    print_status "Running tests..."
    
    case $test_type in
        "unit")
            print_status "Running unit tests only..."
            go test -v -run "TestFunctionHandler|TestUserStruct|TestFunctionResponseHeaders|BenchmarkFunctionHandler"
            ;;
        "infrastructure")
            print_status "Running infrastructure tests only..."
            go test -v -run "TestYdbInfrastructure|TestYdbInfrastructureVariables"
            ;;
        "integration")
            print_status "Running integration tests only..."
            go test -v -run "TestGoYdbExample" -timeout 30m
            ;;
        "all"|"")
            print_status "Running all tests..."
            go test -v -timeout 30m
            ;;
        *)
            print_error "Unknown test type: $test_type"
            print_status "Available test types: unit, infrastructure, integration, all"
            exit 1
            ;;
    esac
}

# Main execution
main() {
    local test_type=${1:-"all"}
    
    print_status "Starting YDB example tests..."
    
    # Check prerequisites
    check_prerequisites
    
    # Run tests
    run_tests $test_type
    
    print_status "Tests completed successfully!"
}

# Handle script arguments
case "${1:-}" in
    "unit"|"infrastructure"|"integration"|"all"|"")
        main "$1"
        ;;
    "help"|"-h"|"--help")
        echo "Usage: $0 [test_type]"
        echo ""
        echo "Test types:"
        echo "  unit           - Run unit tests only"
        echo "  infrastructure - Run infrastructure tests only"
        echo "  integration    - Run integration tests only"
        echo "  all            - Run all tests (default)"
        echo "  help           - Show this help message"
        echo ""
        echo "Environment variables required:"
        echo "  CLOUD_ID   - Yandex Cloud ID"
        echo "  FOLDER_ID  - Yandex Cloud folder ID"
        echo "  YC_TOKEN   - Yandex Cloud authentication token"
        ;;
    *)
        print_error "Unknown argument: $1"
        echo "Use '$0 help' for usage information"
        exit 1
        ;;
esac 
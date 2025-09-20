#!/bin/bash

set -euo pipefail

echo "Building Lambda functions..."

# Get to the right directory
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )"
cd "$SCRIPT_DIR/.."

# Lambda functions to build
FUNCTIONS=(
    "users"
    "accounts" 
    "transactions"
    "categories"
    "budgets"
    "reports"
    "auth"
    "notifications"
)

# Build each Go Lambda function
for func in "${FUNCTIONS[@]}"; do
    echo "Building $func function..."
    
    cd "lambda/$func"
    
    # Set Go environment
    export GOOS=linux
    export GOARCH=amd64
    export CGO_ENABLED=0
    export GO111MODULE=on
    
    # Build the function
    go mod tidy
    go build -ldflags="-s -w" -o bootstrap main.go
    
    # Verify the binary exists
    if [[ ! -f bootstrap ]]; then
        echo "Error: Failed to build $func function"
        exit 1
    fi
    
    echo "Successfully built $func function"
    cd - > /dev/null
done

# Build Node.js authorizer function
echo "Building authorizer function..."
cd lambda/authorizer

# Install dependencies and build
npm ci --production

echo "Successfully built authorizer function"
cd - > /dev/null

# Build database initialization function  
echo "Building database init function..."
cd lambda/db-init

npm ci --production

echo "Successfully built database init function"
cd - > /dev/null

echo "All Lambda functions built successfully!"
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

# Build each Go Lambda function from backend directory
for func in "${FUNCTIONS[@]}"; do
    echo "Building $func function..."
    
    # Set Go environment
    export GOOS=linux
    export GOARCH=amd64
    export CGO_ENABLED=0
    export GO111MODULE=on
    
    # Check if source exists first
    if [[ -f "lambda/$func/main.go" ]]; then
        # Copy source to backend temporary directory and build there
        temp_dir="../backend/tmp/lambda/$func"
        mkdir -p "$temp_dir"
        cp "lambda/$func/main.go" "$temp_dir/"
        
        # Build from backend directory to use financetracker module
        cd "../backend"
        
        # Build the function and output to CDK lambda directory
        if go build -ldflags="-s -w" -o "../cdk/lambda/$func/bootstrap" "./tmp/lambda/$func/main.go"; then
            echo "Successfully built $func function"
        else
            echo "Error: Failed to build $func function"
            exit 1
        fi
        
        # Cleanup temporary files
        rm -rf "./tmp/lambda/$func"
        
        cd "../cdk"
    else
        echo "Skipping $func function (source not found)"
    fi
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
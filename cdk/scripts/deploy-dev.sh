#!/bin/bash

set -euo pipefail

echo "=== FinSight Development Deployment ==="
echo "Date: $(date)"
echo "Environment: development"
echo ""

# Check AWS credentials
echo "Checking AWS credentials..."
aws sts get-caller-identity > /dev/null 2>&1 || {
    echo "Error: AWS credentials not configured"
    exit 1
}

# Build TypeScript
echo "Building TypeScript..."
npm run build

# Build Lambda functions
echo "Building Lambda functions..."
./scripts/build-lambda.sh

# Run CDK synth to verify
echo "Running CDK synthesis..."
npx cdk synth --context env=dev

# Deploy all stacks
echo "Deploying to development environment..."
npx cdk deploy --all --context env=dev --require-approval never

echo ""
echo "=== Deployment Complete ==="
echo "API Endpoint: $(aws cloudformation describe-stacks --stack-name ApiStack-dev --query 'Stacks[0].Outputs[?OutputKey==`ApiEndpoint`].OutputValue' --output text)"
echo "Dashboard URL: https://ap-northeast-1.console.aws.amazon.com/cloudwatch/home?region=ap-northeast-1#dashboards:name=FinSight-dev"
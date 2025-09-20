#!/bin/bash

set -euo pipefail

echo "=== FinSight Production Deployment ==="
echo "Date: $(date)"
echo "Environment: production"
echo ""

# Confirmation prompt
read -p "Are you sure you want to deploy to PRODUCTION? (yes/no): " confirm
if [ "$confirm" != "yes" ]; then
    echo "Deployment cancelled"
    exit 0
fi

# Check AWS credentials
echo "Checking AWS credentials..."
aws sts get-caller-identity > /dev/null 2>&1 || {
    echo "Error: AWS credentials not configured"
    exit 1
}

# Run tests first
echo "Running tests..."
npm test

# Build TypeScript
echo "Building TypeScript..."
npm run build

# Build Lambda functions
echo "Building Lambda functions..."
./scripts/build-lambda.sh

# Run CDK synth to verify
echo "Running CDK synthesis..."
npx cdk synth --context env=prod

# Create backup of current deployment
echo "Creating deployment backup..."
aws cloudformation describe-stacks \
    --query 'Stacks[?contains(StackName, `prod`)].{Name:StackName,Status:StackStatus,UpdateTime:LastUpdatedTime}' \
    --output table > deployment-backup-$(date +%Y%m%d-%H%M%S).txt

# Deploy all stacks
echo "Deploying to production environment..."
npx cdk deploy --all --context env=prod --require-approval never

# Run post-deployment verification
echo "Running post-deployment verification..."
API_ENDPOINT=$(aws cloudformation describe-stacks --stack-name ApiStack-prod --query 'Stacks[0].Outputs[?OutputKey==`ApiEndpoint`].OutputValue' --output text)
curl -s -o /dev/null -w "%{http_code}" "$API_ENDPOINT/v1/health" | grep -q "200" || {
    echo "Warning: Health check failed"
}

echo ""
echo "=== Deployment Complete ==="
echo "API Endpoint: $API_ENDPOINT"
echo "Dashboard URL: https://ap-northeast-1.console.aws.amazon.com/cloudwatch/home?region=ap-northeast-1#dashboards:name=FinSight-prod"
echo ""
echo "Please monitor the CloudWatch dashboard for any issues."
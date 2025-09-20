#!/bin/bash

set -euo pipefail

ENVIRONMENT="${1:-dev}"

echo "=== Testing $ENVIRONMENT deployment ==="

# Get stack outputs
API_ENDPOINT=$(aws cloudformation describe-stacks --stack-name "ApiStack-$ENVIRONMENT" --query 'Stacks[0].Outputs[?OutputKey==`ApiEndpoint`].OutputValue' --output text)

if [ -z "$API_ENDPOINT" ]; then
    echo "Error: Could not get API endpoint"
    exit 1
fi

echo "API Endpoint: $API_ENDPOINT"
echo ""

# Test health endpoint
echo "Testing health endpoint..."
HEALTH_STATUS=$(curl -s -o /dev/null -w "%{http_code}" "$API_ENDPOINT/v1/health")
if [ "$HEALTH_STATUS" == "200" ]; then
    echo "✓ Health check passed"
else
    echo "✗ Health check failed (Status: $HEALTH_STATUS)"
    exit 1
fi

# Test auth endpoint
echo "Testing auth endpoint..."
AUTH_STATUS=$(curl -s -o /dev/null -w "%{http_code}" -X POST "$API_ENDPOINT/v1/auth/callback")
if [ "$AUTH_STATUS" == "400" ] || [ "$AUTH_STATUS" == "401" ]; then
    echo "✓ Auth endpoint accessible"
else
    echo "✗ Auth endpoint test failed (Status: $AUTH_STATUS)"
fi

# Check Lambda functions
echo ""
echo "Checking Lambda functions..."
aws lambda list-functions --query "Functions[?contains(FunctionName, 'ApiStack-$ENVIRONMENT')].{Name:FunctionName,Runtime:Runtime,State:State}" --output table

# Check RDS database
echo ""
echo "Checking RDS database..."
DB_STATUS=$(aws rds describe-db-instances --query "DBInstances[?contains(DBInstanceIdentifier, '$ENVIRONMENT')].{ID:DBInstanceIdentifier,Status:DBInstanceStatus,Engine:Engine}" --output text)
if [ -n "$DB_STATUS" ]; then
    echo "✓ Database found: $DB_STATUS"
else
    echo "✗ Database not found"
fi

# Check CloudWatch alarms
echo ""
echo "Checking CloudWatch alarms..."
ALARM_COUNT=$(aws cloudwatch describe-alarms --query "MetricAlarms[?contains(AlarmName, '$ENVIRONMENT')] | length(@)")
echo "✓ Found $ALARM_COUNT alarms configured"

echo ""
echo "=== Deployment test complete ==="
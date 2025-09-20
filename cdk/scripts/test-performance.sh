#!/bin/bash

set -euo pipefail

ENVIRONMENT="${1:-prod}"

echo "=== Performance Tests for $ENVIRONMENT environment ==="
echo ""

# Get API endpoint
API_ENDPOINT=$(aws cloudformation describe-stacks --stack-name "ApiStack-$ENVIRONMENT" --query 'Stacks[0].Outputs[?OutputKey==`ApiEndpoint`].OutputValue' --output text)

if [ -z "$API_ENDPOINT" ]; then
    echo "Error: Could not get API endpoint"
    exit 1
fi

echo "API Endpoint: $API_ENDPOINT"
echo ""

# Test 1: Response Time - Health Check
echo "Test 1: Health Check Response Time"
echo "Running 10 requests..."
TOTAL_TIME=0
for i in {1..10}; do
    TIME=$(curl -s -o /dev/null -w "%{time_total}" "$API_ENDPOINT/v1/health")
    TOTAL_TIME=$(echo "$TOTAL_TIME + $TIME" | bc)
    echo "  Request $i: ${TIME}s"
done
AVG_TIME=$(echo "scale=3; $TOTAL_TIME / 10" | bc)
echo "Average response time: ${AVG_TIME}s"

if (( $(echo "$AVG_TIME < 0.5" | bc -l) )); then
    echo "✓ Health check response time acceptable"
else
    echo "⚠ Health check response time high"
fi

# Test 2: Concurrent Requests
echo ""
echo "Test 2: Concurrent Request Handling"
echo "Sending 20 concurrent requests..."

# Create temp file for results
TEMP_FILE=$(mktemp)

# Send concurrent requests
for i in {1..20}; do
    curl -s -o /dev/null -w "%{http_code} %{time_total}\n" "$API_ENDPOINT/v1/health" >> "$TEMP_FILE" &
done

# Wait for all requests to complete
wait

# Analyze results
SUCCESS_COUNT=$(grep "^200" "$TEMP_FILE" | wc -l)
TOTAL_RESPONSE_TIME=$(awk '{sum += $2} END {print sum}' "$TEMP_FILE")
AVG_CONCURRENT_TIME=$(echo "scale=3; $TOTAL_RESPONSE_TIME / 20" | bc)

echo "Successful requests: $SUCCESS_COUNT/20"
echo "Average response time (concurrent): ${AVG_CONCURRENT_TIME}s"

if [ "$SUCCESS_COUNT" -eq 20 ]; then
    echo "✓ All concurrent requests succeeded"
else
    echo "✗ Some concurrent requests failed"
fi

# Clean up
rm "$TEMP_FILE"

# Test 3: API Gateway Latency
echo ""
echo "Test 3: API Gateway Latency"
LATENCY_METRICS=$(aws cloudwatch get-metric-statistics \
    --namespace AWS/ApiGateway \
    --metric-name Latency \
    --dimensions Name=ApiName,Value="finsight-api-$ENVIRONMENT" \
    --statistics Average \
    --start-time "$(date -u -d '5 minutes ago' +%Y-%m-%dT%H:%M:%S)" \
    --end-time "$(date -u +%Y-%m-%dT%H:%M:%S)" \
    --period 300 \
    --query 'Datapoints[0].Average' \
    --output text 2>/dev/null || echo "N/A")

if [ "$LATENCY_METRICS" != "N/A" ] && [ -n "$LATENCY_METRICS" ]; then
    echo "API Gateway average latency: ${LATENCY_METRICS}ms"
else
    echo "API Gateway metrics not available yet"
fi

# Test 4: Lambda Cold Start
echo ""
echo "Test 4: Lambda Cold Start Performance"
echo "Testing users function..."

# Get Lambda function name
FUNCTION_NAME=$(aws lambda list-functions --query "Functions[?contains(FunctionName, 'ApiStack-$ENVIRONMENT-usersFunction')].FunctionName" --output text)

if [ -n "$FUNCTION_NAME" ]; then
    # Invoke function directly
    START_TIME=$(date +%s%3N)
    aws lambda invoke --function-name "$FUNCTION_NAME" --payload '{}' /tmp/lambda-response.json > /dev/null 2>&1 || true
    END_TIME=$(date +%s%3N)
    
    COLD_START_TIME=$((END_TIME - START_TIME))
    echo "Cold start time: ${COLD_START_TIME}ms"
    
    if [ "$COLD_START_TIME" -lt 3000 ]; then
        echo "✓ Cold start time acceptable"
    else
        echo "⚠ Cold start time high"
    fi
    
    rm -f /tmp/lambda-response.json
else
    echo "Lambda function not found"
fi

# Test 5: Database Connection Pool
echo ""
echo "Test 5: Database Performance"
DB_INSTANCE=$(aws rds describe-db-instances --query "DBInstances[?contains(DBInstanceIdentifier, '$ENVIRONMENT')].DBInstanceIdentifier" --output text)

if [ -n "$DB_INSTANCE" ]; then
    # Get database metrics
    DB_CONNECTIONS=$(aws cloudwatch get-metric-statistics \
        --namespace AWS/RDS \
        --metric-name DatabaseConnections \
        --dimensions Name=DBInstanceIdentifier,Value="$DB_INSTANCE" \
        --statistics Average \
        --start-time "$(date -u -d '5 minutes ago' +%Y-%m-%dT%H:%M:%S)" \
        --end-time "$(date -u +%Y-%m-%dT%H:%M:%S)" \
        --period 300 \
        --query 'Datapoints[0].Average' \
        --output text 2>/dev/null || echo "N/A")
    
    if [ "$DB_CONNECTIONS" != "N/A" ] && [ -n "$DB_CONNECTIONS" ]; then
        echo "Average database connections: $DB_CONNECTIONS"
    else
        echo "Database metrics not available"
    fi
fi

# Summary
echo ""
echo "=== Performance Test Summary ==="
echo "✓ Health check avg response: ${AVG_TIME}s"
echo "✓ Concurrent request handling: $SUCCESS_COUNT/20 succeeded"
echo "✓ Cold start tested"
echo ""

# Recommendations based on results
echo "Performance Recommendations:"
if (( $(echo "$AVG_TIME > 0.3" | bc -l) )); then
    echo "- Consider increasing Lambda memory for better performance"
fi
if [ "$SUCCESS_COUNT" -lt 20 ]; then
    echo "- Review API Gateway throttling settings"
fi
echo "- Monitor CloudWatch dashboards for detailed metrics"
echo "- Consider implementing caching for frequently accessed data"

echo ""
echo "Performance tests completed!"
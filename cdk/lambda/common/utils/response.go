package utils

import (
    "encoding/json"
    "github.com/aws/aws-lambda-go/events"
)

// CreateResponse creates an API Gateway response
func CreateResponse(statusCode int, body interface{}) (events.APIGatewayProxyResponse, error) {
    headers := map[string]string{
        "Content-Type":                "application/json",
        "Access-Control-Allow-Origin": "*",
        "Access-Control-Allow-Headers": "Content-Type,Authorization",
        "Access-Control-Allow-Methods": "GET,POST,PUT,DELETE,OPTIONS",
    }
    
    var bodyString string
    if body != nil {
        bodyBytes, err := json.Marshal(body)
        if err != nil {
            return events.APIGatewayProxyResponse{
                StatusCode: 500,
                Headers:    headers,
                Body:       `{"error":"Failed to marshal response"}`,
            }, err
        }
        bodyString = string(bodyBytes)
    }
    
    return events.APIGatewayProxyResponse{
        StatusCode: statusCode,
        Headers:    headers,
        Body:       bodyString,
    }, nil
}

// ErrorResponse creates an error response
func ErrorResponse(statusCode int, message string) events.APIGatewayProxyResponse {
    resp, _ := CreateResponse(statusCode, map[string]string{
        "error": message,
    })
    return resp
}

// SuccessResponse creates a success response
func SuccessResponse(data interface{}) events.APIGatewayProxyResponse {
    resp, _ := CreateResponse(200, data)
    return resp
}
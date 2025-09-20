package main

import (
    "context"
    "encoding/json"
    "fmt"
    "log"
    
    "github.com/aws/aws-lambda-go/events"
    "github.com/aws/aws-lambda-go/lambda"
    
    "finsight/common/db"
    "finsight/common/models"
    "finsight/common/utils"
)

func handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
    // Get user ID from authorizer context
    userID, ok := request.RequestContext.Authorizer["userId"].(string)
    if !ok {
        return utils.ErrorResponse(401, "Unauthorized"), nil
    }
    
    switch request.HTTPMethod {
    case "GET":
        return handleGetUser(userID)
    case "PUT":
        return handleUpdateUser(userID, request.Body)
    default:
        return utils.ErrorResponse(405, "Method not allowed"), nil
    }
}

func handleGetUser(auth0UserID string) (events.APIGatewayProxyResponse, error) {
    database, err := db.GetDB()
    if err != nil {
        log.Printf("Failed to connect to database: %v", err)
        return utils.ErrorResponse(500, "Database connection error"), nil
    }
    
    var user models.User
    query := `SELECT id, auth0_user_id, email, name, created_at, updated_at 
              FROM users WHERE auth0_user_id = $1`
    
    err = database.QueryRow(query, auth0UserID).Scan(
        &user.ID, &user.Auth0UserID, &user.Email, 
        &user.Name, &user.CreatedAt, &user.UpdatedAt,
    )
    
    if err != nil {
        // If user not found, create new user
        if err.Error() == "sql: no rows in result set" {
            return createNewUser(auth0UserID)
        }
        log.Printf("Failed to query user: %v", err)
        return utils.ErrorResponse(500, "Failed to fetch user"), nil
    }
    
    return utils.SuccessResponse(user), nil
}

func createNewUser(auth0UserID string) (events.APIGatewayProxyResponse, error) {
    database, err := db.GetDB()
    if err != nil {
        return utils.ErrorResponse(500, "Database connection error"), nil
    }
    
    var user models.User
    query := `INSERT INTO users (auth0_user_id, email, name) 
              VALUES ($1, $2, $3) 
              RETURNING id, auth0_user_id, email, name, created_at, updated_at`
    
    // Default values for new user
    email := fmt.Sprintf("%s@example.com", auth0UserID)
    name := "New User"
    
    err = database.QueryRow(query, auth0UserID, email, name).Scan(
        &user.ID, &user.Auth0UserID, &user.Email,
        &user.Name, &user.CreatedAt, &user.UpdatedAt,
    )
    
    if err != nil {
        log.Printf("Failed to create user: %v", err)
        return utils.ErrorResponse(500, "Failed to create user"), nil
    }
    
    return utils.SuccessResponse(user), nil
}

func handleUpdateUser(auth0UserID string, body string) (events.APIGatewayProxyResponse, error) {
    var updateData struct {
        Name  string `json:"name"`
        Email string `json:"email"`
    }
    
    if err := json.Unmarshal([]byte(body), &updateData); err != nil {
        return utils.ErrorResponse(400, "Invalid request body"), nil
    }
    
    database, err := db.GetDB()
    if err != nil {
        return utils.ErrorResponse(500, "Database connection error"), nil
    }
    
    var user models.User
    query := `UPDATE users SET name = $1, email = $2, updated_at = NOW() 
              WHERE auth0_user_id = $3
              RETURNING id, auth0_user_id, email, name, created_at, updated_at`
    
    err = database.QueryRow(query, updateData.Name, updateData.Email, auth0UserID).Scan(
        &user.ID, &user.Auth0UserID, &user.Email,
        &user.Name, &user.CreatedAt, &user.UpdatedAt,
    )
    
    if err != nil {
        log.Printf("Failed to update user: %v", err)
        return utils.ErrorResponse(500, "Failed to update user"), nil
    }
    
    return utils.SuccessResponse(user), nil
}

func main() {
    lambda.Start(handler)
}
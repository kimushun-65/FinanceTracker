package main

import (
    "context"
    "encoding/json"
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
    
    // First, get internal user ID from Auth0 user ID
    internalUserID, err := getUserID(userID)
    if err != nil {
        return utils.ErrorResponse(500, "Failed to get user"), nil
    }
    
    switch request.HTTPMethod {
    case "GET":
        return handleGetAccounts(internalUserID)
    case "POST":
        return handleCreateAccount(internalUserID, request.Body)
    case "PUT":
        accountID := request.PathParameters["id"]
        return handleUpdateAccount(internalUserID, accountID, request.Body)
    case "DELETE":
        accountID := request.PathParameters["id"]
        return handleDeleteAccount(internalUserID, accountID)
    default:
        return utils.ErrorResponse(405, "Method not allowed"), nil
    }
}

func getUserID(auth0UserID string) (string, error) {
    database, err := db.GetDB()
    if err != nil {
        return "", err
    }
    
    var userID string
    query := `SELECT id FROM users WHERE auth0_user_id = $1`
    err = database.QueryRow(query, auth0UserID).Scan(&userID)
    
    return userID, err
}

func handleGetAccounts(userID string) (events.APIGatewayProxyResponse, error) {
    database, err := db.GetDB()
    if err != nil {
        log.Printf("Failed to connect to database: %v", err)
        return utils.ErrorResponse(500, "Database connection error"), nil
    }
    
    query := `SELECT id, user_id, name, type, balance, currency, created_at, updated_at 
              FROM accounts WHERE user_id = $1 ORDER BY created_at DESC`
    
    rows, err := database.Query(query, userID)
    if err != nil {
        log.Printf("Failed to query accounts: %v", err)
        return utils.ErrorResponse(500, "Failed to fetch accounts"), nil
    }
    defer rows.Close()
    
    accounts := []models.Account{}
    for rows.Next() {
        var account models.Account
        err := rows.Scan(
            &account.ID, &account.UserID, &account.Name, &account.Type,
            &account.Balance, &account.Currency, &account.CreatedAt, &account.UpdatedAt,
        )
        if err != nil {
            log.Printf("Failed to scan account: %v", err)
            continue
        }
        accounts = append(accounts, account)
    }
    
    return utils.SuccessResponse(accounts), nil
}

func handleCreateAccount(userID string, body string) (events.APIGatewayProxyResponse, error) {
    var createData struct {
        Name     string  `json:"name"`
        Type     string  `json:"type"`
        Balance  float64 `json:"balance"`
        Currency string  `json:"currency"`
    }
    
    if err := json.Unmarshal([]byte(body), &createData); err != nil {
        return utils.ErrorResponse(400, "Invalid request body"), nil
    }
    
    // Set default currency if not provided
    if createData.Currency == "" {
        createData.Currency = "JPY"
    }
    
    database, err := db.GetDB()
    if err != nil {
        return utils.ErrorResponse(500, "Database connection error"), nil
    }
    
    var account models.Account
    query := `INSERT INTO accounts (user_id, name, type, balance, currency) 
              VALUES ($1, $2, $3, $4, $5) 
              RETURNING id, user_id, name, type, balance, currency, created_at, updated_at`
    
    err = database.QueryRow(
        query, userID, createData.Name, createData.Type, 
        createData.Balance, createData.Currency,
    ).Scan(
        &account.ID, &account.UserID, &account.Name, &account.Type,
        &account.Balance, &account.Currency, &account.CreatedAt, &account.UpdatedAt,
    )
    
    if err != nil {
        log.Printf("Failed to create account: %v", err)
        return utils.ErrorResponse(500, "Failed to create account"), nil
    }
    
    return utils.SuccessResponse(account), nil
}

func handleUpdateAccount(userID, accountID string, body string) (events.APIGatewayProxyResponse, error) {
    var updateData struct {
        Name    string  `json:"name"`
        Balance float64 `json:"balance"`
    }
    
    if err := json.Unmarshal([]byte(body), &updateData); err != nil {
        return utils.ErrorResponse(400, "Invalid request body"), nil
    }
    
    database, err := db.GetDB()
    if err != nil {
        return utils.ErrorResponse(500, "Database connection error"), nil
    }
    
    var account models.Account
    query := `UPDATE accounts SET name = $1, balance = $2, updated_at = NOW() 
              WHERE id = $3 AND user_id = $4
              RETURNING id, user_id, name, type, balance, currency, created_at, updated_at`
    
    err = database.QueryRow(
        query, updateData.Name, updateData.Balance, accountID, userID,
    ).Scan(
        &account.ID, &account.UserID, &account.Name, &account.Type,
        &account.Balance, &account.Currency, &account.CreatedAt, &account.UpdatedAt,
    )
    
    if err != nil {
        log.Printf("Failed to update account: %v", err)
        return utils.ErrorResponse(500, "Failed to update account"), nil
    }
    
    return utils.SuccessResponse(account), nil
}

func handleDeleteAccount(userID, accountID string) (events.APIGatewayProxyResponse, error) {
    database, err := db.GetDB()
    if err != nil {
        return utils.ErrorResponse(500, "Database connection error"), nil
    }
    
    query := `DELETE FROM accounts WHERE id = $1 AND user_id = $2`
    result, err := database.Exec(query, accountID, userID)
    
    if err != nil {
        log.Printf("Failed to delete account: %v", err)
        return utils.ErrorResponse(500, "Failed to delete account"), nil
    }
    
    rowsAffected, _ := result.RowsAffected()
    if rowsAffected == 0 {
        return utils.ErrorResponse(404, "Account not found"), nil
    }
    
    return utils.SuccessResponse(map[string]string{"message": "Account deleted successfully"}), nil
}

func main() {
    lambda.Start(handler)
}
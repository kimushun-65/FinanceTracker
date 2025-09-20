package db

import (
    "database/sql"
    "encoding/json"
    "fmt"
    "os"
    
    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/secretsmanager"
    _ "github.com/lib/pq"
)

type DBSecret struct {
    Username string `json:"username"`
    Password string `json:"password"`
}

var db *sql.DB

// GetDB returns a database connection
func GetDB() (*sql.DB, error) {
    if db != nil {
        return db, nil
    }
    
    // Get database credentials from Secrets Manager
    sess, err := session.NewSession(&aws.Config{
        Region: aws.String(os.Getenv("AWS_REGION")),
    })
    if err != nil {
        return nil, fmt.Errorf("failed to create AWS session: %v", err)
    }
    
    svc := secretsmanager.New(sess)
    secretArn := os.Getenv("DB_SECRET_ARN")
    
    result, err := svc.GetSecretValue(&secretsmanager.GetSecretValueInput{
        SecretId: aws.String(secretArn),
    })
    if err != nil {
        return nil, fmt.Errorf("failed to get secret value: %v", err)
    }
    
    var secret DBSecret
    if err := json.Unmarshal([]byte(*result.SecretString), &secret); err != nil {
        return nil, fmt.Errorf("failed to unmarshal secret: %v", err)
    }
    
    // Connect to database
    dbHost := os.Getenv("DB_ENDPOINT")
    connStr := fmt.Sprintf("host=%s port=5432 user=%s password=%s dbname=finsight sslmode=require",
        dbHost, secret.Username, secret.Password)
    
    db, err = sql.Open("postgres", connStr)
    if err != nil {
        return nil, fmt.Errorf("failed to open database: %v", err)
    }
    
    // Test connection
    if err := db.Ping(); err != nil {
        return nil, fmt.Errorf("failed to ping database: %v", err)
    }
    
    return db, nil
}

// CloseDB closes the database connection
func CloseDB() error {
    if db != nil {
        return db.Close()
    }
    return nil
}
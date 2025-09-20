package models

import (
    "time"
)

// User represents a user in the system
type User struct {
    ID           string    `json:"id"`
    Auth0UserID  string    `json:"auth0_user_id"`
    Email        string    `json:"email"`
    Name         string    `json:"name"`
    CreatedAt    time.Time `json:"created_at"`
    UpdatedAt    time.Time `json:"updated_at"`
}

// Account represents a financial account
type Account struct {
    ID        string    `json:"id"`
    UserID    string    `json:"user_id"`
    Name      string    `json:"name"`
    Type      string    `json:"type"`
    Balance   float64   `json:"balance"`
    Currency  string    `json:"currency"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}

// Category represents a transaction category
type Category struct {
    ID        string    `json:"id"`
    UserID    string    `json:"user_id"`
    Name      string    `json:"name"`
    Type      string    `json:"type"`
    Icon      *string   `json:"icon,omitempty"`
    Color     *string   `json:"color,omitempty"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}

// Transaction represents a financial transaction
type Transaction struct {
    ID          string    `json:"id"`
    UserID      string    `json:"user_id"`
    AccountID   string    `json:"account_id"`
    CategoryID  *string   `json:"category_id,omitempty"`
    Type        string    `json:"type"`
    Amount      float64   `json:"amount"`
    Description *string   `json:"description,omitempty"`
    Date        string    `json:"date"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}

// Budget represents a budget
type Budget struct {
    ID         string    `json:"id"`
    UserID     string    `json:"user_id"`
    CategoryID *string   `json:"category_id,omitempty"`
    Name       string    `json:"name"`
    Amount     float64   `json:"amount"`
    Period     string    `json:"period"`
    StartDate  string    `json:"start_date"`
    EndDate    *string   `json:"end_date,omitempty"`
    CreatedAt  time.Time `json:"created_at"`
    UpdatedAt  time.Time `json:"updated_at"`
}
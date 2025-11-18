package models

import (
	"time"
)

type Expense struct {
	ID          string     `json:"id"`
	UserID      string     `json:"user_id"`
	ClientID    *string    `json:"client_id,omitempty"`
	Description string     `json:"description"`
	Amount      float64    `json:"amount"`
	Currency    string     `json:"currency"`
	Category    *string    `json:"category,omitempty"`
	ExpenseDate time.Time  `json:"expense_date"`
	ReceiptURL  *string    `json:"receipt_url,omitempty"`
	Notes       *string    `json:"notes,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	Client      *Client    `json:"client,omitempty"`
}


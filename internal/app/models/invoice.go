package models

import "time"

type InvoiceStatus string

const (
	InvoiceStatusDraft    InvoiceStatus = "draft"
	InvoiceStatusPending InvoiceStatus = "pending"
	InvoiceStatusPaid    InvoiceStatus = "paid"
	InvoiceStatusOverdue InvoiceStatus = "overdue"
	InvoiceStatusCancelled InvoiceStatus = "cancelled"
)

type Invoice struct {
	ID           string         `json:"id"`
	UserID       string         `json:"user_id"`
	ClientID     string         `json:"client_id"`
	InvoiceNumber string        `json:"invoice_number"`
	Status       InvoiceStatus  `json:"status"`
	IssueDate    time.Time      `json:"issue_date"`
	DueDate      *time.Time     `json:"due_date,omitempty"`
	Currency     string         `json:"currency"`
	Subtotal     float64        `json:"subtotal"`
	TaxRate      float64        `json:"tax_rate"`
	TaxAmount    float64        `json:"tax_amount"`
	Total        float64        `json:"total"`
	Notes        *string        `json:"notes,omitempty"`
	PaymentLink  *string        `json:"payment_link,omitempty"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	Items        []InvoiceItem  `json:"items,omitempty"`
	Client       *Client        `json:"client,omitempty"`
	Payments     []Payment      `json:"payments,omitempty"`
}

type InvoiceItem struct {
	ID          string    `json:"id"`
	InvoiceID   string    `json:"invoice_id"`
	Description string    `json:"description"`
	Quantity    float64   `json:"quantity"`
	UnitPrice   float64   `json:"unit_price"`
	Amount      float64   `json:"amount"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Payment struct {
	ID            string         `json:"id"`
	InvoiceID     string         `json:"invoice_id"`
	Amount        float64        `json:"amount"`
	Currency      string         `json:"currency"`
	PaymentMethod *string        `json:"payment_method,omitempty"`
	PaymentDate   time.Time      `json:"payment_date"`
	TransactionID *string        `json:"transaction_id,omitempty"`
	Notes         *string        `json:"notes,omitempty"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
}


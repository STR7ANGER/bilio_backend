package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/nava1525/bilio-backend/internal/app/models"
	"github.com/nava1525/bilio-backend/internal/app/repositories"
)

type InvoiceService struct {
	invoices repositories.InvoiceRepository
	clients  repositories.ClientRepository
}

type CreateInvoiceInput struct {
	ClientID      string                 `json:"client_id"`
	InvoiceNumber string                 `json:"invoice_number"`
	Status        models.InvoiceStatus   `json:"status"`
	IssueDate     time.Time              `json:"issue_date"`
	DueDate       *time.Time             `json:"due_date,omitempty"`
	Currency      string                 `json:"currency"`
	TaxRate       float64                `json:"tax_rate"`
	Notes         *string                `json:"notes,omitempty"`
	Items         []CreateInvoiceItemInput `json:"items"`
}

type CreateInvoiceItemInput struct {
	Description string  `json:"description"`
	Quantity    float64 `json:"quantity"`
	UnitPrice   float64 `json:"unit_price"`
}

type UpdateInvoiceInput struct {
	Status    *models.InvoiceStatus `json:"status,omitempty"`
	IssueDate *time.Time            `json:"issue_date,omitempty"`
	DueDate   *time.Time            `json:"due_date,omitempty"`
	Currency  string                `json:"currency"`
	TaxRate   float64               `json:"tax_rate"`
	Notes     *string               `json:"notes,omitempty"`
	Items     []CreateInvoiceItemInput `json:"items,omitempty"`
}

type InvoiceFilters struct {
	Status   *models.InvoiceStatus
	ClientID *string
	FromDate *time.Time
	ToDate   *time.Time
}

func NewInvoiceService(invoiceRepo repositories.InvoiceRepository, clientRepo repositories.ClientRepository) *InvoiceService {
	return &InvoiceService{
		invoices: invoiceRepo,
		clients:  clientRepo,
	}
}

func (s *InvoiceService) List(ctx context.Context, userID string, filters InvoiceFilters) ([]models.Invoice, error) {
	repoFilters := repositories.InvoiceFilters{
		Status:   filters.Status,
		ClientID: filters.ClientID,
		FromDate: filters.FromDate,
		ToDate:   filters.ToDate,
	}
	return s.invoices.List(ctx, userID, repoFilters)
}

func (s *InvoiceService) GetByID(ctx context.Context, id string, userID string) (*models.Invoice, error) {
	invoice, err := s.invoices.GetByID(ctx, id, userID)
	if err != nil {
		return nil, err
	}
	if invoice == nil {
		return nil, errors.New("invoice not found")
	}

	// Load items
	items, err := s.invoices.GetItems(ctx, id)
	if err != nil {
		return nil, err
	}
	invoice.Items = items

	// Load payments
	payments, err := s.invoices.GetPayments(ctx, id)
	if err != nil {
		return nil, err
	}
	invoice.Payments = payments

	return invoice, nil
}

func (s *InvoiceService) Create(ctx context.Context, userID string, input CreateInvoiceInput) (*models.Invoice, error) {
	// Verify client exists and belongs to user
	client, err := s.clients.GetByID(ctx, input.ClientID, userID)
	if err != nil {
		return nil, err
	}
	if client == nil {
		return nil, errors.New("client not found")
	}

	if input.InvoiceNumber == "" {
		return nil, errors.New("invoice_number is required")
	}
	if len(input.Items) == 0 {
		return nil, errors.New("at least one item is required")
	}
	if input.Currency == "" {
		input.Currency = "USD"
	}
	if input.Status == "" {
		input.Status = models.InvoiceStatusDraft
	}

	// Calculate totals
	subtotal := 0.0
	for _, item := range input.Items {
		amount := item.Quantity * item.UnitPrice
		subtotal += amount
	}

	taxAmount := subtotal * (input.TaxRate / 100)
	total := subtotal + taxAmount

	invoice := &models.Invoice{
		UserID:        userID,
		ClientID:      input.ClientID,
		InvoiceNumber: input.InvoiceNumber,
		Status:        input.Status,
		IssueDate:     input.IssueDate,
		DueDate:       input.DueDate,
		Currency:      input.Currency,
		Subtotal:      subtotal,
		TaxRate:       input.TaxRate,
		TaxAmount:     taxAmount,
		Total:         total,
		Notes:         input.Notes,
	}

	created, err := s.invoices.Create(ctx, invoice)
	if err != nil {
		return nil, fmt.Errorf("failed to create invoice: %w", err)
	}

	// Create items
	for _, itemInput := range input.Items {
		amount := itemInput.Quantity * itemInput.UnitPrice
		item := &models.InvoiceItem{
			InvoiceID:   created.ID,
			Description: itemInput.Description,
			Quantity:    itemInput.Quantity,
			UnitPrice:   itemInput.UnitPrice,
			Amount:      amount,
		}
		if err := s.invoices.CreateItem(ctx, item); err != nil {
			return nil, fmt.Errorf("failed to create invoice item: %w", err)
		}
		created.Items = append(created.Items, *item)
	}

	return created, nil
}

func (s *InvoiceService) Update(ctx context.Context, id string, userID string, input UpdateInvoiceInput) (*models.Invoice, error) {
	invoice, err := s.invoices.GetByID(ctx, id, userID)
	if err != nil {
		return nil, err
	}
	if invoice == nil {
		return nil, errors.New("invoice not found")
	}

	// Only allow updates to draft or pending invoices
	if invoice.Status != models.InvoiceStatusDraft && invoice.Status != models.InvoiceStatusPending {
		return nil, errors.New("can only update draft or pending invoices")
	}

	if input.Status != nil {
		invoice.Status = *input.Status
	}
	if input.IssueDate != nil {
		invoice.IssueDate = *input.IssueDate
	}
	if input.DueDate != nil {
		invoice.DueDate = input.DueDate
	}
	if input.Currency != "" {
		invoice.Currency = input.Currency
	}
	invoice.TaxRate = input.TaxRate
	if input.Notes != nil {
		invoice.Notes = input.Notes
	}

	// Update items if provided
	if len(input.Items) > 0 {
		// Delete existing items
		existingItems, _ := s.invoices.GetItems(ctx, id)
		for _, item := range existingItems {
			_ = s.invoices.DeleteItem(ctx, item.ID)
		}

		// Recalculate totals
		subtotal := 0.0
		for _, itemInput := range input.Items {
			amount := itemInput.Quantity * itemInput.UnitPrice
			subtotal += amount

			item := &models.InvoiceItem{
				InvoiceID:   id,
				Description: itemInput.Description,
				Quantity:    itemInput.Quantity,
				UnitPrice:   itemInput.UnitPrice,
				Amount:      amount,
			}
			if err := s.invoices.CreateItem(ctx, item); err != nil {
				return nil, fmt.Errorf("failed to create invoice item: %w", err)
			}
		}

		invoice.Subtotal = subtotal
		invoice.TaxAmount = subtotal * (input.TaxRate / 100)
		invoice.Total = invoice.Subtotal + invoice.TaxAmount
	}

	return s.invoices.Update(ctx, invoice)
}

func (s *InvoiceService) MarkPaid(ctx context.Context, id string, userID string, paymentInput CreatePaymentInput) (*models.Invoice, error) {
	invoice, err := s.invoices.GetByID(ctx, id, userID)
	if err != nil {
		return nil, err
	}
	if invoice == nil {
		return nil, errors.New("invoice not found")
	}

	payment := &models.Payment{
		InvoiceID:     id,
		Amount:        paymentInput.Amount,
		Currency:      paymentInput.Currency,
		PaymentMethod: paymentInput.PaymentMethod,
		PaymentDate:   paymentInput.PaymentDate,
		TransactionID: paymentInput.TransactionID,
		Notes:         paymentInput.Notes,
	}

	if err := s.invoices.CreatePayment(ctx, payment); err != nil {
		return nil, fmt.Errorf("failed to create payment: %w", err)
	}

	// Update invoice status
	invoice.Status = models.InvoiceStatusPaid
	return s.invoices.Update(ctx, invoice)
}

type CreatePaymentInput struct {
	Amount        float64  `json:"amount"`
	Currency      string   `json:"currency"`
	PaymentMethod *string  `json:"payment_method,omitempty"`
	PaymentDate   time.Time `json:"payment_date"`
	TransactionID *string  `json:"transaction_id,omitempty"`
	Notes         *string  `json:"notes,omitempty"`
}


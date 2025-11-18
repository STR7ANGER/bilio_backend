package services

import (
	"context"
	"errors"
	"time"

	"github.com/nava1525/bilio-backend/internal/app/models"
	"github.com/nava1525/bilio-backend/internal/app/repositories"
)

type ExpenseService struct {
	expenses repositories.ExpenseRepository
	clients  repositories.ClientRepository
}

type CreateExpenseInput struct {
	ClientID    *string   `json:"client_id,omitempty"`
	Description string    `json:"description"`
	Amount      float64   `json:"amount"`
	Currency    string    `json:"currency"`
	Category    *string   `json:"category,omitempty"`
	ExpenseDate time.Time `json:"expense_date"`
	ReceiptURL  *string   `json:"receipt_url,omitempty"`
	Notes       *string   `json:"notes,omitempty"`
}

type UpdateExpenseInput struct {
	ClientID    *string   `json:"client_id,omitempty"`
	Description string    `json:"description"`
	Amount      float64   `json:"amount"`
	Currency    string    `json:"currency"`
	Category    *string   `json:"category,omitempty"`
	ExpenseDate time.Time `json:"expense_date"`
	ReceiptURL  *string   `json:"receipt_url,omitempty"`
	Notes       *string   `json:"notes,omitempty"`
}

type ExpenseFilters struct {
	ClientID *string
	Category *string
	FromDate *time.Time
	ToDate   *time.Time
}

func NewExpenseService(expenseRepo repositories.ExpenseRepository, clientRepo repositories.ClientRepository) *ExpenseService {
	return &ExpenseService{
		expenses: expenseRepo,
		clients:  clientRepo,
	}
}

func (s *ExpenseService) List(ctx context.Context, userID string, filters ExpenseFilters) ([]models.Expense, error) {
	repoFilters := repositories.ExpenseFilters{
		ClientID: filters.ClientID,
		Category: filters.Category,
		FromDate: filters.FromDate,
		ToDate:   filters.ToDate,
	}
	return s.expenses.List(ctx, userID, repoFilters)
}

func (s *ExpenseService) GetByID(ctx context.Context, id string, userID string) (*models.Expense, error) {
	expense, err := s.expenses.GetByID(ctx, id, userID)
	if err != nil {
		return nil, err
	}
	if expense == nil {
		return nil, errors.New("expense not found")
	}
	return expense, nil
}

func (s *ExpenseService) Create(ctx context.Context, userID string, input CreateExpenseInput) (*models.Expense, error) {
	if input.Description == "" {
		return nil, errors.New("description is required")
	}
	if input.Amount <= 0 {
		return nil, errors.New("amount must be greater than 0")
	}
	if input.Currency == "" {
		input.Currency = "USD"
	}

	// Verify client exists if provided
	if input.ClientID != nil {
		client, err := s.clients.GetByID(ctx, *input.ClientID, userID)
		if err != nil {
			return nil, err
		}
		if client == nil {
			return nil, errors.New("client not found")
		}
	}

	expense := &models.Expense{
		UserID:      userID,
		ClientID:    input.ClientID,
		Description: input.Description,
		Amount:      input.Amount,
		Currency:    input.Currency,
		Category:    input.Category,
		ExpenseDate: input.ExpenseDate,
		ReceiptURL:  input.ReceiptURL,
		Notes:       input.Notes,
	}

	return s.expenses.Create(ctx, expense)
}

func (s *ExpenseService) Update(ctx context.Context, id string, userID string, input UpdateExpenseInput) (*models.Expense, error) {
	expense, err := s.expenses.GetByID(ctx, id, userID)
	if err != nil {
		return nil, err
	}
	if expense == nil {
		return nil, errors.New("expense not found")
	}

	// Verify client exists if provided
	if input.ClientID != nil {
		client, err := s.clients.GetByID(ctx, *input.ClientID, userID)
		if err != nil {
			return nil, err
		}
		if client == nil {
			return nil, errors.New("client not found")
		}
	}

	expense.Description = input.Description
	expense.Amount = input.Amount
	expense.Currency = input.Currency
	expense.Category = input.Category
	expense.ExpenseDate = input.ExpenseDate
	expense.ReceiptURL = input.ReceiptURL
	expense.Notes = input.Notes
	expense.ClientID = input.ClientID

	return s.expenses.Update(ctx, expense)
}


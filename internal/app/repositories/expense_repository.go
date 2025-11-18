package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/nava1525/bilio-backend/internal/app/models"
)

type ExpenseRepository interface {
	List(ctx context.Context, userID string, filters ExpenseFilters) ([]models.Expense, error)
	GetByID(ctx context.Context, id string, userID string) (*models.Expense, error)
	Create(ctx context.Context, expense *models.Expense) (*models.Expense, error)
	Update(ctx context.Context, expense *models.Expense) (*models.Expense, error)
}

type ExpenseFilters struct {
	ClientID *string
	Category *string
	FromDate *time.Time
	ToDate   *time.Time
}

type postgresExpenseRepository struct {
	db *sql.DB
}

func NewExpenseRepository(db *sql.DB) ExpenseRepository {
	return &postgresExpenseRepository{db: db}
}

func (r *postgresExpenseRepository) List(ctx context.Context, userID string, filters ExpenseFilters) ([]models.Expense, error) {
	query := `SELECT id, user_id, client_id, description, amount, currency, category, expense_date, receipt_url, notes, created_at, updated_at
			  FROM expenses WHERE user_id = $1`
	args := []interface{}{userID}
	argPos := 2

	if filters.ClientID != nil {
		query += ` AND client_id = $` + fmt.Sprintf("%d", argPos)
		args = append(args, *filters.ClientID)
		argPos++
	}
	if filters.Category != nil {
		query += ` AND category = $` + fmt.Sprintf("%d", argPos)
		args = append(args, *filters.Category)
		argPos++
	}
	if filters.FromDate != nil {
		query += ` AND expense_date >= $` + fmt.Sprintf("%d", argPos)
		args = append(args, *filters.FromDate)
		argPos++
	}
	if filters.ToDate != nil {
		query += ` AND expense_date <= $` + fmt.Sprintf("%d", argPos)
		args = append(args, *filters.ToDate)
		argPos++
	}

	query += ` ORDER BY expense_date DESC, created_at DESC`

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var expenses []models.Expense
	for rows.Next() {
		var e models.Expense
		var clientID, category, receiptURL, notes sql.NullString

		if err := rows.Scan(&e.ID, &e.UserID, &clientID, &e.Description, &e.Amount, &e.Currency,
			&category, &e.ExpenseDate, &receiptURL, &notes, &e.CreatedAt, &e.UpdatedAt); err != nil {
			return nil, err
		}

		if clientID.Valid {
			e.ClientID = &clientID.String
		}
		if category.Valid {
			e.Category = &category.String
		}
		if receiptURL.Valid {
			e.ReceiptURL = &receiptURL.String
		}
		if notes.Valid {
			e.Notes = &notes.String
		}

		expenses = append(expenses, e)
	}

	return expenses, rows.Err()
}

func (r *postgresExpenseRepository) GetByID(ctx context.Context, id string, userID string) (*models.Expense, error) {
	var e models.Expense
	var clientID, category, receiptURL, notes sql.NullString

	err := r.db.QueryRowContext(ctx,
		`SELECT id, user_id, client_id, description, amount, currency, category, expense_date, receipt_url, notes, created_at, updated_at
		 FROM expenses WHERE id = $1 AND user_id = $2`,
		id, userID).Scan(&e.ID, &e.UserID, &clientID, &e.Description, &e.Amount, &e.Currency,
		&category, &e.ExpenseDate, &receiptURL, &notes, &e.CreatedAt, &e.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	if clientID.Valid {
		e.ClientID = &clientID.String
	}
	if category.Valid {
		e.Category = &category.String
	}
	if receiptURL.Valid {
		e.ReceiptURL = &receiptURL.String
	}
	if notes.Valid {
		e.Notes = &notes.String
	}

	return &e, nil
}

func (r *postgresExpenseRepository) Create(ctx context.Context, expense *models.Expense) (*models.Expense, error) {
	id := uuid.NewString()
	now := time.Now().UTC()

	_, err := r.db.ExecContext(ctx,
		`INSERT INTO expenses (id, user_id, client_id, description, amount, currency, category, expense_date, receipt_url, notes, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $11)`,
		id, expense.UserID, expense.ClientID, expense.Description, expense.Amount, expense.Currency,
		expense.Category, expense.ExpenseDate, expense.ReceiptURL, expense.Notes, now)
	if err != nil {
		return nil, err
	}

	expense.ID = id
	expense.CreatedAt = now
	expense.UpdatedAt = now
	return expense, nil
}

func (r *postgresExpenseRepository) Update(ctx context.Context, expense *models.Expense) (*models.Expense, error) {
	now := time.Now().UTC()

	_, err := r.db.ExecContext(ctx,
		`UPDATE expenses SET description = $1, amount = $2, currency = $3, category = $4, expense_date = $5,
		 receipt_url = $6, notes = $7, client_id = $8, updated_at = $9
		 WHERE id = $10 AND user_id = $11`,
		expense.Description, expense.Amount, expense.Currency, expense.Category, expense.ExpenseDate,
		expense.ReceiptURL, expense.Notes, expense.ClientID, now, expense.ID, expense.UserID)
	if err != nil {
		return nil, err
	}

	expense.UpdatedAt = now
	return expense, nil
}


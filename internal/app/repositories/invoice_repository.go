package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/nava1525/bilio-backend/internal/app/models"
)

type InvoiceRepository interface {
	List(ctx context.Context, userID string, filters InvoiceFilters) ([]models.Invoice, error)
	GetByID(ctx context.Context, id string, userID string) (*models.Invoice, error)
	Create(ctx context.Context, invoice *models.Invoice) (*models.Invoice, error)
	Update(ctx context.Context, invoice *models.Invoice) (*models.Invoice, error)
	GetItems(ctx context.Context, invoiceID string) ([]models.InvoiceItem, error)
	CreateItem(ctx context.Context, item *models.InvoiceItem) error
	UpdateItem(ctx context.Context, item *models.InvoiceItem) error
	DeleteItem(ctx context.Context, itemID string) error
	GetPayments(ctx context.Context, invoiceID string) ([]models.Payment, error)
	CreatePayment(ctx context.Context, payment *models.Payment) error
}

type InvoiceFilters struct {
	Status   *models.InvoiceStatus
	ClientID *string
	FromDate *time.Time
	ToDate   *time.Time
}

type postgresInvoiceRepository struct {
	db *sql.DB
}

func NewInvoiceRepository(db *sql.DB) InvoiceRepository {
	return &postgresInvoiceRepository{db: db}
}

func (r *postgresInvoiceRepository) List(ctx context.Context, userID string, filters InvoiceFilters) ([]models.Invoice, error) {
	query := `SELECT id, user_id, client_id, invoice_number, status, issue_date, due_date, currency, 
			  subtotal, tax_rate, tax_amount, total, notes, payment_link, created_at, updated_at
			  FROM invoices WHERE user_id = $1`
	args := []interface{}{userID}
	argPos := 2

	if filters.Status != nil {
		query += ` AND status = $` + fmt.Sprintf("%d", argPos)
		args = append(args, *filters.Status)
		argPos++
	}
	if filters.ClientID != nil {
		query += ` AND client_id = $` + fmt.Sprintf("%d", argPos)
		args = append(args, *filters.ClientID)
		argPos++
	}
	if filters.FromDate != nil {
		query += ` AND issue_date >= $` + fmt.Sprintf("%d", argPos)
		args = append(args, *filters.FromDate)
		argPos++
	}
	if filters.ToDate != nil {
		query += ` AND issue_date <= $` + fmt.Sprintf("%d", argPos)
		args = append(args, *filters.ToDate)
		argPos++
	}

	query += ` ORDER BY created_at DESC`

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var invoices []models.Invoice
	for rows.Next() {
		var inv models.Invoice
		var dueDate sql.NullTime
		var notes, paymentLink sql.NullString

		err := rows.Scan(&inv.ID, &inv.UserID, &inv.ClientID, &inv.InvoiceNumber, &inv.Status,
			&inv.IssueDate, &dueDate, &inv.Currency, &inv.Subtotal, &inv.TaxRate, &inv.TaxAmount,
			&inv.Total, &notes, &paymentLink, &inv.CreatedAt, &inv.UpdatedAt)
		if err != nil {
			return nil, err
		}

		if dueDate.Valid {
			inv.DueDate = &dueDate.Time
		}
		if notes.Valid {
			inv.Notes = &notes.String
		}
		if paymentLink.Valid {
			inv.PaymentLink = &paymentLink.String
		}

		invoices = append(invoices, inv)
	}

	return invoices, rows.Err()
}

func (r *postgresInvoiceRepository) GetByID(ctx context.Context, id string, userID string) (*models.Invoice, error) {
	var inv models.Invoice
	var dueDate sql.NullTime
	var notes, paymentLink sql.NullString

	err := r.db.QueryRowContext(ctx,
		`SELECT id, user_id, client_id, invoice_number, status, issue_date, due_date, currency,
		 subtotal, tax_rate, tax_amount, total, notes, payment_link, created_at, updated_at
		 FROM invoices WHERE id = $1 AND user_id = $2`,
		id, userID).Scan(&inv.ID, &inv.UserID, &inv.ClientID, &inv.InvoiceNumber, &inv.Status,
		&inv.IssueDate, &dueDate, &inv.Currency, &inv.Subtotal, &inv.TaxRate, &inv.TaxAmount,
		&inv.Total, &notes, &paymentLink, &inv.CreatedAt, &inv.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	if dueDate.Valid {
		inv.DueDate = &dueDate.Time
	}
	if notes.Valid {
		inv.Notes = &notes.String
	}
	if paymentLink.Valid {
		inv.PaymentLink = &paymentLink.String
	}

	return &inv, nil
}

func (r *postgresInvoiceRepository) Create(ctx context.Context, invoice *models.Invoice) (*models.Invoice, error) {
	id := uuid.NewString()
	now := time.Now().UTC()

	_, err := r.db.ExecContext(ctx,
		`INSERT INTO invoices (id, user_id, client_id, invoice_number, status, issue_date, due_date, currency,
		 subtotal, tax_rate, tax_amount, total, notes, payment_link, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $15)`,
		id, invoice.UserID, invoice.ClientID, invoice.InvoiceNumber, invoice.Status, invoice.IssueDate,
		invoice.DueDate, invoice.Currency, invoice.Subtotal, invoice.TaxRate, invoice.TaxAmount,
		invoice.Total, invoice.Notes, invoice.PaymentLink, now)
	if err != nil {
		return nil, err
	}

	invoice.ID = id
	invoice.CreatedAt = now
	invoice.UpdatedAt = now
	return invoice, nil
}

func (r *postgresInvoiceRepository) Update(ctx context.Context, invoice *models.Invoice) (*models.Invoice, error) {
	now := time.Now().UTC()

	_, err := r.db.ExecContext(ctx,
		`UPDATE invoices SET status = $1, issue_date = $2, due_date = $3, currency = $4,
		 subtotal = $5, tax_rate = $6, tax_amount = $7, total = $8, notes = $9, payment_link = $10, updated_at = $11
		 WHERE id = $12 AND user_id = $13`,
		invoice.Status, invoice.IssueDate, invoice.DueDate, invoice.Currency, invoice.Subtotal,
		invoice.TaxRate, invoice.TaxAmount, invoice.Total, invoice.Notes, invoice.PaymentLink, now, invoice.ID, invoice.UserID)
	if err != nil {
		return nil, err
	}

	invoice.UpdatedAt = now
	return invoice, nil
}

func (r *postgresInvoiceRepository) GetItems(ctx context.Context, invoiceID string) ([]models.InvoiceItem, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, invoice_id, description, quantity, unit_price, amount, created_at, updated_at
		 FROM invoice_items WHERE invoice_id = $1 ORDER BY created_at`,
		invoiceID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []models.InvoiceItem
	for rows.Next() {
		var item models.InvoiceItem
		if err := rows.Scan(&item.ID, &item.InvoiceID, &item.Description, &item.Quantity,
			&item.UnitPrice, &item.Amount, &item.CreatedAt, &item.UpdatedAt); err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	return items, rows.Err()
}

func (r *postgresInvoiceRepository) CreateItem(ctx context.Context, item *models.InvoiceItem) error {
	id := uuid.NewString()
	now := time.Now().UTC()

	_, err := r.db.ExecContext(ctx,
		`INSERT INTO invoice_items (id, invoice_id, description, quantity, unit_price, amount, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $7)`,
		id, item.InvoiceID, item.Description, item.Quantity, item.UnitPrice, item.Amount, now)
	if err != nil {
		return err
	}

	item.ID = id
	item.CreatedAt = now
	item.UpdatedAt = now
	return nil
}

func (r *postgresInvoiceRepository) UpdateItem(ctx context.Context, item *models.InvoiceItem) error {
	now := time.Now().UTC()

	_, err := r.db.ExecContext(ctx,
		`UPDATE invoice_items SET description = $1, quantity = $2, unit_price = $3, amount = $4, updated_at = $5
		 WHERE id = $6`,
		item.Description, item.Quantity, item.UnitPrice, item.Amount, now, item.ID)
	return err
}

func (r *postgresInvoiceRepository) DeleteItem(ctx context.Context, itemID string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM invoice_items WHERE id = $1`, itemID)
	return err
}

func (r *postgresInvoiceRepository) GetPayments(ctx context.Context, invoiceID string) ([]models.Payment, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, invoice_id, amount, currency, payment_method, payment_date, transaction_id, notes, created_at, updated_at
		 FROM payments WHERE invoice_id = $1 ORDER BY payment_date DESC`,
		invoiceID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var payments []models.Payment
	for rows.Next() {
		var p models.Payment
		var paymentMethod, transactionID, notes sql.NullString

		if err := rows.Scan(&p.ID, &p.InvoiceID, &p.Amount, &p.Currency, &paymentMethod,
			&p.PaymentDate, &transactionID, &notes, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, err
		}

		if paymentMethod.Valid {
			p.PaymentMethod = &paymentMethod.String
		}
		if transactionID.Valid {
			p.TransactionID = &transactionID.String
		}
		if notes.Valid {
			p.Notes = &notes.String
		}

		payments = append(payments, p)
	}

	return payments, rows.Err()
}

func (r *postgresInvoiceRepository) CreatePayment(ctx context.Context, payment *models.Payment) error {
	id := uuid.NewString()
	now := time.Now().UTC()

	_, err := r.db.ExecContext(ctx,
		`INSERT INTO payments (id, invoice_id, amount, currency, payment_method, payment_date, transaction_id, notes, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $9)`,
		id, payment.InvoiceID, payment.Amount, payment.Currency, payment.PaymentMethod,
		payment.PaymentDate, payment.TransactionID, payment.Notes, now)
	if err != nil {
		return err
	}

	payment.ID = id
	payment.CreatedAt = now
	payment.UpdatedAt = now
	return nil
}


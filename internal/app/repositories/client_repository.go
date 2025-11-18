package repositories

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"

	"github.com/nava1525/bilio-backend/internal/app/models"
)

type ClientRepository interface {
	List(ctx context.Context, userID string) ([]models.Client, error)
	GetByID(ctx context.Context, id string, userID string) (*models.Client, error)
	Create(ctx context.Context, client *models.Client) (*models.Client, error)
	Update(ctx context.Context, client *models.Client) (*models.Client, error)
	Delete(ctx context.Context, id string, userID string) error
}

type postgresClientRepository struct {
	db *sql.DB
}

func NewClientRepository(db *sql.DB) ClientRepository {
	return &postgresClientRepository{db: db}
}

func (r *postgresClientRepository) List(ctx context.Context, userID string) ([]models.Client, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, user_id, name, email, company, phone, address, tax_id, currency, created_at, updated_at
		 FROM clients WHERE user_id = $1 ORDER BY created_at DESC`,
		userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var clients []models.Client
	for rows.Next() {
		var c models.Client
		var email, company, phone, address, taxID sql.NullString
		if err := rows.Scan(&c.ID, &c.UserID, &c.Name, &email, &company, &phone, &address, &taxID, &c.Currency, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, err
		}
		if email.Valid {
			c.Email = &email.String
		}
		if company.Valid {
			c.Company = &company.String
		}
		if phone.Valid {
			c.Phone = &phone.String
		}
		if address.Valid {
			c.Address = &address.String
		}
		if taxID.Valid {
			c.TaxID = &taxID.String
		}
		clients = append(clients, c)
	}

	return clients, rows.Err()
}

func (r *postgresClientRepository) GetByID(ctx context.Context, id string, userID string) (*models.Client, error) {
	var c models.Client
	var email, company, phone, address, taxID sql.NullString

	err := r.db.QueryRowContext(ctx,
		`SELECT id, user_id, name, email, company, phone, address, tax_id, currency, created_at, updated_at
		 FROM clients WHERE id = $1 AND user_id = $2`,
		id, userID).Scan(&c.ID, &c.UserID, &c.Name, &email, &company, &phone, &address, &taxID, &c.Currency, &c.CreatedAt, &c.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	if email.Valid {
		c.Email = &email.String
	}
	if company.Valid {
		c.Company = &company.String
	}
	if phone.Valid {
		c.Phone = &phone.String
	}
	if address.Valid {
		c.Address = &address.String
	}
	if taxID.Valid {
		c.TaxID = &taxID.String
	}

	return &c, nil
}

func (r *postgresClientRepository) Create(ctx context.Context, client *models.Client) (*models.Client, error) {
	id := uuid.NewString()
	now := time.Now().UTC()

	_, err := r.db.ExecContext(ctx,
		`INSERT INTO clients (id, user_id, name, email, company, phone, address, tax_id, currency, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $10)`,
		id, client.UserID, client.Name, client.Email, client.Company, client.Phone, client.Address, client.TaxID, client.Currency, now)
	if err != nil {
		return nil, err
	}

	client.ID = id
	client.CreatedAt = now
	client.UpdatedAt = now
	return client, nil
}

func (r *postgresClientRepository) Update(ctx context.Context, client *models.Client) (*models.Client, error) {
	now := time.Now().UTC()

	_, err := r.db.ExecContext(ctx,
		`UPDATE clients SET name = $1, email = $2, company = $3, phone = $4, address = $5, tax_id = $6, currency = $7, updated_at = $8
		 WHERE id = $9 AND user_id = $10`,
		client.Name, client.Email, client.Company, client.Phone, client.Address, client.TaxID, client.Currency, now, client.ID, client.UserID)
	if err != nil {
		return nil, err
	}

	client.UpdatedAt = now
	return client, nil
}

func (r *postgresClientRepository) Delete(ctx context.Context, id string, userID string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM clients WHERE id = $1 AND user_id = $2`, id, userID)
	return err
}


package repositories

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/nava1525/bilio-backend/internal/app/models"
)

type PromocodeRepository interface {
	Create(ctx context.Context, promocode *models.Promocode) (*models.Promocode, error)
	FindAndDelete(ctx context.Context, code string) (*models.Promocode, error)
}

type postgresPromocodeRepository struct {
	db *sql.DB
}

var ErrPromocodeNotFound = errors.New("promocode not found")
var ErrPromocodeAlreadyUsed = errors.New("promocode already used")

func NewPromocodeRepository(db *sql.DB) PromocodeRepository {
	return &postgresPromocodeRepository{db: db}
}

func (r *postgresPromocodeRepository) Create(ctx context.Context, promocode *models.Promocode) (*models.Promocode, error) {
	row := r.db.QueryRowContext(ctx,
		`INSERT INTO promocodes (code, created_at)
         VALUES ($1, $2)
         RETURNING code, created_at, used_at`,
		promocode.Code, promocode.CreatedAt,
	)

	var created models.Promocode
	if err := row.Scan(&created.Code, &created.CreatedAt, &created.UsedAt); err != nil {
		return nil, err
	}

	return &created, nil
}

func (r *postgresPromocodeRepository) FindAndDelete(ctx context.Context, code string) (*models.Promocode, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// Find the promocode
	row := tx.QueryRowContext(ctx,
		`SELECT code, created_at, used_at
         FROM promocodes
         WHERE code = $1`,
		code,
	)

	var promocode models.Promocode
	if err := row.Scan(&promocode.Code, &promocode.CreatedAt, &promocode.UsedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrPromocodeNotFound
		}
		return nil, err
	}

	// Check if already used
	if promocode.UsedAt != nil {
		return nil, ErrPromocodeAlreadyUsed
	}

	// Mark as used (soft delete by setting used_at) and get the timestamp back
	now := time.Now()
	row = tx.QueryRowContext(ctx,
		`UPDATE promocodes
         SET used_at = $1
         WHERE code = $2
         RETURNING code, created_at, used_at`,
		now, code,
	)

	var updated models.Promocode
	if err := row.Scan(&updated.Code, &updated.CreatedAt, &updated.UsedAt); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &updated, nil
}


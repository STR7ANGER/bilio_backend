package repositories

import (
	"context"
	"database/sql"
	"errors"

	"github.com/nava1525/bilio-backend/internal/app/models"
)

type WaitlistRepository interface {
	Create(ctx context.Context, entry *models.WaitlistEntry) (*models.WaitlistEntry, error)
}

type postgresWaitlistRepository struct {
	db *sql.DB
}

var ErrWaitlistEntryExists = errors.New("waitlist entry already exists")

func NewWaitlistRepository(db *sql.DB) WaitlistRepository {
	return &postgresWaitlistRepository{db: db}
}

func (r *postgresWaitlistRepository) Create(ctx context.Context, entry *models.WaitlistEntry) (*models.WaitlistEntry, error) {
	row := r.db.QueryRowContext(ctx,
		`INSERT INTO waitlist_entries (email)
         VALUES ($1)
         ON CONFLICT (email) DO NOTHING
         RETURNING email, joined_at`,
		entry.Email,
	)

	var created models.WaitlistEntry
	if err := row.Scan(&created.Email, &created.JoinedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrWaitlistEntryExists
		}
		return nil, err
	}

	return &created, nil
}

package repositories

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"

	"github.com/nava1525/bilio-backend/internal/app/models"
)

type UserRepository interface {
	List(ctx context.Context) ([]models.User, error)
	GetByID(ctx context.Context, id string) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	Create(ctx context.Context, user *models.User) (*models.User, error)
}

type postgresUserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &postgresUserRepository{db: db}
}

func (r *postgresUserRepository) List(ctx context.Context) ([]models.User, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, email, COALESCE(name, ''), workspace_name, created_at, updated_at 
		 FROM users ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		var workspaceName sql.NullString
		if err := rows.Scan(&user.ID, &user.Email, &user.Name, &workspaceName, &user.CreatedAt, &user.UpdatedAt); err != nil {
			return nil, err
		}
		if workspaceName.Valid {
			user.WorkspaceName = &workspaceName.String
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (r *postgresUserRepository) GetByID(ctx context.Context, id string) (*models.User, error) {
	var user models.User
	var passwordHash, workspaceName sql.NullString

	err := r.db.QueryRowContext(ctx,
		`SELECT id, email, COALESCE(name, ''), password_hash, workspace_name, created_at, updated_at
		 FROM users WHERE id = $1`,
		id).Scan(&user.ID, &user.Email, &user.Name, &passwordHash, &workspaceName, &user.CreatedAt, &user.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	if passwordHash.Valid {
		user.PasswordHash = &passwordHash.String
	}
	if workspaceName.Valid {
		user.WorkspaceName = &workspaceName.String
	}

	return &user, nil
}

func (r *postgresUserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	var passwordHash, workspaceName sql.NullString

	err := r.db.QueryRowContext(ctx,
		`SELECT id, email, COALESCE(name, ''), password_hash, workspace_name, created_at, updated_at
		 FROM users WHERE email = $1`,
		email).Scan(&user.ID, &user.Email, &user.Name, &passwordHash, &workspaceName, &user.CreatedAt, &user.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	if passwordHash.Valid {
		user.PasswordHash = &passwordHash.String
	}
	if workspaceName.Valid {
		user.WorkspaceName = &workspaceName.String
	}

	return &user, nil
}

func (r *postgresUserRepository) Create(ctx context.Context, user *models.User) (*models.User, error) {
	id := uuid.NewString()
	now := time.Now().UTC()

	_, err := r.db.ExecContext(ctx,
		`INSERT INTO users (id, email, name, password_hash, workspace_name, created_at, updated_at) 
		 VALUES ($1, $2, $3, $4, $5, $6, $6)`,
		id, user.Email, user.Name, user.PasswordHash, user.WorkspaceName, now,
	)
	if err != nil {
		return nil, err
	}

	return &models.User{
		ID:           id,
		Email:        user.Email,
		Name:         user.Name,
		WorkspaceName: user.WorkspaceName,
		CreatedAt:    now,
		UpdatedAt:    now,
	}, nil
}

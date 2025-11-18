package services

import (
	"context"
	"fmt"

	"github.com/nava1525/bilio-backend/internal/app/models"
	"github.com/nava1525/bilio-backend/internal/app/repositories"
)

type UserService struct {
	users repositories.UserRepository
}

type CreateUserInput struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

func NewUserService(repo repositories.UserRepository) *UserService {
	return &UserService{users: repo}
}

func (s *UserService) List(ctx context.Context) ([]models.User, error) {
	return s.users.List(ctx)
}

func (s *UserService) Create(ctx context.Context, input CreateUserInput) (*models.User, error) {
	if input.Email == "" {
		return nil, fmt.Errorf("email is required")
	}

	user := &models.User{
		Email: input.Email,
		Name:  input.Name,
	}
	return s.users.Create(ctx, user)
}

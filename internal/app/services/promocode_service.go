package services

import (
	"context"
	"crypto/rand"
	"fmt"
	"time"

	"github.com/nava1525/bilio-backend/internal/app/models"
	"github.com/nava1525/bilio-backend/internal/app/repositories"
)

type PromocodeService struct {
	repo repositories.PromocodeRepository
}

func NewPromocodeService(repo repositories.PromocodeRepository) *PromocodeService {
	return &PromocodeService{
		repo: repo,
	}
}

func (s *PromocodeService) Generate(ctx context.Context) (*models.Promocode, error) {
	const maxRetries = 5
	var lastErr error

	for i := 0; i < maxRetries; i++ {
		// Generate a random 8-character promocode (4 bytes = 8 hex characters)
		bytes := make([]byte, 4)
		if _, err := rand.Read(bytes); err != nil {
			return nil, fmt.Errorf("generate random code: %w", err)
		}
		// Convert to uppercase hex string
		code := fmt.Sprintf("%X", bytes)

		promocode := &models.Promocode{
			Code:      code,
			CreatedAt: time.Now(),
		}

		created, err := s.repo.Create(ctx, promocode)
		if err == nil {
			return created, nil
		}

		// If it's a duplicate key error, retry with a new code
		lastErr = err
	}

	return nil, fmt.Errorf("failed to generate unique promocode after %d attempts: %w", maxRetries, lastErr)
}


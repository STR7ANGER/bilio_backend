package services

import (
	"context"
	"errors"

	"github.com/nava1525/bilio-backend/internal/app/models"
	"github.com/nava1525/bilio-backend/internal/app/repositories"
)

type ClientService struct {
	clients repositories.ClientRepository
}

type CreateClientInput struct {
	Name     string  `json:"name"`
	Email    *string `json:"email,omitempty"`
	Company  *string `json:"company,omitempty"`
	Phone    *string `json:"phone,omitempty"`
	Address  *string `json:"address,omitempty"`
	TaxID    *string `json:"tax_id,omitempty"`
	Currency string  `json:"currency"`
}

type UpdateClientInput struct {
	Name     string  `json:"name"`
	Email    *string `json:"email,omitempty"`
	Company  *string `json:"company,omitempty"`
	Phone    *string `json:"phone,omitempty"`
	Address  *string `json:"address,omitempty"`
	TaxID    *string `json:"tax_id,omitempty"`
	Currency string  `json:"currency"`
}

func NewClientService(clientRepo repositories.ClientRepository) *ClientService {
	return &ClientService{clients: clientRepo}
}

func (s *ClientService) List(ctx context.Context, userID string) ([]models.Client, error) {
	return s.clients.List(ctx, userID)
}

func (s *ClientService) GetByID(ctx context.Context, id string, userID string) (*models.Client, error) {
	client, err := s.clients.GetByID(ctx, id, userID)
	if err != nil {
		return nil, err
	}
	if client == nil {
		return nil, errors.New("client not found")
	}
	return client, nil
}

func (s *ClientService) Create(ctx context.Context, userID string, input CreateClientInput) (*models.Client, error) {
	if input.Name == "" {
		return nil, errors.New("name is required")
	}
	if input.Currency == "" {
		input.Currency = "USD"
	}

	client := &models.Client{
		UserID:   userID,
		Name:     input.Name,
		Email:    input.Email,
		Company:  input.Company,
		Phone:    input.Phone,
		Address:  input.Address,
		TaxID:    input.TaxID,
		Currency: input.Currency,
	}

	return s.clients.Create(ctx, client)
}

func (s *ClientService) Update(ctx context.Context, id string, userID string, input UpdateClientInput) (*models.Client, error) {
	client, err := s.clients.GetByID(ctx, id, userID)
	if err != nil {
		return nil, err
	}
	if client == nil {
		return nil, errors.New("client not found")
	}

	client.Name = input.Name
	client.Email = input.Email
	client.Company = input.Company
	client.Phone = input.Phone
	client.Address = input.Address
	client.TaxID = input.TaxID
	client.Currency = input.Currency

	return s.clients.Update(ctx, client)
}

func (s *ClientService) Delete(ctx context.Context, id string, userID string) error {
	client, err := s.clients.GetByID(ctx, id, userID)
	if err != nil {
		return err
	}
	if client == nil {
		return errors.New("client not found")
	}
	return s.clients.Delete(ctx, id, userID)
}


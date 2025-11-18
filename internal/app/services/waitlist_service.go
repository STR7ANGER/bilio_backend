package services

import (
	"context"
	"errors"
	"fmt"
	"net/mail"
	"strings"

	"github.com/nava1525/bilio-backend/internal/app/models"
	"github.com/nava1525/bilio-backend/internal/app/repositories"
	templates "github.com/nava1525/bilio-backend/internal/templates"
	"github.com/nava1525/bilio-backend/pkg/mailer"
)

type ValidationError struct {
	Message string
}

func (e ValidationError) Error() string {
	return e.Message
}

func newValidationError(message string) error {
	return ValidationError{Message: message}
}

func AsValidationError(err error) (ValidationError, bool) {
	var vErr ValidationError
	if errors.As(err, &vErr) {
		return vErr, true
	}
	return ValidationError{}, false
}

type WaitlistService struct {
	repo            repositories.WaitlistRepository
	promocodeRepo   repositories.PromocodeRepository
	mailer          mailer.Sender
}

type JoinWaitlistInput struct {
	Email     string `json:"email"`
	Promocode string `json:"promocode,omitempty"`
}

func NewWaitlistService(repo repositories.WaitlistRepository, promocodeRepo repositories.PromocodeRepository, sender mailer.Sender) *WaitlistService {
	return &WaitlistService{
		repo:          repo,
		promocodeRepo: promocodeRepo,
		mailer:        sender,
	}
}

func (s *WaitlistService) Join(ctx context.Context, input JoinWaitlistInput) (*models.WaitlistEntry, error) {
	if s.mailer == nil {
		return nil, fmt.Errorf("email sender not configured")
	}

	// Validate promocode if provided (normalize to uppercase for case-insensitive matching)
	if input.Promocode != "" {
		normalizedPromocode := strings.ToUpper(strings.TrimSpace(input.Promocode))
		_, err := s.promocodeRepo.FindAndDelete(ctx, normalizedPromocode)
		if err != nil {
			if errors.Is(err, repositories.ErrPromocodeNotFound) {
				return nil, newValidationError("invalid promocode")
			}
			if errors.Is(err, repositories.ErrPromocodeAlreadyUsed) {
				return nil, newValidationError("promocode already used")
			}
			return nil, fmt.Errorf("validate promocode: %w", err)
		}
	} else {
		return nil, newValidationError("promocode is required")
	}

	normalizedEmail, err := normalizeGmail(input.Email)
	if err != nil {
		return nil, err
	}

	entry := &models.WaitlistEntry{
		Email: normalizedEmail,
	}

	created, err := s.repo.Create(ctx, entry)
	if err != nil {
		if errors.Is(err, repositories.ErrWaitlistEntryExists) {
			return entry, nil
		}
		return nil, err
	}

	subject := "Welcome to the BillStack Waitlist!"
	textBody := `Hi there,

Thanks for raising your hand for the BillStack waitlist—we're excited to help you reclaim your billing workflow.

Here's what you're getting early access to:
- Branded invoice creation with one-click payment links
- Smart reminders that nudge clients automatically
- Expense tracking that ties every rupee and dollar back to clients
- Real-time profitability reports so you can see what’s working

We’ll be in touch soon with early access perks, onboarding, and partner offers designed just for waitlisters.

Talk soon,
Team BillStack`

	msg := mailer.Message{
		To:       created.Email,
		Subject:  subject,
		TextBody: textBody,
		HTMLBody: templates.WelcomeEmailHTML,
	}

	if err := s.mailer.Send(ctx, msg); err != nil {
		return nil, fmt.Errorf("send welcome email: %w", err)
	}

	return created, nil
}

func normalizeGmail(email string) (string, error) {
	trimmed := strings.TrimSpace(strings.ToLower(email))
	if trimmed == "" {
		return "", newValidationError("email is required")
	}

	parsed, err := mail.ParseAddress(trimmed)
	if err != nil {
		return "", newValidationError("invalid email address")
	}

	address := strings.ToLower(parsed.Address)
	if !strings.HasSuffix(address, "@gmail.com") {
		return "", newValidationError("only gmail.com addresses are allowed")
	}

	return address, nil
}

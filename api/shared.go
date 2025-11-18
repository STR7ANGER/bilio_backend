package api

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/nava1525/bilio-backend/internal/app/models"
	"github.com/nava1525/bilio-backend/internal/app/repositories"
	"github.com/nava1525/bilio-backend/internal/app/services"
	"github.com/nava1525/bilio-backend/internal/config"
	"github.com/nava1525/bilio-backend/internal/database"
	"github.com/nava1525/bilio-backend/internal/logger"
	pkgmailer "github.com/nava1525/bilio-backend/pkg/mailer"
)

var (
	initOnce     sync.Once
	sharedConfig *config.Config
	sharedDB     *sql.DB
	sharedLogger logger.Logger

	// Services
	authService      *services.AuthService
	clientService    *services.ClientService
	invoiceService   *services.InvoiceService
	expenseService   *services.ExpenseService
	reportService    *services.ReportService
	waitlistService  *services.WaitlistService
	promocodeService *services.PromocodeService
	userService      *services.UserService
)

func initServices() error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}
	sharedConfig = cfg

	logProvider := logger.New(cfg.Logging.Level)
	sharedLogger = logProvider

	dbClient, err := database.NewClient(cfg.Database.URL)
	if err != nil {
		return fmt.Errorf("initialize database: %w", err)
	}
	sharedDB = dbClient.DB()

	// Repositories
	userRepo := repositories.NewUserRepository(sharedDB)
	clientRepo := repositories.NewClientRepository(sharedDB)
	invoiceRepo := repositories.NewInvoiceRepository(sharedDB)
	expenseRepo := repositories.NewExpenseRepository(sharedDB)
	waitlistRepo := repositories.NewWaitlistRepository(sharedDB)
	promocodeRepo := repositories.NewPromocodeRepository(sharedDB)

	// Services
	authService = services.NewAuthService(userRepo)
	clientService = services.NewClientService(clientRepo)
	invoiceService = services.NewInvoiceService(invoiceRepo, clientRepo)
	expenseService = services.NewExpenseService(expenseRepo, clientRepo)
	reportService = services.NewReportService(invoiceRepo, expenseRepo, clientRepo)
	userService = services.NewUserService(userRepo)

	// Email service
	if cfg.Email.SMTP.Username == "" || cfg.Email.SMTP.Password == "" {
		return fmt.Errorf("email smtp credentials missing; set EMAIL_USER and EMAIL_PASSWORD")
	}

	mailer, err := pkgmailer.NewSMTPMailer(pkgmailer.SMTPConfig{
		Host:     cfg.Email.SMTP.Host,
		Port:     cfg.Email.SMTP.Port,
		Username: cfg.Email.SMTP.Username,
		Password: cfg.Email.SMTP.Password,
		From:     cfg.Email.From,
	})
	if err != nil {
		return fmt.Errorf("initialize mailer: %w", err)
	}

	waitlistService = services.NewWaitlistService(waitlistRepo, promocodeRepo, mailer)
	promocodeService = services.NewPromocodeService(promocodeRepo)

	return nil
}

func EnsureInitialized() error {
	var initErr error
	initOnce.Do(func() {
		initErr = initServices()
	})
	return initErr
}

// Helper functions

func HandleCORS(w http.ResponseWriter, r *http.Request) {
	origins := sharedConfig.CORS.AllowedOrigins
	if len(origins) == 0 {
		origins = []string{"*"}
	}

	origin := r.Header.Get("Origin")
	allowed := false
	for _, o := range origins {
		if o == "*" || o == origin {
			allowed = true
			break
		}
	}

	if allowed {
		w.Header().Set("Access-Control-Allow-Origin", origin)
	}

	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Authorization, Content-Type, X-CSRF-Token")
	w.Header().Set("Access-Control-Expose-Headers", "Link")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Max-Age", "300")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusNoContent)
		return
	}
}

func RespondJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func RespondError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]string{"error": message})
}

func GetUserID(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", fmt.Errorf("missing authorization header")
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return "", fmt.Errorf("invalid authorization header format")
	}

	token := parts[1]
	claims, err := authService.ValidateToken(token)
	if err != nil {
		return "", fmt.Errorf("invalid token: %w", err)
	}

	userID, ok := claims["user_id"].(string)
	if !ok {
		return "", fmt.Errorf("invalid token claims")
	}

	return userID, nil
}

func RequireAuth(w http.ResponseWriter, r *http.Request) (string, bool) {
	userID, err := GetUserID(r)
	if err != nil {
		RespondError(w, http.StatusUnauthorized, err.Error())
		return "", false
	}
	return userID, true
}

// GetAuthService returns the initialized auth service
func GetAuthService() *services.AuthService {
	_ = EnsureInitialized()
	return authService
}

// GetClientService returns the initialized client service
func GetClientService() *services.ClientService {
	_ = EnsureInitialized()
	return clientService
}

// GetInvoiceService returns the initialized invoice service
func GetInvoiceService() *services.InvoiceService {
	_ = EnsureInitialized()
	return invoiceService
}

// GetExpenseService returns the initialized expense service
func GetExpenseService() *services.ExpenseService {
	_ = EnsureInitialized()
	return expenseService
}

// GetReportService returns the initialized report service
func GetReportService() *services.ReportService {
	_ = EnsureInitialized()
	return reportService
}

// GetWaitlistService returns the initialized waitlist service
func GetWaitlistService() *services.WaitlistService {
	_ = EnsureInitialized()
	return waitlistService
}

// GetPromocodeService returns the initialized promocode service
func GetPromocodeService() *services.PromocodeService {
	_ = EnsureInitialized()
	return promocodeService
}

// GetUserService returns the initialized user service
func GetUserService() *services.UserService {
	_ = EnsureInitialized()
	return userService
}

// GetLogger returns the initialized logger
func GetLogger() logger.Logger {
	_ = EnsureInitialized()
	return sharedLogger
}

func ExtractID(r *http.Request) string {
	// Vercel passes path parameters via query or we can parse from URL
	// For now, we'll use a helper that works with Vercel's routing
	path := r.URL.Path
	parts := strings.Split(strings.Trim(path, "/"), "/")
	if len(parts) > 0 {
		// Try to find ID in path - this is a simple approach
		// Vercel may pass it differently
		for i, part := range parts {
			if part == "clients" || part == "invoices" || part == "expenses" {
				if i+1 < len(parts) {
					return parts[i+1]
				}
			}
		}
	}
	return ""
}

func GetContextWithUserID(r *http.Request, userID string) context.Context {
	return context.WithValue(r.Context(), "user_id", userID)
}

// Re-export service types to avoid importing internal packages from function files
type (
	// Auth service types
	LoginInput    = services.LoginInput
	RegisterInput = services.RegisterInput

	// Client service types
	CreateClientInput = services.CreateClientInput
	UpdateClientInput = services.UpdateClientInput

	// Invoice service types
	CreateInvoiceInput = services.CreateInvoiceInput
	UpdateInvoiceInput = services.UpdateInvoiceInput
	InvoiceFilters     = services.InvoiceFilters
	CreatePaymentInput = services.CreatePaymentInput

	// Expense service types
	CreateExpenseInput = services.CreateExpenseInput
	UpdateExpenseInput = services.UpdateExpenseInput
	ExpenseFilters     = services.ExpenseFilters

	// User service types
	CreateUserInput = services.CreateUserInput

	// Waitlist service types
	JoinWaitlistInput = services.JoinWaitlistInput
)

// Re-export model types
type InvoiceStatus = models.InvoiceStatus

// Re-export service functions
func AsValidationError(err error) (services.ValidationError, bool) {
	return services.AsValidationError(err)
}


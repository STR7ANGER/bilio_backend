package transport

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/rs/zerolog"

	appHandlers "github.com/nava1525/bilio-backend/internal/app/handlers"
	appRepositories "github.com/nava1525/bilio-backend/internal/app/repositories"
	appServices "github.com/nava1525/bilio-backend/internal/app/services"
	"github.com/nava1525/bilio-backend/internal/config"
	pkgmailer "github.com/nava1525/bilio-backend/pkg/mailer"
	pkgmiddleware "github.com/nava1525/bilio-backend/pkg/middleware"
)

func NewRouter(cfg *config.Config, logger zerolog.Logger, db *sql.DB) (http.Handler, error) {
	r := chi.NewRouter()

	r.Use(pkgmiddleware.RequestLogger(logger))
	r.Use(pkgmiddleware.Recovery(logger))
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Timeout(60 * time.Second))
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   cfg.CORS.AllowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	healthHandler := appHandlers.NewHealthHandler()
	
	// Repositories
	userRepo := appRepositories.NewUserRepository(db)
	clientRepo := appRepositories.NewClientRepository(db)
	invoiceRepo := appRepositories.NewInvoiceRepository(db)
	expenseRepo := appRepositories.NewExpenseRepository(db)
	waitlistRepo := appRepositories.NewWaitlistRepository(db)
	promocodeRepo := appRepositories.NewPromocodeRepository(db)

	// Services
	authService := appServices.NewAuthService(userRepo)
	clientService := appServices.NewClientService(clientRepo)
	invoiceService := appServices.NewInvoiceService(invoiceRepo, clientRepo)
	expenseService := appServices.NewExpenseService(expenseRepo, clientRepo)
	reportService := appServices.NewReportService(invoiceRepo, expenseRepo, clientRepo)

	// Handlers
	authHandler := appHandlers.NewAuthHandler(authService)
	clientHandler := appHandlers.NewClientHandler(clientService)
	invoiceHandler := appHandlers.NewInvoiceHandler(invoiceService)
	expenseHandler := appHandlers.NewExpenseHandler(expenseService)
	reportHandler := appHandlers.NewReportHandler(reportService)
	userHandler := appHandlers.NewUserHandler(userRepo)

	if cfg.Email.SMTP.Username == "" || cfg.Email.SMTP.Password == "" {
		return nil, fmt.Errorf("email smtp credentials missing; set EMAIL_USER and EMAIL_PASSWORD")
	}

	mailer, err := pkgmailer.NewSMTPMailer(pkgmailer.SMTPConfig{
		Host:     cfg.Email.SMTP.Host,
		Port:     cfg.Email.SMTP.Port,
		Username: cfg.Email.SMTP.Username,
		Password: cfg.Email.SMTP.Password,
		From:     cfg.Email.From,
	})
	if err != nil {
		return nil, err
	}

	waitlistService := appServices.NewWaitlistService(waitlistRepo, promocodeRepo, mailer)
	waitlistHandler := appHandlers.NewWaitlistHandler(waitlistService, logger)

	promocodeService := appServices.NewPromocodeService(promocodeRepo)
	promocodeHandler := appHandlers.NewPromocodeHandler(promocodeService, logger)

	// Auth middleware
	authMiddleware := pkgmiddleware.AuthMiddleware(authService)

	r.Get("/health", healthHandler.Check)

	r.Route("/api/v1", func(r chi.Router) {
		// Public endpoints
		r.Post("/auth/register", authHandler.Register)
		r.Post("/auth/login", authHandler.Login)

		// Legacy endpoints (keep for backward compatibility)
		r.Route("/users", func(r chi.Router) {
			r.Get("/", userHandler.List)
			r.Post("/", userHandler.Create)
		})
		r.Get("/promocode", promocodeHandler.Generate)
		r.Post("/waitlist", waitlistHandler.Join)

		// Protected endpoints - require authentication
		r.Group(func(r chi.Router) {
			r.Use(authMiddleware)

			// Clients
			r.Route("/clients", func(r chi.Router) {
				r.Get("/", clientHandler.List)
				r.Post("/", clientHandler.Create)
				r.Get("/{id}", clientHandler.Get)
				r.Put("/{id}", clientHandler.Update)
				r.Delete("/{id}", clientHandler.Delete)
			})

			// Invoices
			r.Route("/invoices", func(r chi.Router) {
				r.Get("/", invoiceHandler.List)
				r.Post("/", invoiceHandler.Create)
				r.Get("/{id}", invoiceHandler.Get)
				r.Put("/{id}", invoiceHandler.Update)
				r.Post("/{id}/send", invoiceHandler.Send)
				r.Post("/{id}/mark-paid", invoiceHandler.MarkPaid)
				r.Get("/{id}/pdf", invoiceHandler.GetPDF)
			})

			// Expenses
			r.Route("/expenses", func(r chi.Router) {
				r.Get("/", expenseHandler.List)
				r.Post("/", expenseHandler.Create)
				r.Get("/{id}", expenseHandler.Get)
				r.Put("/{id}", expenseHandler.Update)
			})

			// Reports
			r.Route("/reports", func(r chi.Router) {
				r.Get("/summary", reportHandler.GetSummary)
				r.Get("/client-profit/{id}", reportHandler.GetClientProfitability)
				r.Get("/tax-summary", reportHandler.GetTaxSummary)
			})
		})
	})

	return r, nil
}

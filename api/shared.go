package api

// Re-export everything from pkg/api to maintain compatibility
// This allows api/v1/* files to import from github.com/nava1525/bilio-backend/api
// while the actual implementation that imports internal/ is in pkg/api
import (
	apipkg "github.com/nava1525/bilio-backend/pkg/api"
)

// Re-export all functions
var (
	HandleCORS        = apipkg.HandleCORS
	RespondJSON       = apipkg.RespondJSON
	RespondError      = apipkg.RespondError
	GetUserID         = apipkg.GetUserID
	RequireAuth       = apipkg.RequireAuth
	EnsureInitialized = apipkg.EnsureInitialized
	ExtractID         = apipkg.ExtractID
	GetContextWithUserID = apipkg.GetContextWithUserID
	GetAuthService    = apipkg.GetAuthService
	GetClientService  = apipkg.GetClientService
	GetInvoiceService = apipkg.GetInvoiceService
	GetExpenseService = apipkg.GetExpenseService
	GetReportService  = apipkg.GetReportService
	GetWaitlistService = apipkg.GetWaitlistService
	GetPromocodeService = apipkg.GetPromocodeService
	GetUserService    = apipkg.GetUserService
	GetLogger         = apipkg.GetLogger
	AsValidationError = apipkg.AsValidationError
)

// Re-export types
type (
	LoginInput         = apipkg.LoginInput
	RegisterInput      = apipkg.RegisterInput
	CreateClientInput  = apipkg.CreateClientInput
	UpdateClientInput  = apipkg.UpdateClientInput
	CreateInvoiceInput = apipkg.CreateInvoiceInput
	UpdateInvoiceInput = apipkg.UpdateInvoiceInput
	InvoiceFilters     = apipkg.InvoiceFilters
	CreatePaymentInput = apipkg.CreatePaymentInput
	CreateExpenseInput = apipkg.CreateExpenseInput
	UpdateExpenseInput = apipkg.UpdateExpenseInput
	ExpenseFilters     = apipkg.ExpenseFilters
	CreateUserInput    = apipkg.CreateUserInput
	JoinWaitlistInput  = apipkg.JoinWaitlistInput
	InvoiceStatus      = apipkg.InvoiceStatus
)


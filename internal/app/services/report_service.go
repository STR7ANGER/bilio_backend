package services

import (
	"context"
	"errors"
	"time"

	"github.com/nava1525/bilio-backend/internal/app/models"
	"github.com/nava1525/bilio-backend/internal/app/repositories"
)

type ReportService struct {
	invoices repositories.InvoiceRepository
	expenses repositories.ExpenseRepository
	clients  repositories.ClientRepository
}

type SummaryReport struct {
	TotalRevenue    float64 `json:"total_revenue"`
	TotalExpenses   float64 `json:"total_expenses"`
	NetProfit       float64 `json:"net_profit"`
	OutstandingInvoices int `json:"outstanding_invoices"`
	PaidInvoices    int    `json:"paid_invoices"`
	TotalInvoices    int    `json:"total_invoices"`
}

type ClientProfitability struct {
	ClientID      string  `json:"client_id"`
	ClientName    string  `json:"client_name"`
	TotalRevenue  float64 `json:"total_revenue"`
	TotalExpenses float64 `json:"total_expenses"`
	NetProfit     float64 `json:"net_profit"`
	ProfitMargin  float64 `json:"profit_margin"` // percentage
}

type TaxSummary struct {
	Period       string            `json:"period"`
	TotalRevenue float64           `json:"total_revenue"`
	TotalExpenses float64          `json:"total_expenses"`
	NetIncome    float64           `json:"net_income"`
	Invoices     []TaxInvoiceEntry `json:"invoices"`
	Expenses     []TaxExpenseEntry `json:"expenses"`
}

type TaxInvoiceEntry struct {
	InvoiceNumber string    `json:"invoice_number"`
	Date          time.Time `json:"date"`
	ClientName    string    `json:"client_name"`
	Amount        float64   `json:"amount"`
	TaxAmount     float64   `json:"tax_amount"`
}

type TaxExpenseEntry struct {
	Description string    `json:"description"`
	Date        time.Time `json:"date"`
	Category    string    `json:"category"`
	Amount      float64   `json:"amount"`
}

func NewReportService(invoiceRepo repositories.InvoiceRepository, expenseRepo repositories.ExpenseRepository, clientRepo repositories.ClientRepository) *ReportService {
	return &ReportService{
		invoices: invoiceRepo,
		expenses: expenseRepo,
		clients:  clientRepo,
	}
}

func (s *ReportService) GetSummary(ctx context.Context, userID string, fromDate, toDate *time.Time) (*SummaryReport, error) {
	filters := repositories.InvoiceFilters{
		FromDate: fromDate,
		ToDate:   toDate,
	}
	invoices, err := s.invoices.List(ctx, userID, filters)
	if err != nil {
		return nil, err
	}

	expenseFilters := repositories.ExpenseFilters{
		FromDate: fromDate,
		ToDate:   toDate,
	}
	expenses, err := s.expenses.List(ctx, userID, expenseFilters)
	if err != nil {
		return nil, err
	}

	totalRevenue := 0.0
	totalExpenses := 0.0
	outstandingCount := 0
	paidCount := 0

	for _, inv := range invoices {
		if inv.Status == models.InvoiceStatusPaid {
			totalRevenue += inv.Total
			paidCount++
		} else if inv.Status == models.InvoiceStatusPending || inv.Status == models.InvoiceStatusOverdue {
			outstandingCount++
		}
	}

	for _, exp := range expenses {
		totalExpenses += exp.Amount
	}

	return &SummaryReport{
		TotalRevenue:       totalRevenue,
		TotalExpenses:      totalExpenses,
		NetProfit:          totalRevenue - totalExpenses,
		OutstandingInvoices: outstandingCount,
		PaidInvoices:       paidCount,
		TotalInvoices:      len(invoices),
	}, nil
}

func (s *ReportService) GetClientProfitability(ctx context.Context, userID string, clientID string, fromDate, toDate *time.Time) (*ClientProfitability, error) {
	client, err := s.clients.GetByID(ctx, clientID, userID)
	if err != nil {
		return nil, err
	}
	if client == nil {
		return nil, errors.New("client not found")
	}

	clientIDPtr := &clientID
	filters := repositories.InvoiceFilters{
		ClientID: clientIDPtr,
		FromDate: fromDate,
		ToDate:   toDate,
	}
	invoices, err := s.invoices.List(ctx, userID, filters)
	if err != nil {
		return nil, err
	}

	expenseFilters := repositories.ExpenseFilters{
		ClientID: clientIDPtr,
		FromDate: fromDate,
		ToDate:   toDate,
	}
	expenses, err := s.expenses.List(ctx, userID, expenseFilters)
	if err != nil {
		return nil, err
	}

	totalRevenue := 0.0
	for _, inv := range invoices {
		if inv.Status == models.InvoiceStatusPaid {
			totalRevenue += inv.Total
		}
	}

	totalExpenses := 0.0
	for _, exp := range expenses {
		totalExpenses += exp.Amount
	}

	netProfit := totalRevenue - totalExpenses
	profitMargin := 0.0
	if totalRevenue > 0 {
		profitMargin = (netProfit / totalRevenue) * 100
	}

	return &ClientProfitability{
		ClientID:     clientID,
		ClientName:    client.Name,
		TotalRevenue:  totalRevenue,
		TotalExpenses: totalExpenses,
		NetProfit:     netProfit,
		ProfitMargin:  profitMargin,
	}, nil
}

func (s *ReportService) GetTaxSummary(ctx context.Context, userID string, fromDate, toDate time.Time) (*TaxSummary, error) {
	filters := repositories.InvoiceFilters{
		FromDate: &fromDate,
		ToDate:   &toDate,
	}
	invoices, err := s.invoices.List(ctx, userID, filters)
	if err != nil {
		return nil, err
	}

	expenseFilters := repositories.ExpenseFilters{
		FromDate: &fromDate,
		ToDate:   &toDate,
	}
	expenses, err := s.expenses.List(ctx, userID, expenseFilters)
	if err != nil {
		return nil, err
	}

	totalRevenue := 0.0
	totalExpenses := 0.0
	taxInvoices := []TaxInvoiceEntry{}
	taxExpenses := []TaxExpenseEntry{}

	for _, inv := range invoices {
		if inv.Status == models.InvoiceStatusPaid {
			totalRevenue += inv.Total
			taxInvoices = append(taxInvoices, TaxInvoiceEntry{
				InvoiceNumber: inv.InvoiceNumber,
				Date:          inv.IssueDate,
				ClientName:    "", // Would need to join with clients table
				Amount:        inv.Total,
				TaxAmount:     inv.TaxAmount,
			})
		}
	}

	for _, exp := range expenses {
		totalExpenses += exp.Amount
		category := ""
		if exp.Category != nil {
			category = *exp.Category
		}
		taxExpenses = append(taxExpenses, TaxExpenseEntry{
			Description: exp.Description,
			Date:        exp.ExpenseDate,
			Category:    category,
			Amount:      exp.Amount,
		})
	}

	return &TaxSummary{
		Period:       fromDate.Format("2006-01") + " to " + toDate.Format("2006-01"),
		TotalRevenue: totalRevenue,
		TotalExpenses: totalExpenses,
		NetIncome:    totalRevenue - totalExpenses,
		Invoices:     taxInvoices,
		Expenses:     taxExpenses,
	}, nil
}


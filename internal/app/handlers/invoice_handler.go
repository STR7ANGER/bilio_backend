package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	pkgmiddleware "github.com/nava1525/bilio-backend/pkg/middleware"
	"github.com/nava1525/bilio-backend/internal/app/models"
	"github.com/nava1525/bilio-backend/internal/app/services"
)

type InvoiceHandler struct {
	service *services.InvoiceService
}

func NewInvoiceHandler(service *services.InvoiceService) *InvoiceHandler {
	return &InvoiceHandler{service: service}
}

func (h *InvoiceHandler) List(w http.ResponseWriter, r *http.Request) {
	userID := pkgmiddleware.GetUserID(r.Context())
	if userID == "" {
		respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	filters := services.InvoiceFilters{}

	if status := r.URL.Query().Get("status"); status != "" {
		s := models.InvoiceStatus(status)
		filters.Status = &s
	}
	if clientID := r.URL.Query().Get("client_id"); clientID != "" {
		filters.ClientID = &clientID
	}
	if fromDateStr := r.URL.Query().Get("from_date"); fromDateStr != "" {
		if fromDate, err := time.Parse("2006-01-02", fromDateStr); err == nil {
			filters.FromDate = &fromDate
		}
	}
	if toDateStr := r.URL.Query().Get("to_date"); toDateStr != "" {
		if toDate, err := time.Parse("2006-01-02", toDateStr); err == nil {
			filters.ToDate = &toDate
		}
	}

	invoices, err := h.service.List(r.Context(), userID, filters)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, invoices)
}

func (h *InvoiceHandler) Get(w http.ResponseWriter, r *http.Request) {
	userID := pkgmiddleware.GetUserID(r.Context())
	if userID == "" {
		respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	id := chi.URLParam(r, "id")
	invoice, err := h.service.GetByID(r.Context(), id, userID)
	if err != nil {
		respondError(w, http.StatusNotFound, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, invoice)
}

func (h *InvoiceHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID := pkgmiddleware.GetUserID(r.Context())
	if userID == "" {
		respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var input services.CreateInvoiceInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		respondError(w, http.StatusBadRequest, "invalid payload")
		return
	}

	invoice, err := h.service.Create(r.Context(), userID, input)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondJSON(w, http.StatusCreated, invoice)
}

func (h *InvoiceHandler) Update(w http.ResponseWriter, r *http.Request) {
	userID := pkgmiddleware.GetUserID(r.Context())
	if userID == "" {
		respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	id := chi.URLParam(r, "id")
	var input services.UpdateInvoiceInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		respondError(w, http.StatusBadRequest, "invalid payload")
		return
	}

	invoice, err := h.service.Update(r.Context(), id, userID, input)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, invoice)
}

func (h *InvoiceHandler) MarkPaid(w http.ResponseWriter, r *http.Request) {
	userID := pkgmiddleware.GetUserID(r.Context())
	if userID == "" {
		respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	id := chi.URLParam(r, "id")
	var input services.CreatePaymentInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		respondError(w, http.StatusBadRequest, "invalid payload")
		return
	}

	invoice, err := h.service.MarkPaid(r.Context(), id, userID, input)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, invoice)
}

func (h *InvoiceHandler) Send(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement email sending
	respondJSON(w, http.StatusOK, map[string]string{"message": "Invoice send functionality will be implemented"})
}

func (h *InvoiceHandler) GetPDF(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement PDF generation
	respondJSON(w, http.StatusOK, map[string]string{"message": "PDF generation will be implemented"})
}


package handlers

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	pkgmiddleware "github.com/nava1525/bilio-backend/pkg/middleware"
	"github.com/nava1525/bilio-backend/internal/app/services"
)

type ReportHandler struct {
	service *services.ReportService
}

func NewReportHandler(service *services.ReportService) *ReportHandler {
	return &ReportHandler{service: service}
}

func (h *ReportHandler) GetSummary(w http.ResponseWriter, r *http.Request) {
	userID := pkgmiddleware.GetUserID(r.Context())
	if userID == "" {
		respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var fromDate, toDate *time.Time
	if fromDateStr := r.URL.Query().Get("from_date"); fromDateStr != "" {
		if parsed, err := time.Parse("2006-01-02", fromDateStr); err == nil {
			fromDate = &parsed
		}
	}
	if toDateStr := r.URL.Query().Get("to_date"); toDateStr != "" {
		if parsed, err := time.Parse("2006-01-02", toDateStr); err == nil {
			toDate = &parsed
		}
	}

	summary, err := h.service.GetSummary(r.Context(), userID, fromDate, toDate)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, summary)
}

func (h *ReportHandler) GetClientProfitability(w http.ResponseWriter, r *http.Request) {
	userID := pkgmiddleware.GetUserID(r.Context())
	if userID == "" {
		respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	clientID := chi.URLParam(r, "id")
	var fromDate, toDate *time.Time
	if fromDateStr := r.URL.Query().Get("from_date"); fromDateStr != "" {
		if parsed, err := time.Parse("2006-01-02", fromDateStr); err == nil {
			fromDate = &parsed
		}
	}
	if toDateStr := r.URL.Query().Get("to_date"); toDateStr != "" {
		if parsed, err := time.Parse("2006-01-02", toDateStr); err == nil {
			toDate = &parsed
		}
	}

	profitability, err := h.service.GetClientProfitability(r.Context(), userID, clientID, fromDate, toDate)
	if err != nil {
		respondError(w, http.StatusNotFound, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, profitability)
}

func (h *ReportHandler) GetTaxSummary(w http.ResponseWriter, r *http.Request) {
	userID := pkgmiddleware.GetUserID(r.Context())
	if userID == "" {
		respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	fromDateStr := r.URL.Query().Get("from_date")
	toDateStr := r.URL.Query().Get("to_date")

	if fromDateStr == "" || toDateStr == "" {
		respondError(w, http.StatusBadRequest, "from_date and to_date are required")
		return
	}

	fromDate, err := time.Parse("2006-01-02", fromDateStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid from_date format (use YYYY-MM-DD)")
		return
	}

	toDate, err := time.Parse("2006-01-02", toDateStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid to_date format (use YYYY-MM-DD)")
		return
	}

	summary, err := h.service.GetTaxSummary(r.Context(), userID, fromDate, toDate)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, summary)
}


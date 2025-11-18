package invoices

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/nava1525/bilio-backend/api"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	api.HandleCORS(w, r)
	if r.Method == "OPTIONS" {
		return
	}

	if err := api.EnsureInitialized(); err != nil {
		api.RespondError(w, http.StatusInternalServerError, "service initialization failed")
		return
	}

	userID, ok := api.RequireAuth(w, r)
	if !ok {
		return
	}

	id := extractIDFromPath(r.URL.Path)
	if id != "" {
		switch r.Method {
		case http.MethodGet:
			invoice, err := api.GetInvoiceService().GetByID(r.Context(), id, userID)
			if err != nil {
				api.RespondError(w, http.StatusNotFound, err.Error())
				return
			}
			api.RespondJSON(w, http.StatusOK, invoice)
		case http.MethodPut:
			var input api.UpdateInvoiceInput
			if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
				api.RespondError(w, http.StatusBadRequest, "invalid payload")
				return
			}

			invoice, err := api.GetInvoiceService().Update(r.Context(), id, userID, input)
			if err != nil {
				api.RespondError(w, http.StatusBadRequest, err.Error())
				return
			}
			api.RespondJSON(w, http.StatusOK, invoice)
		default:
			api.RespondError(w, http.StatusMethodNotAllowed, "method not allowed")
		}
		return
	}

	switch r.Method {
	case http.MethodGet:
		filters := api.InvoiceFilters{}

		if status := r.URL.Query().Get("status"); status != "" {
			s := api.InvoiceStatus(status)
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

		invoices, err := api.GetInvoiceService().List(r.Context(), userID, filters)
		if err != nil {
			api.RespondError(w, http.StatusInternalServerError, err.Error())
			return
		}
		api.RespondJSON(w, http.StatusOK, invoices)
	case http.MethodPost:
		var input api.CreateInvoiceInput
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			api.RespondError(w, http.StatusBadRequest, "invalid payload")
			return
		}

		invoice, err := api.GetInvoiceService().Create(r.Context(), userID, input)
		if err != nil {
			api.RespondError(w, http.StatusBadRequest, err.Error())
			return
		}
		api.RespondJSON(w, http.StatusCreated, invoice)
	default:
		api.RespondError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

func extractIDFromPath(path string) string {
	parts := strings.Split(strings.Trim(path, "/"), "/")
	for i, part := range parts {
		if part == "invoices" && i+1 < len(parts) {
			nextPart := parts[i+1]
			if nextPart != "index" && nextPart != "" && nextPart != "handler" {
				return nextPart
			}
		}
	}
	return ""
}


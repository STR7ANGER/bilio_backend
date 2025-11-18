package expenses

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/nava1525/bilio-backend/api"
	"github.com/nava1525/bilio-backend/internal/app/services"
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
			expense, err := api.GetExpenseService().GetByID(r.Context(), id, userID)
			if err != nil {
				api.RespondError(w, http.StatusNotFound, err.Error())
				return
			}
			api.RespondJSON(w, http.StatusOK, expense)
		case http.MethodPut:
			var input services.UpdateExpenseInput
			if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
				api.RespondError(w, http.StatusBadRequest, "invalid payload")
				return
			}

			expense, err := api.GetExpenseService().Update(r.Context(), id, userID, input)
			if err != nil {
				api.RespondError(w, http.StatusBadRequest, err.Error())
				return
			}
			api.RespondJSON(w, http.StatusOK, expense)
		default:
			api.RespondError(w, http.StatusMethodNotAllowed, "method not allowed")
		}
		return
	}

	switch r.Method {
	case http.MethodGet:
		filters := services.ExpenseFilters{}

		if clientID := r.URL.Query().Get("client_id"); clientID != "" {
			filters.ClientID = &clientID
		}
		if category := r.URL.Query().Get("category"); category != "" {
			filters.Category = &category
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

		expenses, err := api.GetExpenseService().List(r.Context(), userID, filters)
		if err != nil {
			api.RespondError(w, http.StatusInternalServerError, err.Error())
			return
		}
		api.RespondJSON(w, http.StatusOK, expenses)
	case http.MethodPost:
		var input services.CreateExpenseInput
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			api.RespondError(w, http.StatusBadRequest, "invalid payload")
			return
		}

		expense, err := api.GetExpenseService().Create(r.Context(), userID, input)
		if err != nil {
			api.RespondError(w, http.StatusBadRequest, err.Error())
			return
		}
		api.RespondJSON(w, http.StatusCreated, expense)
	default:
		api.RespondError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

func extractIDFromPath(path string) string {
	parts := strings.Split(strings.Trim(path, "/"), "/")
	for i, part := range parts {
		if part == "expenses" && i+1 < len(parts) {
			nextPart := parts[i+1]
			if nextPart != "index" && nextPart != "" {
				return nextPart
			}
		}
	}
	return ""
}


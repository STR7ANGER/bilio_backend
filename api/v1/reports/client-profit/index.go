package reports

import (
	"net/http"
	"strings"
	"time"

	"github.com/nava1525/bilio-backend/pkg/api"
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

	if r.Method != http.MethodGet {
		api.RespondError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	clientID := extractIDFromPath(r.URL.Path)

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

	profitability, err := api.GetReportService().GetClientProfitability(r.Context(), userID, clientID, fromDate, toDate)
	if err != nil {
		api.RespondError(w, http.StatusNotFound, err.Error())
		return
	}

	api.RespondJSON(w, http.StatusOK, profitability)
}

func extractIDFromPath(path string) string {
	parts := strings.Split(strings.Trim(path, "/"), "/")
	for i, part := range parts {
		if part == "client-profit" && i+1 < len(parts) {
			return parts[i+1]
		}
	}
	return ""
}


package summary

import (
	"net/http"
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

	summary, err := api.GetReportService().GetSummary(r.Context(), userID, fromDate, toDate)
	if err != nil {
		api.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	api.RespondJSON(w, http.StatusOK, summary)
}


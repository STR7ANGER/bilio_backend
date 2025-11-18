package reports

import (
	"net/http"
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

	if r.Method != http.MethodGet {
		api.RespondError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	fromDateStr := r.URL.Query().Get("from_date")
	toDateStr := r.URL.Query().Get("to_date")

	if fromDateStr == "" || toDateStr == "" {
		api.RespondError(w, http.StatusBadRequest, "from_date and to_date are required")
		return
	}

	fromDate, err := time.Parse("2006-01-02", fromDateStr)
	if err != nil {
		api.RespondError(w, http.StatusBadRequest, "invalid from_date format (use YYYY-MM-DD)")
		return
	}

	toDate, err := time.Parse("2006-01-02", toDateStr)
	if err != nil {
		api.RespondError(w, http.StatusBadRequest, "invalid to_date format (use YYYY-MM-DD)")
		return
	}

	summary, err := api.GetReportService().GetTaxSummary(r.Context(), userID, fromDate, toDate)
	if err != nil {
		api.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	api.RespondJSON(w, http.StatusOK, summary)
}


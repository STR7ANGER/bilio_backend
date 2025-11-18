package api

import (
	"net/http"

	"github.com/nava1525/bilio-backend/pkg/api"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	api.HandleCORS(w, r)
	if r.Method == "OPTIONS" {
		return
	}

	if r.Method != http.MethodGet {
		api.RespondError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	api.RespondJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}


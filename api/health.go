package api

import (
	"net/http"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	HandleCORS(w, r)
	if r.Method == "OPTIONS" {
		return
	}

	if r.Method != http.MethodGet {
		RespondError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	RespondJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}


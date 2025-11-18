package auth

import (
	"encoding/json"
	"net/http"

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

	if r.Method != http.MethodPost {
		api.RespondError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var input api.RegisterInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		api.RespondError(w, http.StatusBadRequest, "invalid payload")
		return
	}

	result, err := api.GetAuthService().Register(r.Context(), input)
	if err != nil {
		api.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}

	api.RespondJSON(w, http.StatusCreated, result)
}


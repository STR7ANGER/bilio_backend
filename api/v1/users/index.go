package users

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

	switch r.Method {
	case http.MethodGet:
		users, err := api.GetUserService().List(r.Context())
		if err != nil {
			api.RespondError(w, http.StatusInternalServerError, err.Error())
			return
		}
		api.RespondJSON(w, http.StatusOK, users)
	case http.MethodPost:
		var input api.CreateUserInput
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			api.RespondError(w, http.StatusBadRequest, "invalid payload")
			return
		}

		user, err := api.GetUserService().Create(r.Context(), input)
		if err != nil {
			api.RespondError(w, http.StatusInternalServerError, err.Error())
			return
		}
		api.RespondJSON(w, http.StatusCreated, user)
	default:
		api.RespondError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}


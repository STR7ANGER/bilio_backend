package waitlist

import (
	"encoding/json"
	"net/http"

	"github.com/nava1525/bilio-backend/api"
)

type waitlistJoinRequest struct {
	Email     string `json:"email"`
	Promocode string `json:"promocode,omitempty"`
}

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

	var req waitlistJoinRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		api.RespondError(w, http.StatusBadRequest, "invalid payload")
		return
	}

	entry, err := api.GetWaitlistService().Join(r.Context(), api.JoinWaitlistInput{
		Email:     req.Email,
		Promocode: req.Promocode,
	})
	if err != nil {
		if validationErr, ok := api.AsValidationError(err); ok {
			api.RespondError(w, http.StatusBadRequest, validationErr.Message)
			return
		}

		logger := api.GetLogger()
		logger.Error().Err(err).Msg("waitlist join failed")
		api.RespondError(w, http.StatusInternalServerError, "failed to join waitlist")
		return
	}

	response := map[string]interface{}{
		"email":   entry.Email,
		"message": "welcome email sent",
	}

	api.RespondJSON(w, http.StatusCreated, response)
}


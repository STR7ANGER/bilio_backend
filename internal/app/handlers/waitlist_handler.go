package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/rs/zerolog"

	"github.com/nava1525/bilio-backend/internal/app/services"
)

type WaitlistHandler struct {
	service *services.WaitlistService
	logger  zerolog.Logger
}

type waitlistJoinRequest struct {
	Email     string `json:"email"`
	Promocode string `json:"promocode,omitempty"`
}

func NewWaitlistHandler(service *services.WaitlistService, logger zerolog.Logger) *WaitlistHandler {
	return &WaitlistHandler{
		service: service,
		logger:  logger,
	}
}

func (h *WaitlistHandler) Join(w http.ResponseWriter, r *http.Request) {
	var req waitlistJoinRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid payload", http.StatusBadRequest)
		return
	}

	entry, err := h.service.Join(r.Context(), services.JoinWaitlistInput{
		Email:     req.Email,
		Promocode: req.Promocode,
	})
	if err != nil {
		if validationErr, ok := services.AsValidationError(err); ok {
			http.Error(w, validationErr.Message, http.StatusBadRequest)
			return
		}

		h.logger.Error().Err(err).Msg("waitlist join failed")
		http.Error(w, "failed to join waitlist", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"email":   entry.Email,
		"message": "welcome email sent",
	}

	respondJSON(w, http.StatusCreated, response)
}

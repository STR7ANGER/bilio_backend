package handlers

import (
	"net/http"

	"github.com/rs/zerolog"

	"github.com/nava1525/bilio-backend/internal/app/services"
)

type PromocodeHandler struct {
	service *services.PromocodeService
	logger  zerolog.Logger
}

func NewPromocodeHandler(service *services.PromocodeService, logger zerolog.Logger) *PromocodeHandler {
	return &PromocodeHandler{
		service: service,
		logger:  logger,
	}
}

func (h *PromocodeHandler) Generate(w http.ResponseWriter, r *http.Request) {
	promocode, err := h.service.Generate(r.Context())
	if err != nil {
		h.logger.Error().Err(err).Msg("promocode generation failed")
		http.Error(w, "failed to generate promocode", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"code":       promocode.Code,
		"created_at": promocode.CreatedAt,
	}

	respondJSON(w, http.StatusOK, response)
}


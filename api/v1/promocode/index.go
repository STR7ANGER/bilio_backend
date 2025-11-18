package promocode

import (
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

	if r.Method != http.MethodGet {
		api.RespondError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	promocode, err := api.GetPromocodeService().Generate(r.Context())
	if err != nil {
		logger := api.GetLogger()
		logger.Error().Err(err).Msg("promocode generation failed")
		api.RespondError(w, http.StatusInternalServerError, "failed to generate promocode")
		return
	}

	response := map[string]interface{}{
		"code":       promocode.Code,
		"created_at": promocode.CreatedAt,
	}

	api.RespondJSON(w, http.StatusOK, response)
}


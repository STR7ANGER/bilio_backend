package clients

import (
	"encoding/json"
	"net/http"
	"strings"

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

	// Check if this is an ID operation (path contains an ID after /clients/)
	id := extractIDFromPath(r.URL.Path)
	if id != "" {
		// Handle ID-based operations
		switch r.Method {
		case http.MethodGet:
			client, err := api.GetClientService().GetByID(r.Context(), id, userID)
			if err != nil {
				api.RespondError(w, http.StatusNotFound, err.Error())
				return
			}
			api.RespondJSON(w, http.StatusOK, client)
		case http.MethodPut:
			var input api.UpdateClientInput
			if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
				api.RespondError(w, http.StatusBadRequest, "invalid payload")
				return
			}

			client, err := api.GetClientService().Update(r.Context(), id, userID, input)
			if err != nil {
				api.RespondError(w, http.StatusBadRequest, err.Error())
				return
			}
			api.RespondJSON(w, http.StatusOK, client)
		case http.MethodDelete:
			if err := api.GetClientService().Delete(r.Context(), id, userID); err != nil {
				api.RespondError(w, http.StatusNotFound, err.Error())
				return
			}
			w.WriteHeader(http.StatusNoContent)
		default:
			api.RespondError(w, http.StatusMethodNotAllowed, "method not allowed")
		}
		return
	}

	// Handle collection operations
	switch r.Method {
	case http.MethodGet:
		clients, err := api.GetClientService().List(r.Context(), userID)
		if err != nil {
			api.RespondError(w, http.StatusInternalServerError, err.Error())
			return
		}
		api.RespondJSON(w, http.StatusOK, clients)
	case http.MethodPost:
		var input api.CreateClientInput
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			api.RespondError(w, http.StatusBadRequest, "invalid payload")
			return
		}

		client, err := api.GetClientService().Create(r.Context(), userID, input)
		if err != nil {
			api.RespondError(w, http.StatusBadRequest, err.Error())
			return
		}
		api.RespondJSON(w, http.StatusCreated, client)
	default:
		api.RespondError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

func extractIDFromPath(path string) string {
	parts := strings.Split(strings.Trim(path, "/"), "/")
	for i, part := range parts {
		if part == "clients" && i+1 < len(parts) {
			nextPart := parts[i+1]
			// Don't treat "index" as an ID
			if nextPart != "index" && nextPart != "" {
				return nextPart
			}
		}
	}
	return ""
}


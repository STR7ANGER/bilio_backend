package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	pkgmiddleware "github.com/nava1525/bilio-backend/pkg/middleware"
	"github.com/nava1525/bilio-backend/internal/app/services"
)

type ExpenseHandler struct {
	service *services.ExpenseService
}

func NewExpenseHandler(service *services.ExpenseService) *ExpenseHandler {
	return &ExpenseHandler{service: service}
}

func (h *ExpenseHandler) List(w http.ResponseWriter, r *http.Request) {
	userID := pkgmiddleware.GetUserID(r.Context())
	if userID == "" {
		respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	filters := services.ExpenseFilters{}

	if clientID := r.URL.Query().Get("client_id"); clientID != "" {
		filters.ClientID = &clientID
	}
	if category := r.URL.Query().Get("category"); category != "" {
		filters.Category = &category
	}
	if fromDateStr := r.URL.Query().Get("from_date"); fromDateStr != "" {
		if fromDate, err := time.Parse("2006-01-02", fromDateStr); err == nil {
			filters.FromDate = &fromDate
		}
	}
	if toDateStr := r.URL.Query().Get("to_date"); toDateStr != "" {
		if toDate, err := time.Parse("2006-01-02", toDateStr); err == nil {
			filters.ToDate = &toDate
		}
	}

	expenses, err := h.service.List(r.Context(), userID, filters)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, expenses)
}

func (h *ExpenseHandler) Get(w http.ResponseWriter, r *http.Request) {
	userID := pkgmiddleware.GetUserID(r.Context())
	if userID == "" {
		respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	id := chi.URLParam(r, "id")
	expense, err := h.service.GetByID(r.Context(), id, userID)
	if err != nil {
		respondError(w, http.StatusNotFound, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, expense)
}

func (h *ExpenseHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID := pkgmiddleware.GetUserID(r.Context())
	if userID == "" {
		respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var input services.CreateExpenseInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		respondError(w, http.StatusBadRequest, "invalid payload")
		return
	}

	expense, err := h.service.Create(r.Context(), userID, input)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondJSON(w, http.StatusCreated, expense)
}

func (h *ExpenseHandler) Update(w http.ResponseWriter, r *http.Request) {
	userID := pkgmiddleware.GetUserID(r.Context())
	if userID == "" {
		respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	id := chi.URLParam(r, "id")
	var input services.UpdateExpenseInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		respondError(w, http.StatusBadRequest, "invalid payload")
		return
	}

	expense, err := h.service.Update(r.Context(), id, userID, input)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, expense)
}


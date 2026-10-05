package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/ednanf/school-api/internal/domain"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

type staffHandler struct {
	service  domain.StaffService
	validate *validator.Validate
}

func NewStaffHandler(service domain.StaffService, validate *validator.Validate) *staffHandler {
	return &staffHandler{service: service, validate: validate}
}

func (h *staffHandler) StaffRoutes() chi.Router {
	r := chi.NewRouter()

	r.Get("/", h.HandleList)
	r.Post("/", h.HandleCreate)
	r.Delete("/{id}", h.HandleDelete)
	r.Get("/{id}", h.HandleGetByID)
	r.Patch("/{id}", h.HandleUpdate)

	return r
}

func (h *staffHandler) HandleCreate(w http.ResponseWriter, r *http.Request) {
	var employee domain.Staff

	// Decode
	if err := json.NewDecoder(r.Body).Decode(&employee); err != nil {
		fmt.Println(err)
		sendError(w, http.StatusBadRequest, "Invalid JSON payload", nil)
		return
	}

	// Validate
	if err := h.validate.StructCtx(r.Context(), &employee); err != nil {
		if validationErrs, ok := err.(validator.ValidationErrors); ok {
			sendError(w, http.StatusUnprocessableEntity, "Validation failed", formatValidationErrors(validationErrs))
			return
		}
		sendError(w, http.StatusBadRequest, "Validation failed", nil)
		return
	}

	// Delegate normalizationg, timestamping and persistence to the service layer
	if err := h.service.Create(r.Context(), &employee); err != nil {
		sendError(w, http.StatusInternalServerError, "Failed to create staff member entry", nil)
		return
	}

	sendSuccess(w, http.StatusCreated, "Staff member created successfully", employee)
}

func (h *staffHandler) HandleDelete(w http.ResponseWriter, r *http.Request) {
	// Extract and convert url param
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil || id <= 0 {
		sendError(w, http.StatusBadRequest, "Invalid staff member ID", nil)
		return
	}

	if err := h.service.Delete(r.Context(), id); err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			sendError(w, http.StatusNotFound, "Staff member not found", nil)
			return
		}
		sendError(w, http.StatusInternalServerError, "Failed to delete staff member", nil)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *staffHandler) HandleGetByID(w http.ResponseWriter, r *http.Request) {
	sendSuccess(w, http.StatusOK, "get by id hit", nil)
}

func (h *staffHandler) HandleList(w http.ResponseWriter, r *http.Request) {
	sendSuccess(w, http.StatusOK, "list hit", nil)
}

func (h *staffHandler) HandleUpdate(w http.ResponseWriter, r *http.Request) {
	sendSuccess(w, http.StatusOK, "update hit", nil)
}

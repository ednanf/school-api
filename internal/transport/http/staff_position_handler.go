package http

import (
	"encoding/json"
	"net/http"

	"github.com/ednanf/school-api/internal/domain"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

type StaffPositionHandler struct {
	service  domain.StaffPositionService
	validate *validator.Validate
}

func NewStaffPositionHandler(service domain.StaffPositionService, validate *validator.Validate) *StaffPositionHandler {
	return &StaffPositionHandler{service: service, validate: validate}
}

func (h *StaffPositionHandler) StaffPositionRoutes() chi.Router {
	r := chi.NewRouter()

	r.Get("/", h.HandleList)
	r.Post("/", h.HandleCreate)
	r.Delete("/{id}", h.HandleDelete)
	r.Get("/{id}", h.HandleGetByID)
	r.Patch("/{id}", h.HandleUpdate)

	return r
}

func (h *StaffPositionHandler) HandleCreate(w http.ResponseWriter, r *http.Request) {
	var position domain.StaffPosition

	if err := json.NewDecoder(r.Body).Decode(&position); err != nil {
		sendError(w, http.StatusBadRequest, "Invalid JSON payload", nil)
		return
	}

	// Validate struct rules
	if err := h.validate.StructCtx(r.Context(), &position); err != nil {
		if validationsErrs, ok := err.(validator.ValidationErrors); ok {
			sendError(w, http.StatusUnprocessableEntity, "Validations failed", formatValidationErrors(validationsErrs))
			return
		}
		sendError(w, http.StatusBadRequest, "Validation failed", nil)
		return
	}

	// Delegate normalization, timing and persistence to service layer
	if err := h.service.Create(r.Context(), &position); err != nil {
		sendError(w, http.StatusInternalServerError, "Failed to create position entry", nil)
		return
	}

	sendSuccess(w, http.StatusCreated, "Position created successfully", position)
}

func (h *StaffPositionHandler) HandleDelete(w http.ResponseWriter, r *http.Request) {
	sendSuccess(w, http.StatusOK, "delete hit", nil)
}

func (h *StaffPositionHandler) HandleGetByID(w http.ResponseWriter, r *http.Request) {
	sendSuccess(w, http.StatusOK, "get by id hit", nil)
}

func (h *StaffPositionHandler) HandleList(w http.ResponseWriter, r *http.Request) {
	sendSuccess(w, http.StatusOK, "list hit", nil)
}

func (h *StaffPositionHandler) HandleUpdate(w http.ResponseWriter, r *http.Request) {
	sendSuccess(w, http.StatusOK, "update hit", nil)
}

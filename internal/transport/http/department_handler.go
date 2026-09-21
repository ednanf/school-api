package http

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ednanf/school-api/internal/domain"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

type DepartmentHandler struct {
	service  domain.DepartmentService
	validate *validator.Validate
}

func NewDepartmentHandler(service domain.DepartmentService, validate *validator.Validate) *DepartmentHandler {
	return &DepartmentHandler{service: service, validate: validate}
}

func (h *DepartmentHandler) DepartmentRoutes() chi.Router {
	r := chi.NewRouter()

	r.Get("/", h.HandleList)
	r.Post("/", h.HandleCreate)
	r.Delete("/{id}", h.HandleDelete)
	r.Get("/{id}", h.HandleGetById)
	r.Patch("/{id}", h.HandleUpdate)

	return r
}

func (h *DepartmentHandler) HandleCreate(w http.ResponseWriter, r *http.Request) {
	var department domain.Department

	if err := json.NewDecoder(r.Body).Decode(&department); err != nil {
		sendError(w, http.StatusBadRequest, "Invalid JSON payload", nil)
		return
	}

	if err := h.validate.StructCtx(r.Context(), &department); err != nil {
		if validationErrs, ok := err.(validator.ValidationErrors); ok {
			sendError(w, http.StatusUnprocessableEntity, "Validation failed", formatValidationErrors(validationErrs))
			return
		}
		sendError(w, http.StatusBadRequest, "Validation failed", nil)
		return
	}

	// Calls service now, which runs Normalize() + Timestamps -> then calls Repo.Create()
	if err := h.service.Create(r.Context(), &department); err != nil {
		fmt.Println(err)
		sendError(w, http.StatusInternalServerError, "Failed to create department entry", nil)
		return
	}

	sendSuccess(w, http.StatusCreated, "Department created successfully", department)
}

func (h *DepartmentHandler) HandleDelete(w http.ResponseWriter, r *http.Request) {
	sendSuccess(w, http.StatusOK, "delete hit", nil)
}

func (h *DepartmentHandler) HandleGetById(w http.ResponseWriter, r *http.Request) {
	sendSuccess(w, http.StatusOK, "get by id hit", nil)
}

func (h *DepartmentHandler) HandleList(w http.ResponseWriter, r *http.Request) {
	sendSuccess(w, http.StatusOK, "list hit", nil)
}

func (h *DepartmentHandler) HandleUpdate(w http.ResponseWriter, r *http.Request) {
	sendSuccess(w, http.StatusOK, "update hit", nil)
}

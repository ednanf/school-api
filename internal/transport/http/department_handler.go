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

	if err := h.service.Create(r.Context(), &department); err != nil {
		fmt.Println(err)
		sendError(w, http.StatusInternalServerError, "Failed to create department entry", nil)
		return
	}

	sendSuccess(w, http.StatusCreated, "Department created successfully", department)
}

func (h *DepartmentHandler) HandleDelete(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		sendError(w, http.StatusBadRequest, "Invalid department ID", nil)
		return
	}

	if err := h.service.Delete(r.Context(), id); err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			sendError(w, http.StatusNotFound, "Department not found", nil)
			return
		}
		sendError(w, http.StatusInternalServerError, "Failed to delete department", nil)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *DepartmentHandler) HandleGetById(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		sendError(w, http.StatusBadRequest, "Invalid department id", nil)
		return
	}

	dept, err := h.service.GetById(r.Context(), id)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			sendError(w, http.StatusNotFound, "Department not found", nil)
			return
		}
		sendError(w, http.StatusInternalServerError, "Failed to retrieve department", nil)
		return
	}

	sendSuccess(w, http.StatusOK, "Department retrieved successfully", dept)
}

func (h *DepartmentHandler) HandleList(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()

	reqPage, _ := strconv.Atoi(queryParams.Get("page"))
	reqLimit, _ := strconv.Atoi(queryParams.Get("limit"))

	// Service returns the actual normalized page and limit used
	departments, totalItems, page, limit, err := h.service.List(r.Context(), reqPage, reqLimit)
	if err != nil {
		sendError(w, http.StatusInternalServerError, "Failed to fetch departments", nil)
		return
	}

	totalPages := 0
	if totalItems > 0 {
		totalPages = (totalItems + limit - 1) / limit
	}

	meta := PaginatedMeta{
		Page:       page,
		Limit:      limit,
		Count:      len(departments),
		TotalItems: totalItems,
		TotalPages: totalPages,
	}

	sendPaginated(w, http.StatusOK, departments, meta)
}

func (h *DepartmentHandler) HandleUpdate(w http.ResponseWriter, r *http.Request) {
	// Extract and convert the id
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		sendError(w, http.StatusBadRequest, "Invalid department ID", nil)
		return
	}

	// Decode the request
	var input domain.PatchDepartmentInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		sendError(w, http.StatusBadRequest, "Invalid JSON payload", nil)
		return
	}

	// Validate the payload
	if err := h.validate.StructCtx(r.Context(), &input); err != nil {
		if validationErrs, ok := err.(validator.ValidationErrors); ok {
			sendError(w, http.StatusUnprocessableEntity, "Validation failed", formatValidationErrors(validationErrs))
			return
		}
		sendError(w, http.StatusBadRequest, "Validation failed", nil)
		return
	}

	// Execute the service
	updatedDept, err := h.service.Update(r.Context(), id, input)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			sendError(w, http.StatusNotFound, "Department not found", nil)
			return
		}
		sendError(w, http.StatusInternalServerError, "Failed to update department", nil)
		return
	}

	sendSuccess(w, http.StatusOK, "Department updated successfully", updatedDept)
}

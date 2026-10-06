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
	// Extract and conver the url param
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil || id <= 0 {
		sendError(w, http.StatusBadRequest, "Invalid staff member ID", nil)
		return
	}

	// Delegate fetching to the service layer
	employee, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			sendError(w, http.StatusNotFound, "Staff member not found", nil)
			return
		}
		sendError(w, http.StatusInternalServerError, "Failed to fetch staff member", nil)
		return
	}

	sendSuccess(w, http.StatusOK, "Staff member retrieved successfully", employee)
}

func (h *staffHandler) HandleList(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()

	// Parse page
	page := 1
	if pageStr := queryParams.Get("page"); pageStr != "" {
		if parsedPage, err := strconv.Atoi(pageStr); err == nil && parsedPage > 0 {
			page = parsedPage
		}
	}

	// Parse limit
	limit := 10
	if limitStr := queryParams.Get("limit"); limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err != nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	// Derive SQL offset
	offset := (page - 1) * limit

	// Delegate paginated fetch and count to service layer
	staff, totalItems, err := h.service.List(r.Context(), limit, offset)
	if err != nil {
		sendError(w, http.StatusInternalServerError, "Failed to fetch staff members", nil)
		return
	}

	// Calculate total pages
	totalPages := 0
	if totalItems > 0 {
		totalPages = (totalItems + limit - 1) / limit
	}

	// Construct metadata
	meta := PaginatedMeta{
		Page:       page,
		Limit:      limit,
		Count:      len(staff),
		TotalItems: totalItems,
		TotalPages: totalPages,
	}

	sendPaginated(w, http.StatusOK, staff, meta)
}

func (h *staffHandler) HandleUpdate(w http.ResponseWriter, r *http.Request) {
	// Extract and convert url param
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil || id <= 0 {
		sendError(w, http.StatusBadRequest, "Invalid staff member id", nil)
		return
	}

	// Decode the HTTP body
	var input domain.PatchStaffInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		sendError(w, http.StatusBadRequest, "Invalid JSON payload", nil)
		return
	}

	// Ensure there's at least one field to HandleUpdate
	if !input.HasUpdates() {
		sendError(w, http.StatusBadRequest, "At least one field must be provided for update", nil)
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

	// Delegate business logic to the service layer
	updatedEmployee, err := h.service.Update(r.Context(), id, input)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			sendError(w, http.StatusNotFound, "Staff member not found", nil)
			return
		}
		sendError(w, http.StatusInternalServerError, "Failed to update staff member", nil)
		return
	}

	sendSuccess(w, http.StatusOK, "Staff member updated successfully", updatedEmployee)
}

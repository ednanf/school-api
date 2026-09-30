package http

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

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
	// Extract and convert url param
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		sendError(w, http.StatusBadRequest, "Invalid position ID", nil)
		return
	}

	// Delegate deletion to service layer
	if err := h.service.Delete(r.Context(), id); err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			sendError(w, http.StatusNotFound, "Position not found", nil)
			return
		}

		sendError(w, http.StatusInternalServerError, "Failed to delete position", nil)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *StaffPositionHandler) HandleGetByID(w http.ResponseWriter, r *http.Request) {
	// Extract and convert the url param
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		sendError(w, http.StatusBadRequest, "Invalid position ID", nil)
		return
	}

	// Delegate fetching to the service layer
	position, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			sendError(w, http.StatusNotFound, "Position not found", nil)
		}

		sendError(w, http.StatusInternalServerError, "Failed to fetch position", nil)
		return
	}

	sendSuccess(w, http.StatusOK, "Position retrieved successfully", position)
}

func (h *StaffPositionHandler) HandleList(w http.ResponseWriter, r *http.Request) {
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

	// Delegate paginated fetch and total count to service layer
	positions, totalItems, err := h.service.List(r.Context(), limit, offset)
	if err != nil {
		log.Printf("[DEBUG] %s", err)
		sendError(w, http.StatusInternalServerError, "Failed to fetch staff positions list", nil)
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
		Count:      len(positions),
		TotalItems: totalItems,
		TotalPages: totalPages,
	}

	sendPaginated(w, http.StatusOK, positions, meta)
}

func (h *StaffPositionHandler) HandleUpdate(w http.ResponseWriter, r *http.Request) {
	sendSuccess(w, http.StatusOK, "update hit", nil)
}

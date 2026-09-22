package http

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/ednanf/school-api/internal/domain"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

// ClassHandler contains `repo` with a way to communicate with the database and the pointer to the validator instantiated once in `main.go`
type ClassHandler struct {
	service  domain.ClassService
	validate *validator.Validate
}

// NewClassHandler is a constructor that returns a pointer to a ClassHandler struct, initializing it with the injected repository and validator dependencies
func NewClassHandler(service domain.ClassService, validate *validator.Validate) *ClassHandler {
	return &ClassHandler{service: service, validate: validate}
}

// Route paths
func (h *ClassHandler) ClassRoutes() chi.Router {
	r := chi.NewRouter()

	r.Get("/", h.HandleList)
	r.Post("/", h.HandleCreate)
	r.Get("/{id}/students", h.HandleListStudentsByClassId)
	r.Delete("/{id}", h.HandleDelete)
	r.Get("/{id}", h.HandleGetById)
	r.Patch("/{id}", h.HandleUpdate)

	return r
}

func (h *ClassHandler) HandleCreate(w http.ResponseWriter, r *http.Request) {
	// Initialize a Class struct
	var class domain.Class

	// Decode the JSON body directly into the struct via pointer
	if err := json.NewDecoder(r.Body).Decode(&class); err != nil {
		sendError(w, http.StatusBadRequest, "Invalid JSON payload", nil)
		return
	}

	// Validate struct using the injected validator instance
	if err := h.validate.StructCtx(r.Context(), &class); err != nil {
		if validationErrs, ok := err.(validator.ValidationErrors); ok {
			sendError(w, http.StatusUnprocessableEntity, "Validation failed", formatValidationErrors(validationErrs))
			return
		}
		sendError(w, http.StatusBadRequest, "Validation failed", nil)
		return
	}

	// Save to the db via the repository
	if err := h.service.Create(r.Context(), &class); err != nil {
		sendError(w, http.StatusInternalServerError, "Failed to create class entry", nil)
		return
	}

	// Return 201 with the full record (includes message)
	sendSuccess(w, http.StatusCreated, "Class created successfully", class)
}

func (h *ClassHandler) HandleDelete(w http.ResponseWriter, r *http.Request) {
	// Extract and convert the URL param to int
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		sendError(w, http.StatusBadRequest, "Invalid class ID", nil)
		return
	}

	// Execute the db operation
	if err := h.service.Delete(r.Context(), id); err != nil {
		// 404 when class was not found
		if errors.Is(err, domain.ErrNotFound) {
			sendError(w, http.StatusNotFound, "Class not found", nil)
			return
		}

		// 500 for db connection or syntax errors
		sendError(w, http.StatusInternalServerError, "Failed to delete student", nil)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *ClassHandler) HandleGetById(w http.ResponseWriter, r *http.Request) {
	// Extract and convert the URL param
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		sendError(w, http.StatusBadRequest, "Invalid class ID", nil)
		return
	}

	// Search for class
	class, err := h.service.GetById(r.Context(), id)
	if err != nil {
		sendError(w, http.StatusInternalServerError, "Database error", nil)
		return
	}

	// If the student does not exist
	if class == nil {
		sendError(w, http.StatusNotFound, "Class not found", nil)
		return
	}

	sendSuccess(w, http.StatusOK, "", class)
}

func (h *ClassHandler) HandleList(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()

	reqPage, _ := strconv.Atoi(queryParams.Get("page"))
	reqLimit, _ := strconv.Atoi(queryParams.Get("limit"))

	// Service returns the actual normalized page and limit used
	classes, totalItems, page, limit, err := h.service.List(r.Context(), reqPage, reqLimit)
	if err != nil {
		sendError(w, http.StatusInternalServerError, "Failed to fetch students", nil)
		return
	}

	// Calculate total pages safely
	totalPages := 0
	if totalItems > 0 {
		totalPages = (totalItems + limit - 1) / limit
	}

	meta := PaginatedMeta{
		Page:       page,
		Limit:      limit,
		Count:      len(classes),
		TotalItems: totalItems,
		TotalPages: totalPages,
	}

	sendPaginated(w, http.StatusOK, classes, meta)
}

func (h *ClassHandler) HandleListStudentsByClassId(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	classID, err := strconv.Atoi(idStr)
	if err != nil || classID <= 0 {
		sendError(w, http.StatusBadRequest, "Invalid class ID", nil)
		return
	}

	queryParams := r.URL.Query()
	reqPage, _ := strconv.Atoi(queryParams.Get("page"))
	reqLimit, _ := strconv.Atoi(queryParams.Get("limit"))

	students, totalItems, page, limit, err := h.service.ListStudentsByClassId(r.Context(), classID, reqPage, reqLimit)
	if err != nil {
		sendError(w, http.StatusInternalServerError, "Failed to fetch students", nil)
		return
	}

	totalPages := 0
	if totalItems > 0 {
		totalPages = (totalItems + limit - 1) / limit
	}

	meta := PaginatedMeta{
		Page:       page,
		Limit:      limit,
		Count:      len(students),
		TotalItems: totalItems,
		TotalPages: totalPages,
	}

	sendPaginated(w, http.StatusOK, students, meta)
}

func (h *ClassHandler) HandleUpdate(w http.ResponseWriter, r *http.Request) {
	// Extract and convert the id to int
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		sendError(w, http.StatusBadRequest, "Invalid class ID", nil)
		return
	}

	// Decode the request's body into the pointer-based PATCH DTO
	var input domain.PatchClassInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		sendError(w, http.StatusBadRequest, "Invalid JSON payload", nil)
		return
	}

	// Validate optional field constraints (using omitempty rules)
	if err := h.validate.StructCtx(r.Context(), &input); err != nil {
		if validationErrs, ok := err.(validator.ValidationErrors); ok {
			sendError(w, http.StatusUnprocessableEntity, "Validation failed", formatValidationErrors(validationErrs))
			return
		}
		sendError(w, http.StatusBadRequest, "Validation failed", nil)
		return
	}

	// Perform the update in the db
	updatedClass, err := h.service.Update(r.Context(), id, input)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			sendError(w, http.StatusNotFound, "Class not found", nil)
			return
		}
		sendError(w, http.StatusInternalServerError, "Failed to update the class", nil)
		return
	}

	sendSuccess(w, http.StatusOK, "", updatedClass)
}

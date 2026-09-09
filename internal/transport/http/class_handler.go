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
	repo     domain.ClassRepository
	validate *validator.Validate
}

// NewClassHandler is a constructor that returns a pointer to a ClassHandler struct, initializing it with the injected repository and validator dependencies
func NewClassHandler(repo domain.ClassRepository, validate *validator.Validate) *ClassHandler {
	return &ClassHandler{repo: repo, validate: validate}
}

// Route paths
func (h *ClassHandler) ClassRoutes() chi.Router {
	r := chi.NewRouter()

	r.Get("/", h.HandleList)
	r.Post("/", h.HandleCreate)
	r.Delete("/{id}", h.HandleDelete)
	r.Get("/{id}", h.HandleGetById)
	r.Patch("/{id}", h.HandlePatch)

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
	if err := h.repo.Create(r.Context(), &class); err != nil {
		sendError(w, http.StatusInternalServerError, "Failed to create class entry", nil)
		return
	}

	// Return 201 with the full record
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
	if err := h.repo.Delete(r.Context(), id); err != nil {
		// 404 when class was not found
		if errors.Is(err, sql.ErrNoRows) {
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
	class, err := h.repo.GetById(r.Context(), id)
	if err != nil {
		sendError(w, http.StatusInternalServerError, "Database error", nil)
		return
	}

	// If the student does not exist
	if class == nil {
		sendError(w, http.StatusNotFound, "Class not found", nil)
		return
	}

	sendSuccess(w, http.StatusOK, "Class retrieved successfully", class)
}

func (h *ClassHandler) HandleList(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters correctly from r.URL.Query()
	queryParams := r.URL.Query()
	limitStr := queryParams.Get("limit")
	pageStr := queryParams.Get("page")

	// Defaults
	limit := 20
	page := 1

	// Convert string query params to integers
	if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
		limit = l
	}

	if p, err := strconv.Atoi(pageStr); err == nil && p >= 0 {
		page = p
	}

	// Limit cap to prevent abuse
	if limit > 100 {
		limit = 100
	}

	// Calculate the database offset derived from page number
	offset := (page - 1) * limit

	// Call the repository with context and parsed pagination
	classes, err := h.repo.List(r.Context(), limit, offset)
	if err != nil {
		sendError(w, http.StatusInternalServerError, "Failed to fetch students", nil)
		return
	}

	sendSuccess(w, http.StatusOK, "Fetched classes successfully", classes)
}

func (h *ClassHandler) HandlePatch(w http.ResponseWriter, r *http.Request) {
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
	updatedClass, err := h.repo.Update(r.Context(), id, input)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			sendError(w, http.StatusNotFound, "Class not found", nil)
			return
		}
		sendError(w, http.StatusInternalServerError, "Failed to update the class", nil)
		return
	}

	sendSuccess(w, http.StatusOK, "Class updated successfully", updatedClass)
}

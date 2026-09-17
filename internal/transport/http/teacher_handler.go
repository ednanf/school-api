package http

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/ednanf/school-api/internal/domain"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

// TeacherHandler contains `repo` with a way to communicate with the database and the pointer to the validator instantiated in `main.go`
type TeacherHandler struct {
	repo     domain.TeacherRepository
	validate *validator.Validate
}

// NewTeacherHandler is a constructor that returns a pointer to a TeacherHandler struct, initializing it with the injected repository and validator dependencies
func NewTeacherHandler(repo domain.TeacherRepository, validate *validator.Validate) *TeacherHandler {
	return &TeacherHandler{repo: repo, validate: validate}
}

// Route paths
func (h *TeacherHandler) TeacherRoutes() chi.Router {
	r := chi.NewRouter()

	r.Post("/", h.HandleCreate)
	r.Get("/", h.HandleList)
	r.Delete("/{id}", h.HandleDelete)
	r.Get("/{id}", h.HandleGetById)
	r.Patch("/{id}", h.HandleUpdate)

	return r
}

func (h *TeacherHandler) HandleCreate(w http.ResponseWriter, r *http.Request) {
	// Initialize a Teacher struct
	var teacher domain.Teacher

	// Decode the JSON body directly into the struct via pointer
	if err := json.NewDecoder(r.Body).Decode(&teacher); err != nil {
		sendError(w, http.StatusBadRequest, "Invalid JSON payload", nil)
		return
	}

	// Validate struct rules using the injected validator instance
	if err := h.validate.StructCtx(r.Context(), &teacher); err != nil {
		if validationErrs, ok := err.(validator.ValidationErrors); ok {
			sendError(w, http.StatusUnprocessableEntity, "Validation failed", formatValidationErrors(validationErrs))
			return
		}
		sendError(w, http.StatusBadRequest, "Validation failed", nil)
		return
	}

	// Save to the db via the repository
	if err := h.repo.Create(r.Context(), &teacher); err != nil {
		sendError(w, http.StatusInternalServerError, "Failed to create teacher entry", nil)
		return
	}

	sendSuccess(w, http.StatusCreated, "Teacher created successfully", teacher)
}

func (h *TeacherHandler) HandleDelete(w http.ResponseWriter, r *http.Request) {
	// Extract and convert the URL param into int
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		sendError(w, http.StatusBadRequest, "Invalid student ID", nil)
		return
	}

	// Execute the db operation
	if err := h.repo.Delete(r.Context(), id); err != nil {
		// 404 when teacher was not found
		if errors.Is(err, sql.ErrNoRows) {
			sendError(w, http.StatusNotFound, "Student not found", nil)
			return
		}

		// 500 for db connection or syntax errors
		sendError(w, http.StatusInternalServerError, "Failed to delete student", nil)
		return
	}

	// 204 status requires no JSON body
	w.WriteHeader(http.StatusNoContent)
}

func (h *TeacherHandler) HandleGetById(w http.ResponseWriter, r *http.Request) {
	// Obtain id and convert to int
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		sendError(w, http.StatusBadRequest, "Invalid teacher ID", nil)
		return
	}

	// Search for the teacher
	teacher, err := h.repo.GetById(r.Context(), id)
	if err != nil {
		sendError(w, http.StatusInternalServerError, "Database error", nil)
		return
	}

	// If the teacher does not exist
	if teacher == nil {
		sendError(w, http.StatusNotFound, "Teacher not found", nil)
		return
	}

	sendSuccess(w, http.StatusOK, "", teacher)
}

func (h *TeacherHandler) HandleList(w http.ResponseWriter, r *http.Request) {
	// Parse query params
	queryParams := r.URL.Query()
	limitStr := queryParams.Get("limit")
	pageStr := queryParams.Get("page")

	// Defaults
	limit := 10
	page := 1

	// Convert string params to integers
	if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
		limit = l
	}

	// Limit cap
	if limit > 100 {
		limit = 100
	}

	if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
		page = p
	}

	// Calculate the db offset
	offset := (page - 1) * limit

	// Execute db operation
	teachers, totalItems, err := h.repo.List(r.Context(), limit, offset)
	if err != nil {
		fmt.Printf("[DEBUG] ERROR: %v", err)
		sendError(w, http.StatusInternalServerError, "Failed to fetch teachers", nil)
		return
	}

	// Calculate total pages
	totalPages := 0
	if totalItems > 0 {
		totalPages = (totalItems + limit - 1) / limit
	}

	meta := PaginatedMeta{
		Page:       page,
		Limit:      limit,
		Count:      len(teachers),
		TotalItems: totalItems,
		TotalPages: totalPages,
	}

	sendPaginated(w, http.StatusOK, teachers, meta)
}

func (h *TeacherHandler) HandleUpdate(w http.ResponseWriter, r *http.Request) {
	// Extract and convert the id to int
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		sendError(w, http.StatusBadRequest, "Invalid teacher ID", nil)
		return
	}

	// Decode the request body
	var input domain.PatchTeacherInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		sendError(w, http.StatusBadRequest, "Invalid JSON payload", nil)
		return
	}

	// Validate optional field constraints
	if err := h.validate.StructCtx(r.Context(), &input); err != nil {
		if validationErrs, ok := err.(validator.ValidationErrors); ok {
			sendError(w, http.StatusUnprocessableEntity, "Validation failed", formatValidationErrors(validationErrs))
			return
		}
		sendError(w, http.StatusBadRequest, "Validation failed", nil)
		return
	}

	// Perform the update
	updatedTeacher, err := h.repo.Update(r.Context(), id, input)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			sendError(w, http.StatusNotFound, "Teacher not found", nil)
			return
		}
		sendError(w, http.StatusInternalServerError, "Failed to update teacher", nil)
		return
	}

	// Return 200 with the updated record
	sendSuccess(w, http.StatusOK, "", updatedTeacher)
}

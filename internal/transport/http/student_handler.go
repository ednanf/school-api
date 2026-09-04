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

// StudentHandler holds routes
type StudentHandler struct {
	repo     domain.StudentRepository
	validate *validator.Validate
}

// Constructor
func NewStudentHandler(repo domain.StudentRepository, validate *validator.Validate) *StudentHandler {
	return &StudentHandler{repo: repo, validate: validator.New()}
}

// Route paths
func (h *StudentHandler) StudentRoutes() chi.Router {
	r := chi.NewRouter()

	r.Get("/", h.HandleList)
	r.Post("/", h.HandleCreate)
	r.Delete("/batch", h.HandleBatchDelete)
	r.Post("/batch", h.HandleBatchCreate)
	r.Delete("/{id}", h.HandleDelete)
	r.Get("/{id}", h.HandleGetByID)
	r.Patch("/{id}", h.HandlePatch)

	return r
}

func (h *StudentHandler) HandleBatchCreate(w http.ResponseWriter, r *http.Request) {
	// Instantiate a variable to hold the payload according to the DTO (for validation and structure)
	var input domain.BatchCreateInput

	// Decode the body and store into the variable
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

	// Create the new students
	createdStudents, err := h.repo.BatchCreate(r.Context(), input.Students)
	if err != nil {
		sendError(w, http.StatusInternalServerError, "Failed to create students", nil)
		return
	}

	// Return 201 with the created students
	sendSuccess(w, http.StatusCreated, "Batch creation successful", createdStudents)
}

func (h *StudentHandler) HandleBatchDelete(w http.ResponseWriter, r *http.Request) {
	// Instantiate a variable to hold the IDs to be deleted
	var input domain.BatchDeleteInput

	// Decode the body and store in the variable
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

	// BatchDelete will return the number of rows affected by the query
	deletedCount, err := h.repo.BatchDelete(r.Context(), input.IDs)
	if err != nil {
		sendError(w, http.StatusInternalServerError, "Failed to delete students in batch", nil)
		return
	}

	// Send a simple 200 with the number of deleted students
	sendSuccess(w, http.StatusOK, "Batch deletion successful", map[string]any{
		"deleted_count": deletedCount,
	})
}

func (h *StudentHandler) HandleCreate(w http.ResponseWriter, r *http.Request) {
	// Initialize a Student struct
	var student domain.Student

	// Decode the JSON body directly into the struct pointer
	if err := json.NewDecoder(r.Body).Decode(&student); err != nil {
		sendError(w, http.StatusBadRequest, "Invalid JSON payload", nil)
		return
	}

	// Validate struct rules using the injected validator instance
	if err := h.validate.StructCtx(r.Context(), &student); err != nil {
		if validationErrs, ok := err.(validator.ValidationErrors); ok {
			sendError(w, http.StatusUnprocessableEntity, "Validation failed", formatValidationErrors(validationErrs))
			return
		}
		sendError(w, http.StatusBadRequest, "Validation failed", nil)
		return
	}

	// Save to MariaDB via the repository
	if err := h.repo.Create(r.Context(), &student); err != nil {
		sendError(w, http.StatusInternalServerError, "Failed to create student entry", nil)
		return
	}

	// Return 201 Created with the full record (including generated ID and timestamps)
	sendSuccess(w, http.StatusCreated, "Student created successfully", student)
}

func (h *StudentHandler) HandleDelete(w http.ResponseWriter, r *http.Request) {
	// Extract and convert the URL param to int
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		sendError(w, http.StatusBadRequest, "Invalid student ID", nil)
		return
	}

	err = h.repo.Delete(r.Context(), id)
	if err != nil {
		// 404 when student was not found
		if errors.Is(err, sql.ErrNoRows) {
			sendError(w, http.StatusNotFound, "Student not found", nil)
			return
		}

		// 500 for DB connection or syntax errors
		sendError(w, http.StatusInternalServerError, "Failed to delete student", nil)
		return
	}

	// 204 status requires no JSON body, therefore, no `sendSuccess`
	w.WriteHeader(http.StatusNoContent)
}

func (h *StudentHandler) HandleGetByID(w http.ResponseWriter, r *http.Request) {
	// Get id from URL
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		sendError(w, http.StatusBadRequest, "Invalid student ID", nil)
		return
	}

	// Search for student
	student, err := h.repo.GetByID(r.Context(), id)
	if err != nil {
		sendError(w, http.StatusInternalServerError, "Database error", nil)
		return
	}

	// If the student does not exist
	if student == nil {
		sendError(w, http.StatusNotFound, "Student not found", nil)
		return
	}

	sendSuccess(w, http.StatusOK, "Student retrieved successfully", student)
}

func (h *StudentHandler) HandleList(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters correctly from r.URL.Query()
	queryParams := r.URL.Query()
	limitStr := queryParams.Get("limit")
	pageStr := queryParams.Get("page")

	// Defaults
	limit := 10
	page := 1

	// Convert string query params to integers
	if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
		limit = l
	}

	// Limit cap to prevent someone from requesting 1 million rows at once
	if limit > 100 {
		limit = 100
	}

	if p, err := strconv.Atoi(pageStr); err == nil && p >= 0 {
		page = p
	}

	// Calculate the database offset derived from page number
	offset := (page - 1) * limit

	// Call the repository with context and parsed pagination
	students, err := h.repo.List(r.Context(), limit, offset)
	if err != nil {
		sendError(w, http.StatusInternalServerError, "Failed to fetch students", nil)
		return
	}

	// Returns [] instead of null if empty because studentRepo initializes an empty slice
	sendSuccess(w, http.StatusOK, "Fetched students successfully", students)
}

func (h *StudentHandler) HandlePatch(w http.ResponseWriter, r *http.Request) {
	// Extract and convert the id to int
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		sendError(w, http.StatusBadRequest, "Invalid student ID", nil)
		return
	}

	// Decode the request's body into the pointer-based PATCH DTO
	var input domain.PatchStudentInput
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
	updatedStudent, err := h.repo.Update(r.Context(), id, input)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			sendError(w, http.StatusNotFound, "Student not found", nil)
			return
		}
		sendError(w, http.StatusInternalServerError, "Failed to update student", nil)
		return
	}

	// Return 200 with the updated record
	sendSuccess(w, http.StatusOK, "Student updated successfully", updatedStudent)
}

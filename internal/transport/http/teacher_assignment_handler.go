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

// TeacherAssignmentHandler contains `repo` with a way to communicate with the database and the pointer to the validator instantiated once in `main.go`
type TeacherAssignmentHandler struct {
	repo     domain.TeacherAssignmentRepository
	validate *validator.Validate
}

// NewTeacherAssignmentHandler is a constructor that returns a pointer to a TeacherAssignmentHandler struct, initializing it with the injected repository and validator dependencies
func NewTeacherAssignmentHandler(repo domain.TeacherAssignmentRepository, validate *validator.Validate) *TeacherAssignmentHandler {
	return &TeacherAssignmentHandler{repo: repo, validate: validate}
}

// Route paths
func (h *TeacherAssignmentHandler) TeacherAssignmentRoutes() chi.Router {
	r := chi.NewRouter()

	r.Get("/", h.HandleList)
	r.Post("/", h.HandleCreate)
	r.Delete("/{id}", h.HandleDelete)
	r.Get("/{id}", h.HandleGetById)
	r.Patch("/{id}", h.HandleUpdate)

	return r
}

func (h *TeacherAssignmentHandler) HandleCreate(w http.ResponseWriter, r *http.Request) {
	// Initialize a struct
	var assignment domain.TeacherAssignment

	// Decode the JSON body directly into the struct
	if err := json.NewDecoder(r.Body).Decode(&assignment); err != nil {
		sendError(w, http.StatusBadRequest, "Invalid JSON payload", nil)
		return
	}

	// Validate struct rules using injected validator instance
	if err := h.validate.StructCtx(r.Context(), &assignment); err != nil {
		if validationErrs, ok := err.(validator.ValidationErrors); ok {
			sendError(w, http.StatusUnprocessableEntity, "Validation failed", formatValidationErrors(validationErrs))
			return
		}
		sendError(w, http.StatusBadRequest, "Validation failed", nil)
		return
	}

	// Save to the db via the repository
	if err := h.repo.Create(r.Context(), &assignment); err != nil {
		sendError(w, http.StatusInternalServerError, "Failed to create assignment entry", nil)
		return
	}

	// Return 201 Created with the full record
	sendSuccess(w, http.StatusCreated, "Assignment created successfully", assignment)
}

func (h *TeacherAssignmentHandler) HandleDelete(w http.ResponseWriter, r *http.Request) {
	// Extract and convert the URL param to int
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		sendError(w, http.StatusBadRequest, "Invalid assignment ID", nil)
		return
	}

	// Execute the db operation
	if err = h.repo.Delete(r.Context(), id); err != nil {
		// 404 when not found
		if errors.Is(err, sql.ErrNoRows) {
			sendError(w, http.StatusNotFound, "Assignment not found", nil)
			return
		}
		// 500 for other errors
		sendError(w, http.StatusInternalServerError, "Failed to delete assignment", nil)
		return
	}

	// 204 requires no JSON body
	w.WriteHeader(http.StatusNoContent)
}

func (h *TeacherAssignmentHandler) HandleGetById(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	sendSuccess(w, http.StatusOK, "get by id hit", idStr)
}

func (h *TeacherAssignmentHandler) HandleList(w http.ResponseWriter, r *http.Request) {
	sendSuccess(w, http.StatusOK, "list hit", nil)
}

func (h *TeacherAssignmentHandler) HandleUpdate(w http.ResponseWriter, r *http.Request) {
	sendSuccess(w, http.StatusOK, "update hit", nil)
}

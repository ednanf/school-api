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
	idStr := chi.URLParam(r, "id")
	sendSuccess(w, http.StatusOK, "Get by id hit", idStr)
}

func (h *TeacherHandler) HandleList(w http.ResponseWriter, r *http.Request) {
	sendSuccess(w, http.StatusOK, "List hit", nil)
}

func (h *TeacherHandler) HandleUpdate(w http.ResponseWriter, r *http.Request) {
	sendSuccess(w, http.StatusOK, "Update hit", nil)
}

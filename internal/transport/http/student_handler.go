package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/ednanf/school-api/internal/domain"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

// StudentHandler holds routes.
type StudentHandler struct {
	repo     domain.StudentRepository
	validate *validator.Validate
}

// Constructor
func NewStudentHandler(repo domain.StudentRepository, validate *validator.Validate) *StudentHandler {
	return &StudentHandler{repo: repo, validate: validator.New()}
}

// RegisterRoutes groups all student endpoints together, similar to an Express sub-router.
func (h *StudentHandler) RegisterRoutes(r chi.Router) {
	r.Route("/students", func(r chi.Router) {
		r.Post("/", h.HandleCreate)
		r.Get("/{id}", h.HandleGetByID)
	})
}

func (h *StudentHandler) HandleList(w http.ResponseWriter, r *http.Request) {
	sendJSON(w, http.StatusOK, map[string]string{"message": "Student list hit"})
}

func (h *StudentHandler) HandleCreate(w http.ResponseWriter, r *http.Request) {
	// Initialize a Student struct
	var s domain.Student

	// Decode the JSON body directly into the struct pointer
	if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
		sendJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid JSON payload"})
		return
	}

	// Validate struct rules using the injected validator instance
	if err := h.validate.StructCtx(r.Context(), &s); err != nil {
		if validationErrs, ok := err.(validator.ValidationErrors); ok {
			sendJSON(w, http.StatusUnprocessableEntity, map[string]any{
				"error":   "Validation failed",
				"details": formatValidationErrors(validationErrs),
			})
			return
		}
		sendJSON(w, http.StatusBadRequest, map[string]string{"error": "Validation failed"})
		return
	}

	// Save to MariaDB via the repository
	if err := h.repo.Create(r.Context(), &s); err != nil {
		sendJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to create student"})
		return
	}

	// Return 201 Created with the full record (including generated ID and timestamps)
	sendJSON(w, http.StatusCreated, s)
}

func (h *StudentHandler) HandleGetByID(w http.ResponseWriter, r *http.Request) {
	// Get id from URL
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		sendJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid student ID"})
		return
	}

	// Search for student
	student, err := h.repo.GetByID(r.Context(), id)
	if err != nil {
		sendJSON(w, http.StatusInternalServerError, map[string]string{"error": "Database error"})
		return
	}

	// If the student does not exist
	if student == nil {
		sendJSON(w, http.StatusNotFound, map[string]string{"error": "Student not found"})
		return
	}

	sendJSON(w, http.StatusOK, student)
}

func (h *StudentHandler) HandlePatch(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	sendJSON(w, http.StatusOK, map[string]string{"message": "Patch student hit", "id": id})
}

func (h *StudentHandler) HandleDelete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	sendJSON(w, http.StatusOK, map[string]string{"message": "Delete student hit", "id": id})
}

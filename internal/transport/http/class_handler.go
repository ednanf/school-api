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

	fmt.Printf("ID %v, Grade %v, Letter %v, Created at %v, Updated at %v\n", class.ID, class.Grade, class.Letter, class.CreatedAt, class.UpdatedAt)

	// Decode the JSON body directly into the struct via pointer
	if err := json.NewDecoder(r.Body).Decode(&class); err != nil {
		sendError(w, http.StatusBadRequest, "Invalid JSON payload", nil)
		return
	}

	fmt.Printf("ID %v, Grade %v, Letter %v, Created at %v, Updated at %v\n", class.ID, class.Grade, class.Letter, class.CreatedAt, class.UpdatedAt)

	// Validate struct rules using the injected validator instance
	if err := h.validate.StructCtx(r.Context(), &class); err != nil {
		if validationErrs, ok := err.(validator.ValidationErrors); ok {
			sendError(w, http.StatusUnprocessableEntity, "Validation failed", formatValidationErrors(validationErrs))
			return
		}
		sendError(w, http.StatusBadRequest, "Validation failed", nil)
		return
	}

	fmt.Printf("ID %v, Grade %v, Letter %v, Created at %v, Updated at %v\n", class.ID, class.Grade, class.Letter, class.CreatedAt, class.UpdatedAt)

	// Save to the db via the repository
	if err := h.repo.Create(r.Context(), &class); err != nil {
		fmt.Printf("Database error: %v\n", err) // Print exact error to console
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
	id := chi.URLParam(r, "id")
	sendSuccess(w, http.StatusOK, "GetById hit", id)
}

func (h *ClassHandler) HandleList(w http.ResponseWriter, r *http.Request) {
	sendSuccess(w, http.StatusOK, "List hit", nil)
}

func (h *ClassHandler) HandlePatch(w http.ResponseWriter, r *http.Request) {
	sendSuccess(w, http.StatusOK, "Patch hit", nil)
}

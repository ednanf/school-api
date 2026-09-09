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

// SubjectHandler contains `repo` with a way to communicate with the database and the pointer to the validator instantiated once in `main.go`
type SubjectHandler struct {
	repo     domain.SubjectRepository
	validate *validator.Validate
}

// NewSubjectHandler isa constructor that returns a pointer to a SubjectHandler struct, initializing it with the injected repository and validator dependencies
func NewSubjectHandler(repo domain.SubjectRepository, validate *validator.Validate) *SubjectHandler {
	return &SubjectHandler{repo: repo, validate: validate}
}

// Route paths
func (h *SubjectHandler) SubjectRoutes() chi.Router {
	r := chi.NewRouter()

	r.Get("/", h.HandleList)
	r.Post("/", h.HandleCreate)
	r.Delete("/{id}", h.HandleDelete)
	r.Get("/{id}", h.HandleGetById)
	r.Patch("/{id}", h.HandlePatch)

	return r
}

func (h *SubjectHandler) HandleCreate(w http.ResponseWriter, r *http.Request) {
	// Initialize a Subject struct
	var subject domain.Subject

	// Decode the JSON body directly into the struct via pointer
	if err := json.NewDecoder(r.Body).Decode(&subject); err != nil {
		sendError(w, http.StatusBadRequest, "Invalid JSON payload", nil)
		return
	}

	if err := h.validate.StructCtx(r.Context(), &subject); err != nil {
		if validationErrs, ok := err.(validator.ValidationErrors); ok {
			sendError(w, http.StatusUnprocessableEntity, "Validation failed", formatValidationErrors(validationErrs))
			return
		}
		sendError(w, http.StatusBadRequest, "Validation failed", nil)
		return
	}

	// Save to the db via the repository
	if err := h.repo.Create(r.Context(), &subject); err != nil {
		fmt.Printf("[DEBUG] %v\n", err)
		sendError(w, http.StatusInternalServerError, "Failed to create subject entry", nil)
		return
	}

	sendSuccess(w, http.StatusCreated, "Subject created successfully", subject)
}

func (h *SubjectHandler) HandleDelete(w http.ResponseWriter, r *http.Request) {
	// Extract and convert the URL param to int
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		sendError(w, http.StatusBadRequest, "Invalid subject ID", nil)
		return
	}

	// Execute the db operation
	if err := h.repo.Delete(r.Context(), id); err != nil {
		// 404 when subject is not found
		if errors.Is(err, sql.ErrNoRows) {
			sendError(w, http.StatusNotFound, "Subject not found", nil)
			return
		}

		// 500 for db connection or syntax errors
		sendError(w, http.StatusInternalServerError, "Failed to delete student", nil)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *SubjectHandler) HandleGetById(w http.ResponseWriter, r *http.Request) {
	// Extract and convert the URL param
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		sendError(w, http.StatusBadRequest, "Invalid subject id", nil)
		return
	}

	// Search for subject
	subject, err := h.repo.GetById(r.Context(), id)
	if err != nil {
		sendError(w, http.StatusInternalServerError, "Database error", nil)
		return
	}

	// If the subject does not exist
	if subject == nil {
		sendError(w, http.StatusNotFound, "Subject not found", nil)
		return
	}

	sendSuccess(w, http.StatusOK, "Subject retrieved successfully", subject)
}

func (h *SubjectHandler) HandleList(w http.ResponseWriter, r *http.Request) {
	// Parse query params correctly from r.URL.Query()
	queryParams := r.URL.Query()
	limitStr := queryParams.Get("limit")
	pageStr := queryParams.Get("page")

	// Defaults
	limit := 30
	page := 1

	// Convert string query params to integers
	if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
		limit = 1
	}

	if p, err := strconv.Atoi(pageStr); err == nil && p >= 0 {
		page = p
	}

	// Limit cpa to prevent abuse
	if limit > 100 {
		limit = 100
	}

	// Calculate the database offset derived from page number
	offset := (page - 1) * limit

	// Call the repository with context and parsed pagination
	subjects, err := h.repo.List(r.Context(), limit, offset)
	if err != nil {
		sendError(w, http.StatusInternalServerError, "Failed to fetch students", nil)
		return
	}

	sendSuccess(w, http.StatusOK, "Fetched subjects successfully", subjects)
}

func (h *SubjectHandler) HandlePatch(w http.ResponseWriter, r *http.Request) {
	sendSuccess(w, http.StatusOK, "Patch hit", nil)
}

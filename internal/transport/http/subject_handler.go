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

// SubjectHandler contains `repo` with a way to communicate with the database and the pointer to the validator instantiated once in `main.go`
type SubjectHandler struct {
	service  domain.SubjectService
	validate *validator.Validate
}

// NewSubjectHandler isa constructor that returns a pointer to a SubjectHandler struct, initializing it with the injected repository and validator dependencies
func NewSubjectHandler(service domain.SubjectService, validate *validator.Validate) *SubjectHandler {
	return &SubjectHandler{service: service, validate: validate}
}

// Route paths
func (h *SubjectHandler) SubjectRoutes() chi.Router {
	r := chi.NewRouter()

	r.Get("/", h.HandleList)
	r.Post("/", h.HandleCreate)
	r.Delete("/{id}", h.HandleDelete)
	r.Get("/{id}", h.HandleGetById)
	r.Patch("/{id}", h.HandleUpdate)

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
	if err := h.service.Create(r.Context(), &subject); err != nil {
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
	if err := h.service.Delete(r.Context(), id); err != nil {
		// 404 when subject is not found
		if errors.Is(err, sql.ErrNoRows) {
			sendError(w, http.StatusNotFound, "Subject not found", nil)
			return
		}

		// 500 for db connection or syntax errors
		sendError(w, http.StatusInternalServerError, "Failed to delete subject", nil)
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
	subject, err := h.service.GetById(r.Context(), id)
	if err != nil {
		sendError(w, http.StatusInternalServerError, "Database error", nil)
		return
	}

	// If the subject does not exist
	if subject == nil {
		sendError(w, http.StatusNotFound, "Subject not found", nil)
		return
	}

	sendSuccess(w, http.StatusOK, "", subject)
}

func (h *SubjectHandler) HandleList(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()

	// Error is not needed as default values are in place
	reqPage, _ := strconv.Atoi(queryParams.Get("page"))
	reqLimit, _ := strconv.Atoi(queryParams.Get("limit"))

	// Service returns the actual normalized page and limit used
	subjects, totalItems, page, limit, err := h.service.List(r.Context(), reqPage, reqLimit)
	if err != nil {
		sendError(w, http.StatusInternalServerError, "Failed to fetch subjects", nil)
		return
	}

	totalPages := 0
	if totalItems > 0 {
		totalPages = (totalItems + limit - 1) / limit
	}

	meta := PaginatedMeta{
		Page:       page,
		Limit:      limit,
		Count:      len(subjects),
		TotalItems: totalItems,
		TotalPages: totalPages,
	}

	sendPaginated(w, http.StatusOK, subjects, meta)
}

func (h *SubjectHandler) HandleUpdate(w http.ResponseWriter, r *http.Request) {
	// Extract and convert the id to int
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		sendError(w, http.StatusBadRequest, "Invalid subject ID", nil)
		return
	}

	// Decode de request's body into the pointer-based PATCH DTO
	var input domain.PatchSubjectInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		sendError(w, http.StatusBadRequest, "Invalid JSON payload", nil)
		return
	}

	// Validate field constraints
	if err := h.validate.StructCtx(r.Context(), &input); err != nil {
		if validationErrs, ok := err.(validator.ValidationErrors); ok {
			sendError(w, http.StatusUnprocessableEntity, "Validation failed", formatValidationErrors(validationErrs))
			return
		}
		sendError(w, http.StatusBadRequest, "Validation failed", nil)
		return
	}

	// Perform the update in the db
	updatedSubject, err := h.service.Update(r.Context(), id, input)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			sendError(w, http.StatusNotFound, "Subject not found", nil)
			return
		}
		sendError(w, http.StatusInternalServerError, "Failed to update the class", nil)
		return
	}

	sendSuccess(w, http.StatusOK, "Subject updated successfully", updatedSubject)
}

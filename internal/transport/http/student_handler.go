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

// StudentHandler contains `repo` with a way to communicate with the database and the pointer to the validator instantiated once in `main.go`
type StudentHandler struct {
	service  domain.StudentService
	validate *validator.Validate
}

// NewStudentHandler is a constructor that returns a pointer to a StudentHandler struct, initializing it with the injected repository and validator dependencies
func NewStudentHandler(service domain.StudentService, validate *validator.Validate) *StudentHandler {
	return &StudentHandler{service: service, validate: validate}
}

// Route paths
func (h *StudentHandler) StudentRoutes() chi.Router {
	r := chi.NewRouter()

	r.Get("/", h.HandleList)
	r.Post("/", h.HandleCreate)
	r.Delete("/bulk", h.HandleBulkDelete)
	r.Post("/bulk", h.HandleBulkCreate)
	r.Patch("/bulk", h.HandleBulkUpdate)
	r.Patch("/bulk_class", h.HandleBulkUpdateClass)
	r.Delete("/{id}", h.HandleDelete)
	r.Get("/{id}", h.HandleGetByID)
	r.Patch("/{id}", h.HandleUpdate)

	return r
}

func (h *StudentHandler) HandleBulkCreate(w http.ResponseWriter, r *http.Request) {
	var input domain.BulkCreateStudentInput

	// Decode HTTP body into DTO
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		sendError(w, http.StatusBadRequest, "Invalid JSON payload", nil)
		return
	}

	// Validate incoming payload constraints
	if err := h.validate.StructCtx(r.Context(), &input); err != nil {
		if validationErrs, ok := err.(validator.ValidationErrors); ok {
			sendError(w, http.StatusUnprocessableEntity, "Validation failed", formatValidationErrors(validationErrs))
			return
		}
		sendError(w, http.StatusBadRequest, "Validation failed", nil)
		return
	}

	// Delegate business logic, fetching, normalization, and persistence to service
	createdStudents, total, err := h.service.BulkCreate(r.Context(), input.Students)
	if err != nil {
		sendError(w, http.StatusInternalServerError, "Failed to create students", nil)
		return
	}

	result := BulkResult[domain.Student]{
		Items: createdStudents,
		Meta: BulkMeta{
			Count: len(createdStudents),
			Total: total,
		},
	}

	sendSuccess(w, http.StatusCreated, "Students created successfully", result)
}

func (h *StudentHandler) HandleBulkDelete(w http.ResponseWriter, r *http.Request) {
	var input domain.BulkDeleteStudentInput

	// Decode HTTP body into DTO
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		sendError(w, http.StatusBadRequest, "Invalid JSON payload", nil)
		return
	}

	// Validate incoming payload constraints
	if err := h.validate.StructCtx(r.Context(), &input); err != nil {
		if validationErrs, ok := err.(validator.ValidationErrors); ok {
			sendError(w, http.StatusUnprocessableEntity, "Validation failed", formatValidationErrors(validationErrs))
			return
		}
		sendError(w, http.StatusBadRequest, "Validation failed", nil)
		return
	}

	// Delegate business logic, fetching, and persistence to service
	deletedCount, err := h.service.BulkDelete(r.Context(), input.IDs)
	if err != nil {
		sendError(w, http.StatusInternalServerError, "Failed to delete students in batch", nil)
		return
	}

	meta := BulkMeta{
		Count: int(deletedCount),
		Total: len(input.IDs),
	}

	sendSuccess(w, http.StatusOK, "Students deleted successfully", meta)
}

func (h *StudentHandler) HandleBulkUpdate(w http.ResponseWriter, r *http.Request) {
	var input domain.BulkUpdateStudentInput

	// Decode HTTP body into DTO
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		sendError(w, http.StatusBadRequest, "Invalid JSON payload", nil)
		return
	}

	// Validate incoming payload constraints
	if err := h.validate.StructCtx(r.Context(), &input); err != nil {
		if validationErrs, ok := err.(validator.ValidationErrors); ok {
			sendError(w, http.StatusUnprocessableEntity, "Validation failed", formatValidationErrors(validationErrs))
			return
		}
		sendError(w, http.StatusBadRequest, "Validation failed", nil)
		return
	}

	// Delegate business logic, fetching, normalization, and persistence to service
	updatedStudents, total, err := h.service.BulkUpdate(r.Context(), input.Students)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			sendError(w, http.StatusNotFound, "One or more target students were not found", nil)
			return
		}
		sendError(w, http.StatusInternalServerError, "Failed to update students in batch", nil)
		return
	}

	result := BulkResult[domain.Student]{
		Items: updatedStudents,
		Meta: BulkMeta{
			Count: len(updatedStudents),
			Total: total,
		},
	}

	sendSuccess(w, http.StatusOK, "Students updated successfully", result)
}

func (h *StudentHandler) HandleBulkUpdateClass(w http.ResponseWriter, r *http.Request) {
	var input domain.BulkUpdateClassInput

	// Decode HTTP body into DTO
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		sendError(w, http.StatusBadRequest, "Invalid JSON payload", nil)
		return
	}

	// Validate incoming payload constraints
	if err := h.validate.StructCtx(r.Context(), &input); err != nil {
		if validationErrs, ok := err.(validator.ValidationErrors); ok {
			sendError(w, http.StatusUnprocessableEntity, "Validation failed", formatValidationErrors(validationErrs))
			return
		}
		sendError(w, http.StatusBadRequest, "Validation failed", nil)
		return
	}

	// Delegate business logic, fetching, normalization, and persistence to service
	rowsAffected, err := h.service.BulkUpdateClass(r.Context(), input.StudentIDs, input.ClassID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			sendError(w, http.StatusNotFound, "Target class or students not found", nil)
			return
		}
		sendError(w, http.StatusInternalServerError, "Failed to update students class", nil)
		return
	}

	meta := BulkMeta{
		Count: int(rowsAffected),
		Total: len(input.StudentIDs),
	}

	sendSuccess(w, http.StatusOK, "Class updated successfully", meta)
}

func (h *StudentHandler) HandleCreate(w http.ResponseWriter, r *http.Request) {
	var student domain.Student

	// Decode JSON body directly into domain.Student
	if err := json.NewDecoder(r.Body).Decode(&student); err != nil {
		sendError(w, http.StatusBadRequest, "Invalid JSON payload", nil)
		return
	}

	// Validate struct rules using validator instance
	if err := h.validate.StructCtx(r.Context(), &student); err != nil {
		if validationErrs, ok := err.(validator.ValidationErrors); ok {
			sendError(w, http.StatusUnprocessableEntity, "Validation failed", formatValidationErrors(validationErrs))
			return
		}
		sendError(w, http.StatusBadRequest, "Validation failed", nil)
		return
	}

	// Delegate normalization, timing, and persistence to service layer
	if err := h.service.Create(r.Context(), &student); err != nil {
		sendError(w, http.StatusInternalServerError, "Failed to create student entry", nil)
		return
	}

	// Return 201 Created with full record (ID and timestamps attached in place)
	sendSuccess(w, http.StatusCreated, "Student created successfully", student)
}

func (h *StudentHandler) HandleDelete(w http.ResponseWriter, r *http.Request) {
	// Extract and convert URL param to integer
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		sendError(w, http.StatusBadRequest, "Invalid student ID", nil)
		return
	}

	// Delegate deletion to service layer
	if err := h.service.Delete(r.Context(), id); err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			sendError(w, http.StatusNotFound, "Student not found", nil)
			return
		}

		sendError(w, http.StatusInternalServerError, "Failed to delete student", nil)
		return
	}

	// HTTP 204 No Content for successful deletion (no body needed)
	w.WriteHeader(http.StatusNoContent)
}

func (h *StudentHandler) HandleGetByID(w http.ResponseWriter, r *http.Request) {
	// Extract and convert URL param to integer
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		sendError(w, http.StatusBadRequest, "Invalid student ID", nil)
		return
	}

	// Delegate fetching to the service layer
	student, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			sendError(w, http.StatusNotFound, "Student not found", nil)
			return
		}

		sendError(w, http.StatusInternalServerError, "Failed to fetch student", nil)
		return
	}
	sendSuccess(w, http.StatusOK, "Student retrieved successfully", student)
}

func (h *StudentHandler) HandleList(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()

	// Parse page (default: 1)
	page := 1
	if pageStr := queryParams.Get("page"); pageStr != "" {
		if parsedPage, err := strconv.Atoi(pageStr); err == nil && parsedPage > 0 {
			page = parsedPage
		}
	}

	// Parse limit (default: 10)
	limit := 10
	if limitStr := queryParams.Get("limit"); limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	// Derive SQL offset
	offset := (page - 1) * limit

	// Delegate paginated fetch and total count calculation to service layer
	students, totalItems, err := h.service.List(r.Context(), limit, offset)
	if err != nil {
		fmt.Println(err)
		sendError(w, http.StatusInternalServerError, "Failed to fetch student list", nil)
		return
	}

	// Calculate total pages using integer arithmetic
	totalPages := 0
	if totalItems > 0 {
		totalPages = (totalItems + limit - 1) / limit
	}

	// Construct metadata and issue generic paginated response
	meta := PaginatedMeta{
		Page:       page,
		Limit:      limit,
		Count:      len(students),
		TotalItems: totalItems,
		TotalPages: totalPages,
	}

	sendPaginated(w, http.StatusOK, students, meta)
}

func (h *StudentHandler) HandleUpdate(w http.ResponseWriter, r *http.Request) {
	// Extract and convert student ID from URL path parameter
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		sendError(w, http.StatusBadRequest, "Invalid student ID", nil)
		return
	}

	// Decode HTTP body into PatchStudentInput DTO
	var input domain.PatchStudentInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		sendError(w, http.StatusBadRequest, "Invalid JSON payload", nil)
		return
	}

	// Validate patch payload constraints
	if err := h.validate.StructCtx(r.Context(), &input); err != nil {
		if validationErrs, ok := err.(validator.ValidationErrors); ok {
			sendError(w, http.StatusUnprocessableEntity, "Validation failed", formatValidationErrors(validationErrs))
			return
		}
		sendError(w, http.StatusBadRequest, "Validation failed", nil)
		return
	}

	// Delegate fetching, patch merging, normalization, and updates to service layer
	updatedStudent, err := h.service.Update(r.Context(), id, input)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			sendError(w, http.StatusNotFound, "Student not found", nil)
			return
		}
		sendError(w, http.StatusInternalServerError, "Failed to update student", nil)
		return
	}

	sendSuccess(w, http.StatusOK, "Student updated successfully", updatedStudent)
}

package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/ednanf/school-api/internal/domain"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

type userHandler struct {
	service  domain.UserService
	validate *validator.Validate
}

func NewUserHandler(service domain.UserService, validate *validator.Validate) *userHandler {
	return &userHandler{service: service, validate: validate}
}

func (h *userHandler) UserRoutes() chi.Router {
	r := chi.NewRouter()

	r.Get("/", h.HandleList)
	r.Post("/", h.HandleCreate)
	r.Get("/{id}", h.HandleGetByID)
	r.Patch("/{id}", h.HandleUpdate)
	r.Delete("/{id}", h.HandleDelete)

	return r
}

func (h *userHandler) HandleCreate(w http.ResponseWriter, r *http.Request) {
	var input domain.CreateUserInput

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		sendError(w, http.StatusBadRequest, "Invalid JSON payload", nil)
		return
	}

	if err := h.validate.StructCtx(r.Context(), &input); err != nil {
		if validationErrs, ok := err.(validator.ValidationErrors); ok {
			sendError(w, http.StatusUnprocessableEntity, "Validation failed", formatValidationErrors(validationErrs))
			return
		}
		sendError(w, http.StatusBadRequest, "Validation failed", nil)
		return
	}

	// Create an user struct to be passed to Create service. The password, however, will come directly from the decoded JSON (var input)
	user := domain.User{
		StaffID:  input.StaffID,
		Username: input.Username,
		Email:    input.Email,
		Role:     input.Role,
		IsActive: input.IsActive,
	}

	if err := h.service.Create(r.Context(), &user, input.Password); err != nil {
		sendError(w, http.StatusInternalServerError, "Failed to create user", nil)
		return
	}

	// Response payload is a separate struct from User, to avoid leaking credentials
	response := domain.UserResponse{
		ID:       user.ID,
		StaffID:  user.StaffID,
		Username: user.Username,
		Email:    user.Email,
		Role:     user.Role,
		IsActive: user.IsActive,
	}

	sendSuccess(w, http.StatusCreated, "User created successfully", response)
}

func (h *userHandler) HandleDelete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil || id <= 0 {
		sendError(w, http.StatusBadRequest, "Invalid user ID", nil)
		return
	}

	if err := h.service.Delete(r.Context(), id); err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			sendError(w, http.StatusNotFound, "User not found", nil)
			return
		}
		sendError(w, http.StatusInternalServerError, "Failed to delete user", nil)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *userHandler) HandleGetByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil || id <= 0 {
		sendError(w, http.StatusBadRequest, "Invalid user ID", nil)
		return
	}

	user, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			sendError(w, http.StatusNotFound, "User not found", nil)
			return
		}
		sendError(w, http.StatusInternalServerError, "Failed to fetch user", nil)
		return
	}

	sendSuccess(w, http.StatusOK, "User retrieved successfully", user)
}

func (h *userHandler) HandleList(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()

	page := 1
	if pageStr := queryParams.Get("page"); pageStr != "" {
		if parsedPage, err := strconv.Atoi(pageStr); err == nil && parsedPage > 0 {
			page = parsedPage
		}
	}

	limit := 10
	if limitStr := queryParams.Get("limit"); limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	users, totalItems, pageNum, limitNum, err := h.service.List(r.Context(), page, limit)
	if err != nil {
		sendError(w, http.StatusInternalServerError, "Failed to fetch user list", nil)
		return
	}

	totalPages := 0
	if totalItems > 0 {
		totalPages = (totalItems + limitNum - 1) / limitNum
	}

	meta := PaginatedMeta{
		Page:       pageNum,
		Limit:      limitNum,
		Count:      len(users),
		TotalItems: totalItems,
		TotalPages: totalPages,
	}

	sendPaginated(w, http.StatusOK, users, meta)
}

func (h *userHandler) HandleUpdate(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil || id <= 0 {
		sendError(w, http.StatusBadRequest, "Invalid user ID", nil)
		return
	}

	var input domain.PatchUserInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		sendError(w, http.StatusBadRequest, "Invalid JSON payload", nil)
		return
	}

	if !input.HasUpdates() {
		sendError(w, http.StatusBadRequest, "At least one field must be provided for update", nil)
		return
	}

	if err := h.validate.StructCtx(r.Context(), &input); err != nil {
		if validationErrs, ok := err.(validator.ValidationErrors); ok {
			sendError(w, http.StatusUnprocessableEntity, "Validation failed", formatValidationErrors(validationErrs))
			return
		}
		sendError(w, http.StatusBadRequest, "Validation failed", nil)
		return
	}

	updatedUser, err := h.service.Update(r.Context(), id, input)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			sendError(w, http.StatusNotFound, "User not found", nil)
			return
		}
		sendError(w, http.StatusInternalServerError, "Failed to update user", nil)
		return
	}

	// Response uses UserReponse instead of User to avoid password leakage
	response := domain.UserResponse{
		ID:       updatedUser.ID,
		StaffID:  updatedUser.StaffID,
		Username: updatedUser.Username,
		Email:    updatedUser.Email,
		Role:     updatedUser.Role,
		IsActive: updatedUser.IsActive,
	}

	sendSuccess(w, http.StatusOK, "User updated successfully", response)
}

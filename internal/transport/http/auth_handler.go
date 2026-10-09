package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/ednanf/school-api/internal/domain"
	logger "github.com/ednanf/school-api/internal/pkg/loggers"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

type authHandler struct {
	service  domain.AuthService
	validate *validator.Validate
}

func NewAuthHandler(service domain.AuthService, validate *validator.Validate) *authHandler {
	return &authHandler{
		service:  service,
		validate: validate,
	}
}

func (h *authHandler) Routes(tokenService domain.TokenService) chi.Router {
	r := chi.NewRouter()

	// Public auth routes
	r.Post("/login", h.Login)
	r.Post("/forgot-password", h.ForgotPassword)
	r.Post("/reset-password", h.ResetPassword)

	// Protected auth routes (requires valid JWT)
	r.Group(func(r chi.Router) {
		r.Use(AuthMiddleware(tokenService))
		r.Post("/change-password", h.ChangePassword)
	})

	return r
}

func (h *authHandler) Login(w http.ResponseWriter, r *http.Request) {
	var input domain.LoginInput
	fmt.Println(input)
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		logger.DebugLogger(err)
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

	user, token, err := h.service.Login(r.Context(), input)
	if err != nil {
		if errors.Is(err, domain.ErrInvalidCredentials) {
			sendError(w, http.StatusUnauthorized, "Invalid email or password", nil)
			return
		}
		if errors.Is(err, domain.ErrUserInactive) {
			sendError(w, http.StatusForbidden, "User account is inactive", nil)
			return
		}
		sendError(w, http.StatusInternalServerError, "Failed to authenticate", nil)
		return
	}

	// Response structure returning user details and the bearer JWT
	response := map[string]any{
		"token": token,
		"user":  user.ToResponse(),
	}

	sendSuccess(w, http.StatusOK, "Login successful", response)
}

func (h *authHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	// Extract claims injected into r.Context() by AuthMiddleware
	claims, ok := r.Context().Value(userClaimsKey).(*domain.CustomClaims)
	if !ok || claims == nil {
		sendError(w, http.StatusUnauthorized, "Unauthorized access", nil)
		return
	}

	// Decode the incoming JSON request body
	var input domain.ChangePasswordInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		sendError(w, http.StatusBadRequest, "Invalid JSON payload", nil)
		return
	}

	// Validate struct tags
	if err := h.validate.StructCtx(r.Context(), &input); err != nil {
		if validationErrs, ok := err.(validator.ValidationErrors); ok {
			sendError(w, http.StatusUnprocessableEntity, "Validation failed", formatValidationErrors(validationErrs))
			return
		}
		sendError(w, http.StatusBadRequest, "Validation failed", nil)
		return
	}

	// Call the service layer with userID from token claims
	if err := h.service.ChangePassword(r.Context(), claims.UserID, input); err != nil {
		if errors.Is(err, domain.ErrInvalidCredentials) {
			sendError(w, http.StatusUnauthorized, "Current password is incorrect", nil)
			return
		}
		if errors.Is(err, domain.ErrNotFound) {
			sendError(w, http.StatusNotFound, "User not found", nil)
			return
		}
		sendError(w, http.StatusInternalServerError, "Failed to change password", nil)
		return
	}

	// Return success
	sendSuccess(w, http.StatusOK, "Password changed successfully", nil)
}

func (h *authHandler) ForgotPassword(w http.ResponseWriter, r *http.Request) {
	var input domain.ForgotPasswordInput
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

	// Service generates token; returning success unconditionally avoids email enumeration
	_, err := h.service.ForgotPassword(r.Context(), input)
	if err != nil {
		sendError(w, http.StatusInternalServerError, "Failed to process forgot password request", nil)
		return
	}

	sendSuccess(w, http.StatusOK, "If an account exists with that email, a password reset link has been generated", nil)
}

func (h *authHandler) ResetPassword(w http.ResponseWriter, r *http.Request) {
	var input domain.ResetPasswordInput
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

	if err := h.service.ResetPassword(r.Context(), input); err != nil {
		if errors.Is(err, domain.ErrInvalidToken) {
			sendError(w, http.StatusBadRequest, "Invalid or expired reset token", nil)
			return
		}
		sendError(w, http.StatusInternalServerError, "Failed to reset password", nil)
		return
	}

	sendSuccess(w, http.StatusOK, "Password reset successfully", nil)
}

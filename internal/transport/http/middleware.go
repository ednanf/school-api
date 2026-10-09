package http

import (
	"context"
	"net/http"
	"slices"
	"strings"

	"github.com/ednanf/school-api/internal/domain"
)

// Declare the context key
type contextKey string

// Context key that will carry user data via r.Context()
const userClaimsKey contextKey = "userClaims"

// AuthMiddleware verifies the JWT token from the Authorization header
func AuthMiddleware(tokenService domain.TokenService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Extract "Bearer <token>" from Authorization header
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				sendError(w, http.StatusUnauthorized, "Missing authorization header", nil)
				return
			}

			// Split Bearer and Token value apart
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				sendError(w, http.StatusUnauthorized, "Invalid authorization header format", nil)
				return
			}

			tokenStr := parts[1]

			// Validate JWT
			claims, err := tokenService.ValidateToken(tokenStr)
			if err != nil {
				sendError(w, http.StatusUnauthorized, "Invalid or expired token", nil)
				return
			}

			// Attach claims to r.Context()
			ctx := context.WithValue(r.Context(), userClaimsKey, claims)

			// Pass request along with the updated context
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// RequireRole checks if the authenticated user has one of the allowed roles
func RequireRole(roles ...domain.UserRole) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Extract claims put into context by AuthMiddleware
			claims, ok := r.Context().Value(userClaimsKey).(*domain.CustomClaims)
			if !ok || claims == nil {
				sendError(w, http.StatusUnauthorized, "Unauthorized", nil)
				return
			}

			// Check if user's role matches any allowed role
			hasRole := slices.Contains(roles, claims.Role)

			if !hasRole {
				sendError(w, http.StatusForbidden, "Forbidden: insufficient permissions", nil)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

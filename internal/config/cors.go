package config

import (
	"net/http"

	"github.com/go-chi/cors"
)

// SetupCORS returns the CORS middleware handler for the router.
func SetupCORS() func(http.Handler) http.Handler {
	return cors.Handler(cors.Options{
		AllowedOrigins: []string{
			"http://localhost:5173",
			"https://your-production-domain.com", // TODO: Add prod url here
		},
		AllowedMethods: []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodPatch,
			http.MethodDelete,
			http.MethodOptions,
		},
		AllowedHeaders: []string{
			"Accept",
			"Authorization",
			"Content-Type",
		},
		AllowCredentials: false,
		MaxAge:           300,
	})
}

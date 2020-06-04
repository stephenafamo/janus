package middlewares

import (
	"net/http"

	"github.com/go-chi/cors"
)

// CORSMiddleware allows CORS on specific domains set by the server
func CORSMiddleware(allowedDomains []string) func(h http.Handler) http.Handler {
	cors := cors.New(cors.Options{
		AllowedOrigins: allowedDomains,
		AllowedMethods: []string{
			"GET",
			"POST",
			"PATCH",
			"DELETE",
			"OPTIONS",
		},
		AllowedHeaders: []string{
			"Accept",
			"Accept-Encoding",
			"Authorization",
			"Content-Type",
			"Content-Length",
			"X-CSRF-Token",
		},
		AllowCredentials: true,
		MaxAge:           300,
	})
	return cors.Handler
}

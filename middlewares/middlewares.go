package middlewares

import (
	"log"
	"net/http"

	"github.com/go-chi/cors"
	"github.com/justinas/nosurf"
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

// CSRF is a middleware that adds csrf tokens to the request and validates them
func CSRF(h http.Handler) http.Handler {
	surfing := nosurf.New(h)
	surfing.SetFailureHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("Failed to validate XSRF Token:", nosurf.Reason(r))
		w.WriteHeader(http.StatusBadRequest)
	}))

	return surfing
}

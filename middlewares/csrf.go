package middlewares

import (
	"log"
	"net/http"

	"github.com/justinas/nosurf"
)

// CSRF is a csrf middleware and token interface
type CSRF interface {
	Middleware(h http.Handler) http.Handler
	Token(*http.Request) string
}

// Nosurf satisfies the CSRF interface
type Nosurf struct{}

// Middleware gets the csrf middleware
func (Nosurf) Middleware(h http.Handler) http.Handler {
	surfing := nosurf.New(h)
	surfing.SetFailureHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("Failed to validate XSRF Token:", nosurf.Reason(r))
		w.WriteHeader(http.StatusBadRequest)
	}))

	return surfing
}

// Token gets the csrf token for that request
func (Nosurf) Token(r *http.Request) string {
	return nosurf.Token(r)
}

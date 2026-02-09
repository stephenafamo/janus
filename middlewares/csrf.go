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
type Nosurf struct {
	ExemptFunc func(r *http.Request) bool
}

// Middleware gets the csrf middleware
func (n Nosurf) Middleware(h http.Handler) http.Handler {
	surfing := nosurf.New(h)
	surfing.SetFailureHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("Failed to validate XSRF Token:", nosurf.Reason(r))
		w.WriteHeader(http.StatusBadRequest)
	}))

	// necessary so we don't get duplicate cookies which makes the validation fail in some cases
	surfing.SetBaseCookie(http.Cookie{Path: "/", MaxAge: nosurf.MaxAge})

	if n.ExemptFunc != nil {
		surfing.ExemptFunc(n.ExemptFunc)
	}

	return surfing
}

// Token gets the csrf token for that request
func (Nosurf) Token(r *http.Request) string {
	return nosurf.Token(r)
}

package auth

import (
	"net/http"
)

type ctxKey string

// CtxUserID is the context key for the user ID
var CtxUserID ctxKey = "userID"

// Authenticator is the authenticator
type Authenticator interface {

	// Router will be mounted on the auth path.
	// So the authenticator can handle it's own routing
	Router() http.Handler

	// Middlewares to be applied to all routes
	// Also adds these middlewares to the authentication routes since those
	// will not be protected
	DefaultMiddlewares() []func(http.Handler) http.Handler
	// Middlewares to be applied to routes that need authentication
	// One of these middlewares should load the authenticated user ID into
	// the request context with key `ctxUserID`
	ProtectMidelewares() []func(http.Handler) http.Handler

	// Removes all cookies and session for the current user
	// Used to manually log the user out
	Flush(http.ResponseWriter) error

	// Used to redirect something to the auth pages with a message
	Redirect(rw http.ResponseWriter, req *http.Request, redir RedirectOptions) error
}

// RedirectOptions are the values needed for auth packages to redirect
type RedirectOptions interface {
	// Success & Failure are used to set Flash messages / JSON messages
	// if set. They should be mutually exclusive.
	Success() string
	Failure() string
	Path() string
}

// User represents a user
type User interface {
}

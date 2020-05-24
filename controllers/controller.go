package controllers

import (
	"log"
	"net/http"
	"os"

	"github.com/NYTimes/gziphandler"
	"github.com/go-chi/chi/middleware"
	"github.com/stephenafamo/janus/auth"
	"github.com/stephenafamo/janus/middlewares"
	"github.com/stephenafamo/janus/store"
	"github.com/stephenafamo/janus/views"
)

type logger interface {
	Printf(format string, v ...interface{})
}

// Logger to print lgos
var Logger logger = log.New(os.Stderr, "", log.LstdFlags)

// CSRFMiddleware is the middleware that will be returned in c.RecommnededMiddlewares
// It is exported to allow external change to the CSRF middleware
var CSRFMiddleware mid

type mid = func(http.Handler) http.Handler

// Controller contains all our routes and handlers and can return a handler
// for our http.Server
type Controller struct {
	Domains   []string
	Templates views.TemplateExecutor
	Store     store.Store
	Assets    http.FileSystem
	Auth      auth.Authenticator
}

// WriteError is a helper to return an error
func WriteError(w http.ResponseWriter, r *http.Request, err error, code int) {
	Logger.Printf("HTTP ERROR: %v", err)
	http.Error(w, http.StatusText(code), code)
}

// SetDomains for our handler
func (c Controller) SetDomains(d []string) {
	c.Domains = d
}

// SetTemplates for our handler
func (c Controller) SetTemplates(t views.TemplateExecutor) {
	c.Templates = t
}

// SetStore for our handler
func (c Controller) SetStore(s store.Store) {
	c.Store = s
}

// SetAssets for our handler
func (c Controller) SetAssets(s http.FileSystem) {
	c.Assets = s
}

// SetAuth for our handler
func (c Controller) SetAuth(a auth.Authenticator) {
	c.Auth = a
}

// RecommendedMiddlewares for our handler
func (c Controller) RecommendedMiddlewares() []mid {
	UseCSRF := CSRFMiddleware
	if UseCSRF == nil {
		UseCSRF = middlewares.CSRF
	}
	mids := []mid{
		middleware.StripSlashes,
		middleware.RequestID,
		middleware.RealIP,
		middleware.Logger,
		middleware.Recoverer,
		gziphandler.GzipHandler,
		middlewares.CORSMiddleware(c.Domains),
		UseCSRF,
	}

	if c.Auth != nil {
		// Add the default auth middlewares
		mids = append(mids, c.Auth.DefaultMiddlewares()...)
	}

	return mids
}

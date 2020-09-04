package janus

import (
	"log"
	"net/http"
	"os"

	"github.com/NYTimes/gziphandler"
	"github.com/go-chi/chi/middleware"
	"github.com/spf13/afero"
	"github.com/stephenafamo/janus/auth"
	"github.com/stephenafamo/janus/middlewares"
	"github.com/stephenafamo/janus/views/executor"
)

type mid = func(http.Handler) http.Handler

type logger interface {
	Printf(format string, v ...interface{})
}

// Logger to print lgos
var Logger logger = log.New(os.Stderr, "", log.LstdFlags)

// CSRFMiddleware is the middleware that will be returned in c.RecommnededMiddlewares
// It is exported to allow external change to the CSRF middleware
var CSRFMiddleware middlewares.CSRF = middlewares.Nosurf{}

// Broker contains all our routes and handlers and can return a handler
// for our http.Server
type Broker struct {
	Domains   []string
	Templates executor.Executor
	Store     afero.Fs
	Assets    http.FileSystem
	Auth      auth.Authenticator

	csrfMiddleware middlewares.CSRF
}

// WriteError is a helper to return an error
func WriteError(w http.ResponseWriter, r *http.Request, err error, code int) {
	Logger.Printf("HTTP ERROR: %v", err)
	http.Error(w, http.StatusText(code), code)
}

// SetDomains for our handler
func (b *Broker) SetDomains(d []string) {
	b.Domains = d
}

// SetTemplates for our handler
func (b *Broker) SetTemplates(t executor.Executor) {
	b.Templates = t
}

// SetStore for our handler
func (b *Broker) SetStore(s afero.Fs) {
	b.Store = s
}

// SetAssets for our handler
func (b *Broker) SetAssets(s http.FileSystem) {
	b.Assets = s
}

// SetAuth for our handler
func (b *Broker) SetAuth(a auth.Authenticator) {
	b.Auth = a
}

// GetCSRFToken to use for subsequet requests
func (b *Broker) GetCSRFToken(r *http.Request) string {
	m := b.csrfMiddleware
	if m == nil {
		m = CSRFMiddleware
	}
	return m.Token(r)
}

// RecommendedMiddlewares for our handler
func (b Broker) RecommendedMiddlewares() []mid {

	mids := []mid{
		middleware.StripSlashes,
		middleware.RequestID,
		middleware.RealIP,
		middleware.Logger,
		middleware.Recoverer,
		gziphandler.GzipHandler,
		middlewares.CORSMiddleware(b.Domains),
		// CSRFMiddleware.Middleware,
	}

	if b.Auth != nil {
		// Add the default auth middlewares
		mids = append(mids, b.Auth.DefaultMiddlewares()...)
	}

	return mids
}

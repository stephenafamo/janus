package authboss

import (
	"context"
	"errors"
	"log"
	"net/http"

	"github.com/justinas/nosurf"
	"github.com/stephenafamo/janus/auth"
	"github.com/volatiletech/authboss/v3"
)

// GetCsrfToken gets the csrf token that authboss adds to HTMLData
var GetCsrfToken = func(r *http.Request) string {
	return nosurf.Token(r)
}

// Authboss satisfies the Auth interface
// Based on the excellent package github.com/volatiletech/authboss
type Authboss struct {
	*authboss.Authboss

	// For module middlewares
	ExtraDefaultMiddlewares []func(http.Handler) http.Handler
	ExtraProtectMiddlewares []func(http.Handler) http.Handler
}

// Flush satisfies the Auth interface
func (a Authboss) Flush(rw http.ResponseWriter) error {
	authboss.DelAllSession(rw, []string{
		authboss.FlashSuccessKey,
		authboss.FlashErrorKey,
	})
	authboss.DelKnownCookie(rw)
	return nil
}

// Router satisfies the Auth interfaces
func (a Authboss) Router() http.Handler {
	return a.Config.Core.Router
}

// DefaultMiddlewares satisfies the Auth interfaces
func (a Authboss) DefaultMiddlewares() []func(http.Handler) http.Handler {
	mids := []func(http.Handler) http.Handler{
		a.LoadClientStateMiddleware,
		a.AddUserIDToContext,
		a.DataInjector,
		a.RedirectIfLoggedIn,
	}
	if a.ExtraDefaultMiddlewares != nil {
		mids = append(mids, a.ExtraDefaultMiddlewares...)
	}

	return mids
}

// ProtectMidelewares satisfies the Auth interfaces
func (a Authboss) ProtectMidelewares() []func(http.Handler) http.Handler {
	mids := []func(http.Handler) http.Handler{
		authboss.Middleware2(
			a.Authboss,
			authboss.RequireNone,
			authboss.RespondRedirect,
		),
	}
	if a.ExtraProtectMiddlewares != nil {
		mids = append(mids, a.ExtraProtectMiddlewares...)
	}

	return mids
}

// RedirectIfLoggedIn redirects logged in users if visiting the login or register page
func (a Authboss) RedirectIfLoggedIn(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		pid, err := a.CurrentUserID(r)
		if err != nil {
			log.Printf("Error in RedirectIfLoggedIn middleware: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}

		mountPath := a.Config.Paths.Mount
		switch r.URL.Path {
		case mountPath + "/login", mountPath + "/register":
			if pid != "" {
				log.Printf("Redirecting authorized user %q", pid)

				redirTo := r.FormValue("redir")
				if redirTo == "" {
					redirTo = a.Paths.AuthLoginOK
				}

				ro := authboss.RedirectOptions{
					Code:             http.StatusTemporaryRedirect,
					RedirectPath:     redirTo,
					FollowRedirParam: true,
				}
				if err := a.Core.Redirector.Redirect(w, r, ro); err != nil {
					log.Printf("Error in RedirectIfLoggedIn middleware: %v", err)
					http.Error(w, "Internal server error", http.StatusInternalServerError)
				}
			}
		}

		h.ServeHTTP(w, r)
	})
}

// DataInjector is a middleware that adds some auth related values to context
func (a Authboss) DataInjector(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data := authboss.HTMLData{
			"baseUrl":       a.Config.Paths.RootURL,
			"flash_success": authboss.FlashSuccess(w, r),
			"flash_error":   authboss.FlashError(w, r),
			"csrf_token":    GetCsrfToken(r),
		}
		r = r.WithContext(context.WithValue(r.Context(), authboss.CTXKeyData, data))
		handler.ServeHTTP(w, r)
	})
}

// AddUserIDToContext is a middleware that adds some auth related values to context
func (a Authboss) AddUserIDToContext(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, err := a.CurrentUser(r)
		if err != nil && !errors.Is(err, authboss.ErrUserNotFound) {
			log.Printf("Error in AddUserIDToContext middleware: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}

		if !errors.Is(err, authboss.ErrUserNotFound) {
			r = r.WithContext(context.WithValue(r.Context(), auth.CtxUserID, user.GetPID()))
		}

		handler.ServeHTTP(w, r)
	})
}

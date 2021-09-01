package monitor

import (
	"context"
	"net/http"
	"time"
)

type ctxKey string

// CtxUserID is the context key for the scope
var CtxScope ctxKey = "scope"

type Monitor interface {
	// Any implementation must set the scope to the request context in the middleware
	Middleware(http.Handler) http.Handler
	StartSpan(context.Context, string) Span
	CaptureMessage(msg string, tags map[string]string)
	CaptureException(err error, tags map[string]string)
	Recover(ctx context.Context, cause interface{}) error
	Flush(timeout time.Duration)
}

type Scope interface {
	SetTag(key, value string)
	SetUser(id, username, email string)
}

type Span interface {
	SetTag(string, string)
	Finish()
}

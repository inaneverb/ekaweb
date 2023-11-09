package ekaweb_private

import (
	"net/http"
)

// _UkvsManagerMiddleware is the wrapper for UkvsManager, that works flawlessly
// as a middleware, allowing you to initialize UKVS and prepare it
// for being used in next HTTP handlers.
// It also de-initializes and cleans up UKVS after the job is done.
type _UkvsManagerMiddleware struct {
	manager *UkvsManager
}

func (m _UkvsManagerMiddleware) CheckErrorBefore() bool {
	return false
}

func (m _UkvsManagerMiddleware) Callback(next Handler) Handler {
	return HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var ctx = m.manager.InjectUkvs(r.Context())
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r.WithContext(ctx))
		m.manager.ReturnUkvs(ctx)
	})
}

// NewUkvsManagerMiddleware returns a new _UkvsManagerMiddleware
// as a Middleware. That one will initialize UKVS using given UkvsManager
// for next HTTP handler and then performs cleanup after job is done.
func NewUkvsManagerMiddleware(manager *UkvsManager) Middleware {
	return _UkvsManagerMiddleware{manager}
}

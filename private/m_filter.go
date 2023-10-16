package ekaweb_private

import (
	"slices"
)

// FilterNilMiddlewares filters given 'middlewares', putting only not-nil
// objects to the returned list.
// Deprecated: Not needed anymore. BuildHandlerOut() contains all checks.
func FilterNilMiddlewares(middlewares []Middleware) []Middleware {
	return slices.DeleteFunc(middlewares, func(m Middleware) bool {
		return m == nil
	})
}

// FilterNilHandlers filters given 'handlers', putting only not-nil
// objects to the returned list.
// Deprecated: Not necessary anymore. BuildHandlerOut contains all checks.
func FilterNilHandlers(handlers []Handler) []Handler {
	return slices.DeleteFunc(handlers, func(h Handler) bool {
		return h == nil
	})
}

package ekaweb_private

import (
	"net/http"
	"slices"
)

// MergeHandlers just returns one handler that will call all given 'handlers'
// in a row, one-by-one. If empty array is passed, an EmptyHandler is returned.
//
// WARNING! No nil checks before. Filter handlers first if you have to.
// See more: FilterNilHandlers().
func MergeHandlers(handlers []Handler) Handler {

	switch len(handlers) {
	case 0:
		return NewEmptyHandler()

	case 1:
		return handlers[0]

	case 2:
		return mergeTwoHandlers(handlers[0], handlers[1])
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for _, handler := range handlers {
			handler.ServeHTTP(w, r)
		}
	})
}

// MergeMiddlewares just returns one handler that wraps given 'handler'
// by the all provided 'middlewares', one-by-one.
// It means, that the whole chain of inner Middleware.Callback() calls,
// that generates an output Handler is under construction exactly here.
//
// If there's no middlewares, provided handler is returned.
//
// WARNING!
// No nil checks for anything. Filter middlewares/handler first if you have to.
// See more: FilterNilMiddlewares(), FilterNilHandlers().
func MergeMiddlewares(middlewares []Middleware, handler Handler) Handler {

	slices.Reverse(middlewares)
	for _, middleware := range middlewares {
		handler = middleware.Callback(handler)
	}

	return handler
}

////////////////////////////////////////////////////////////////////////////////
///// PRIVATE METHODS //////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////

// mergeTwoHandlers is a special case of MergeHandlers business logic.
func mergeTwoHandlers(handler1, handler2 Handler) Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handler1.ServeHTTP(w, r)
		handler2.ServeHTTP(w, r)
	})
}

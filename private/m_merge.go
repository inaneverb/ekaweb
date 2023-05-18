package ekaweb_private

import (
	"net/http"
)

// MergeHandlers just returns one handler that will call all provided
// in a row, one-by-one.
// If there's no provided handlers (empty array), an EmptyHandler is returned.
//
// WARNING!
// There's no nil check for provided handlers.
// It's your responsibility to call FilterNilHandlers firstly.
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
		for i, n := 0, len(handlers); i < n; i++ {
			handlers[i].ServeHTTP(w, r)
		}
	})
}

// mergeTwoHandlers is a special case of MergeHandlers business logic.
func mergeTwoHandlers(handler1, handler2 Handler) Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handler1.ServeHTTP(w, r)
		handler2.ServeHTTP(w, r)
	})
}

// MergeMiddlewares just returns one handler that will call:
//   - Each provided middleware, one-by-one, passing next middleware
//     to the current middleware function.
//   - The last handler at the end.
//
// If there's no middlewares, a provided handler just returns.
//
// WARNING!
// There's no nil check for provided middlewares & handler.
// It's your responsibility to call FilterNilMiddlewares and/or
// FilterNilHandlers firstly.
func MergeMiddlewares(
	middlewares []Middleware, handler Handler) Handler {

	for i := len(middlewares) - 1; i >= 0; i-- {
		handler = middlewares[i].Callback(handler)
	}

	return handler
}

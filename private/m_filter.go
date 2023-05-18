package ekaweb_private

import (
	"github.com/inaneverb/ekacore/ekaunsafe/v4"
)

// Deprecated: Not necessary anymore. BuildHandlerOut contains all checks.
func FilterNilMiddlewares(middlewares []Middleware) []Middleware {

	if len(middlewares) == 0 {
		return middlewares
	}

	// Avoid pre-copy. Use copy-on-demand. Make a copy only when you're sure
	// there are presented nil middlewares we need filter.

	middlewaresOld, middlewares := middlewares, nil
	for i, n := 0, len(middlewaresOld); i < n; i++ {
		if middlewaresOld[i] == nil {
			if middlewares == nil {
				middlewares = make([]Middleware, 0, len(middlewaresOld)-1)
				for j := 0; j < i; j++ {
					middlewares = append(middlewares, middlewaresOld[j])
				}
			}
			middlewares = append(middlewares, middlewaresOld[i])
		}
	}

	if middlewares == nil {
		return middlewaresOld // There was no copy at all
	} else {
		return middlewares // There were empty middlewares and copy has been made
	}
}

// Deprecated: Not necessary anymore. BuildHandlerOut contains all checks.
func FilterNilHandlers(handlers []Handler) []Handler {

	if len(handlers) == 0 {
		return handlers
	}

	// Avoid pre-copy. Use copy-on-demand. Make a copy only when you're sure
	// there are presented nil middlewares we need filter.

	handlersOld, handlers := handlers, nil
	for i, n := 0, len(handlersOld); i < n; i++ {
		if ekaunsafe.UnpackInterface(handlersOld[i]).Word == nil {
			if handlers == nil {
				handlers = make([]Handler, 0, len(handlersOld)-1)
				for j := 0; j < i; j++ {
					handlers = append(handlers, handlersOld[j])
				}
			}
			handlers = append(handlers, handlersOld[i])
		}
	}

	if handlers == nil {
		return handlersOld // There was no copy at all
	} else {
		return handlers // There were empty handlers and copy has been made
	}
}

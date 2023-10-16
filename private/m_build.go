package ekaweb_private

// BuildHandlerOut prepares outgoing set of middlewares + handler based
// on provided 'components'.
// It also "attaches" 'checkError' before each middleware, if it's allowed
// ('checkError' is not nil, middleware has no "skip error check" behaviour).
//
// You can skip getting final handler, passing true as 3rd arg. In this case,
// final handler will be transformed to middleware and added to outgoing set.
func BuildHandlerOut(
	components []any, checkError Middleware,
	returnMiddlewaresOnly bool) ([]Middleware, Handler) {

	var enableCheckErrorGlobal = checkError != nil

	// Apply components. The rules (about error check skipping) are:
	//
	// 1. First item in whole HTTP call sequence (middlewares + handlers)
	//    IS ALSO may be protected. It could be a call for HTTP group
	//    and there could be an error before.
	//
	// 2. "Protecting" (middleware or handler) means that this item
	//    will be called only if there was no error before.
	//
	// 3. For middlewares we can just add 'checkError' middleware
	//    to the middlewares call stack. Easy.
	//
	// 4. For handlers, we're wrapping handler by the CheckError() call,
	//    getting a new handler.

	var middlewaresInRow []Middleware
	var handlersInRow []Handler
	var handlersOut []Handler

	var addMiddleware = func(middleware Middleware, performCheckError bool) {
		switch {
		case len(handlersInRow) > 0 && len(middlewaresInRow) > 0:
			var tempHandler = MergeHandlers(handlersInRow)
			tempHandler = MergeMiddlewares(middlewaresInRow, tempHandler)
			handlersOut = append(handlersOut, tempHandler)
			handlersInRow = nil
			middlewaresInRow = nil

		case len(handlersInRow) > 0:
			handlersOut = append(handlersOut, handlersInRow...)
			handlersInRow = nil
		}

		var asMiddlewareExtended, ok = middleware.(MiddlewareExtended)

		if performCheckError && (!ok || asMiddlewareExtended.CheckErrorBefore()) {
			middlewaresInRow = append(middlewaresInRow, checkError)
		}
		middlewaresInRow = append(middlewaresInRow, middleware)
	}

	var addHandler = func(handler Handler, performCheckError bool) {
		var asHandlerExtended, ok = handler.(HandlerExtended)

		if performCheckError && (!ok || asHandlerExtended.CheckErrorBefore()) {
			handler = checkError.Callback(handler)
		}
		handlersInRow = append(handlersInRow, handler)
	}

	for _, component := range components {
		var asMiddleware = AsMiddleware(component)
		var asHandler = AsHandler(component)

		switch {
		case asMiddleware != nil:
			addMiddleware(asMiddleware, enableCheckErrorGlobal)

		case asHandler != nil:
			addHandler(asHandler, enableCheckErrorGlobal)
		}
	}

	var handlerOut Handler = nil
	if len(handlersOut) != 0 || len(handlersInRow) != 0 {
		handlerOut = MergeHandlers(append(handlersOut, handlersInRow...))
	}

	if returnMiddlewaresOnly && handlerOut != nil {
		var asMiddleware = ConvertHandlerToMiddleware(handlerOut)
		middlewaresInRow = append(middlewaresInRow, asMiddleware)
		handlerOut = nil
	}

	return middlewaresInRow, handlerOut
}

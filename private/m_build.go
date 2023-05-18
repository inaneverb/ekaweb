package ekaweb_private

// BuildHandlerOut parses all provided arguments and generates an out
// middlewares and handler that could be registered in any cases.
func BuildHandlerOut(
	components []any,
	checkError Middleware,
	returnMiddlewaresOnly bool) ([]Middleware, Handler) {

	enableCheckErrorGlobal := checkError != nil

	// Apply components. The rules (about error check skipping) are:
	//
	// 1. First item in whole HTTP call sequence (middlewares + handlers)
	//    IS ALSO may be protected. It could be a call for HTTP group
	//    and there could be an error before.
	//
	// 2. Protecting (middleware or handler) meaning that this item
	//    will be called only if there was no error before.
	//
	// 3. For middlewares we can just add CheckError middleware
	//    to the middlewares call stack. Easy.
	//
	// 4. For handlers, we're wrapping handler by the CheckError() call,
	//    getting a new handler.

	var middlewaresInRow []Middleware
	var handlersInRow []Handler
	var handlersOut []Handler

	addMiddleware := func(middleware Middleware, performCheckError bool) {
		switch {
		case len(handlersInRow) > 0 && len(middlewaresInRow) > 0:
			tempHandler := MergeHandlers(handlersInRow)
			tempHandler = MergeMiddlewares(middlewaresInRow, tempHandler)
			handlersOut = append(handlersOut, tempHandler)
			handlersInRow = nil
			middlewaresInRow = nil

		case len(handlersInRow) > 0:
			handlersOut = append(handlersOut, handlersInRow...)
			handlersInRow = nil
		}

		asMiddlewareExtended, ok := middleware.(MiddlewareExtended)

		if performCheckError && (!ok || asMiddlewareExtended.CheckErrorBefore()) {
			middlewaresInRow = append(middlewaresInRow, checkError)
		}
		middlewaresInRow = append(middlewaresInRow, middleware)
	}

	addHandler := func(handler Handler, performCheckError bool) {
		asHandlerExtended, ok := handler.(HandlerExtended)

		if performCheckError && (!ok || asHandlerExtended.CheckErrorBefore()) {
			handler = checkError.Callback(handler)
		}
		handlersInRow = append(handlersInRow, handler)
	}

	for i, n := 0, len(components); i < n; i++ {
		asMiddleware := AsMiddleware(components[i])
		asHandler := AsHandler(components[i])

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
		asMiddleware := ConvertHandlerToMiddleware(handlerOut)
		middlewaresInRow = append(middlewaresInRow, asMiddleware)
		handlerOut = nil
	}

	return middlewaresInRow, handlerOut
}

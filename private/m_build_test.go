package ekaweb_private_test

import (
	"testing"

	"github.com/inaneverb/ekaweb/v2/private"
)

func TestBuildHandlerOut(t *testing.T) {

	var components []any
	for _, middleware := range genMiddlewares(3) {
		components = append(components, middleware)
	}
	for _, handler := range genHandlers(2) {
		components = append(components, handler)
	}

	checkError := NewPrintMiddleware("CheckError")

	middlewares, handler := ekaweb_private.BuildHandlerOut(components, checkError, false)
	handler = ekaweb_private.MergeMiddlewares(middlewares, handler)

	handlerInvoke(handler)
}

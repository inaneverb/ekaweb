package ekaweb_private_test

import (
	"strconv"
	"testing"

	"github.com/inaneverb/ekaweb/v2/private"
)

func genHandlers(n int) []ekaweb_private.Handler {
	handlers := make([]ekaweb_private.Handler, 0, n)
	for i := 0; i < n; i++ {
		handler := NewPrintHandler(strconv.Itoa(i))
		handlers = append(handlers, handler)
	}
	return handlers
}

func genMiddlewares(n int) []ekaweb_private.Middleware {
	middlewares := make([]ekaweb_private.Middleware, 0, n)
	for i := 0; i < n; i++ {
		middleware := NewPrintMiddleware(strconv.Itoa(i))
		middlewares = append(middlewares, middleware)
	}
	return middlewares
}

func genTestMergeHandlers(n int) func(*testing.T) {
	return genHandlerInvoke(ekaweb_private.MergeHandlers(genHandlers(n)))
}

func genTestMergeMiddlewares(n int) func(*testing.T) {
	handler := ekaweb_private.MergeHandlers(genHandlers(1))
	return genHandlerInvoke(ekaweb_private.MergeMiddlewares(genMiddlewares(n), handler))
}

func TestMergeHandlers(t *testing.T) {
	t.Run("Merge_Handlers_0", genTestMergeHandlers(0))
	t.Run("Merge_Handlers_1", genTestMergeHandlers(1))
	t.Run("Merge_Handlers_2", genTestMergeHandlers(2))
	t.Run("Merge_Handlers_5", genTestMergeHandlers(5))
}

func TestMergeMiddlewares(t *testing.T) {
	t.Run("Merge_Middlewares_0", genTestMergeMiddlewares(0))
	t.Run("Merge_Middlewares_1", genTestMergeMiddlewares(1))
	t.Run("Merge_Middlewares_2", genTestMergeMiddlewares(2))
	t.Run("Merge_Middlewares_5", genTestMergeMiddlewares(5))
}

package ekaweb_private_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/inaneverb/ekaweb/v2/private"
)

type httpDummyResponseWriter struct{}

func (*httpDummyResponseWriter) Header() http.Header         { return nil }
func (*httpDummyResponseWriter) Write(_ []byte) (int, error) { return 0, nil }
func (*httpDummyResponseWriter) WriteHeader(_ int)           {}

func GetRequestAndResponseWriter() (http.ResponseWriter, *http.Request) {
	return (*httpDummyResponseWriter)(nil), new(http.Request)
}

func NewPrintHandler(s string) ekaweb_private.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("\t--- INVOKE HANDLER: " + s)
	})
}

func NewPrintMiddleware(s string) ekaweb_private.Middleware {
	return ekaweb_private.MiddlewareFunc(func(next ekaweb_private.Handler) ekaweb_private.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Println("\t--- INVOKE MIDDLEWARE: " + s)
			next.ServeHTTP(w, r)
		})
	})
}

func genHandlerInvoke(handler ekaweb_private.Handler) func(*testing.T) {
	return func(_ *testing.T) {
		handlerInvoke(handler)
	}
}

func handlerInvoke(handler ekaweb_private.Handler) {
	w, r := GetRequestAndResponseWriter()
	handler.ServeHTTP(w, r)
}

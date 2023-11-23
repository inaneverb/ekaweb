package main

import (
	"fmt"
	"net/http"

	"github.com/inaneverb/ekaweb/framework/jrpc/v2"
	"github.com/inaneverb/ekaweb/v2"
)

func main() {
	var opts = []ekaweb.RouterOption{
		ekaweb.WithServerName("jrpc.example"),
		ekaweb.WithErrorHandler(errorHandler),
	}

	var r = ekaweb_jrpc.NewRouter(opts...).
		Use(middleware).
		Reg("main", handler)

	panic(http.ListenAndServe(":8081", r.Build()))
}

func middleware(next ekaweb.Handler) ekaweb.Handler {
	return ekaweb.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("HELLO FROM MIDDLEWARE")
		next.ServeHTTP(w, r)
	})
}

func handler(w http.ResponseWriter, r *http.Request) {

	type Response struct {
		RegisteredRoute string `json:"registered_route"`
		ActualRoute     string `json:"actual_route"`
		JRpcMethod      string `json:"jrpc_method"`
	}

	var _, method = ekaweb_jrpc.UkvsGetMeta(r)

	ekaweb.SendEncoded(w, r, ekaweb.StatusOK, Response{
		RegisteredRoute: ekaweb.RoutePath(r),
		ActualRoute:     r.RequestURI,
		JRpcMethod:      method,
	})
}

func errorHandler(w http.ResponseWriter, r *http.Request, err error) {
	var re ekaweb_jrpc.ResponseError
	re.Data = err.Error()
	re.FillMissedFields()
	ekaweb.SendEncoded(w, r, ekaweb.StatusInternalServerError, re)
}

package main

import (
	"net/http"

	"github.com/inaneverb/ekaweb"

	"github.com/inaneverb/ekaweb/extension/binding"
	"github.com/inaneverb/ekaweb/extension/respondent"

	"github.com/inaneverb/ekaweb/framework/chi"
	"github.com/inaneverb/ekaweb/framework/fasthttp"
	"github.com/inaneverb/ekaweb/framework/jrpc"

	"github.com/inaneverb/ekaweb/logger/rz"
)

func main() {

	var exp = ekaweb_respondent.NewCommonExpander()
	var resp = ekaweb_respondent.NewForHTTP(exp)

	var middlewareLog = ekaweb_rz.NewMiddleware(
		ekaweb_rz.WithStringExtractor("jrpc_method", jRpcExtMethod),
	)

	var jRpcHandler = ekaweb_jrpc.NewRouter(ekaweb.WithCoreInit(false)).
		Reg("method.name", handler).
		Build()

	var handler = ekaweb_chi.NewRouter(ekaweb.WithErrorHandler(resp)).
		Use(middlewareLog).
		Post("/method.name", handler).
		Post("/jrpc", jRpcHandler).
		Build()

	var srv = ekaweb_fasthttp.NewServer(
		ekaweb.WithHandler(handler),
		ekaweb.WithListenAddr(":8083"),
	)

	var err = srv.AsyncStart()
	if err != nil {
		panic(err)
	}

	select {}
}

func handler(w http.ResponseWriter, r *http.Request) {
	type Request struct {
		Name string `json:"name" binding:"required"`
		Age  int    `json:"age"  binding:"gte=18"`
	}

	type Response struct {
		Name string `json:"response_name"`
		Age  int    `json:"response_age"`
	}

	var req Request
	if err := ekaweb_bind.ScanAndValidateJSON(r, &req); err != nil {
		ekaweb.ErrorApply(r, err)
		return
	}

	var resp Response
	resp.Name = req.Name
	resp.Age = req.Age

	ekaweb.SendJSON(w, r, ekaweb.StatusOK, &resp)
}

func jRpcExtMethod(r *http.Request) string {
	return ekaweb_jrpc.RequestMethod(r)
}

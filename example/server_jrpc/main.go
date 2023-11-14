package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/inaneverb/ekaweb/extension/binding/v2"
	"github.com/inaneverb/ekaweb/extension/respondent/v2"
	"github.com/inaneverb/ekaweb/framework/chi/v2"
	"github.com/inaneverb/ekaweb/framework/fasthttp/v2"
	"github.com/inaneverb/ekaweb/framework/jrpc/v2"
	"github.com/inaneverb/ekaweb/logger/zerolog/v2"
	"github.com/inaneverb/ekaweb/v2"
)

func main() {

	const _40001 = "Bad request"

	var rcr = ekaweb_respondent.NewCommonReplacer()

	var rce = ekaweb_respondent.NewCommonExpander().
		// --------- HTTP 400: Bad request section --------- //
		WithCustomFillers(ekaweb_bind.ErrMalformedSource, ekaweb.StatusBadRequest, 40001, _40001, ekaweb_respondent.CACF_AutoErrorDetail).
		ExtractorFor(
			ekaweb_bind.ErrValidationFailed, true,
			ekaweb_bind.NewRespondentManifestExtractor(ekaweb.StatusBadRequest, 40001, _40001))

	var resp = ekaweb_respondent.NewForHTTP(rce, ekaweb_respondent.WithReplacer(rcr))

	var middlewareLog = ekaweb_zerolog.NewMiddleware(
		ekaweb_zerolog.WithLogger(log.Output(zerolog.ConsoleWriter{Out: os.Stdout})),
		ekaweb_zerolog.WithStringExtractor("jrpc_method", jRpcExtMethod),
	)

	var jRpcHandler = ekaweb_jrpc.NewRouter(ekaweb.WithCoreInit(false)).
		Reg("method.name", handlerFunc).
		Build()

	var optsRouter = []ekaweb.RouterOption{
		ekaweb.WithErrorHandler(resp),
		ekaweb.WithTrailingSlash(false, true),
		ekaweb.WithServerName("full_example"),
	}

	var handler = ekaweb_chi.NewRouter(optsRouter...).
		Use(middlewareLog).
		Get("/test", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, _ = w.Write([]byte("hello world"))
		})).
		Post("/method.name", handlerFunc).
		Post("/jrpc", jRpcHandler).
		Build()

	var srv = ekaweb_fasthttp.NewServer(
		ekaweb.WithHandler(H{handler}),
		ekaweb.WithListenAddr(":8083"),
	)

	var err = srv.AsyncStart()
	if err != nil {
		panic(err)
	}

	select {}
}

type H struct{ h http.Handler }

func (h H) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var now = time.Now()
	defer func() {
		var dur = time.Since(now)
		fmt.Printf("Full handler, elapse: %s (%d ns)\n", dur.String(), dur.Nanoseconds())
	}()
	h.h.ServeHTTP(w, r)
}

//func handlerFunc(w http.ResponseWriter, r *http.Request) {
//	type Request struct {
//		Name string `json:"name" binding:"required"`
//		Age  int    `json:"age"  binding:"gte=18"`
//	}
//
//	type Response struct {
//		Name string `json:"response_name"`
//		Age  int    `json:"response_age"`
//	}
//
//	var req Request
//	if err := ekaweb_bind.ScanAndValidateJSON(r, &req); err != nil {
//		ekaweb.ErrorApply(r, err)
//		return
//	}
//
//	var resp Response
//	resp.Name = req.Name
//	resp.Age = req.Age
//
//	ekaweb.SendJSON(w, r, ekaweb.StatusOK, &resp)
//}

func handlerFunc(w http.ResponseWriter, r *http.Request) {
	var now = time.Now()
	defer func() {
		var dur = time.Since(now)
		fmt.Printf("Route handler, elapse: %s (%d ns)\n", dur.String(), dur.Nanoseconds())
	}()

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

	var responses []Response = nil
	if true {
		responses = append(responses, resp)
	}

	var code = ekaweb.StatusNoContent
	if len(responses) > 0 || jRpcExtMethod(r) != "" {
		code = ekaweb.StatusOK
	}

	ekaweb.SendEncoded(w, r, code, responses)
}

func jRpcExtMethod(r *http.Request) string {
	var _, method = ekaweb_jrpc.UkvsGetMeta(r)
	return method
}

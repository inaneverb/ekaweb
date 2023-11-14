package main

import (
	"log"
	"net/http"
	"reflect"
	"time"

	"github.com/inaneverb/ekaweb/extension/binding/v2"
	"github.com/inaneverb/ekaweb/extension/respondent/v2"
	"github.com/inaneverb/ekaweb/framework/chi/v2"
	"github.com/inaneverb/ekaweb/framework/fasthttp/v2"
	"github.com/inaneverb/ekaweb/logger/rz/v2"
	"github.com/inaneverb/ekaweb/v2"

	"github.com/go-playground/validator/v10"
	"github.com/gofrs/uuid"
)

func main() {

	var exp = ekaweb_respondent.NewCommonExpander()
	var resp = ekaweb_respondent.NewForHTTP(exp)

	_ = ekaweb_bind.RegisterCustomTypes(cvTypeUUID)

	var middlewareLog = ekaweb_rz.NewMiddleware()

	var handler = ekaweb_chi.NewRouter(ekaweb.WithErrorHandler(resp)).
		Use(middlewareLog).
		Post("/", handler).
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
		// Only one of ID or When must be present

		ID   uuid.UUID `json:"id"   binding:"required_without=When,excluded_with=When"`
		When time.Time `json:"when" binding:"required_without=ID,excluded_with=ID"`
	}

	var req Request
	if err := ekaweb_bind.ScanAndValidateJSON(r, &req); err != nil {
		ekaweb.ErrorApply(r, err)
		return
	}

	if !req.ID.IsNil() {
		log.Println("ID is used")
	} else {
		log.Println("When is used")
	}

	ekaweb.SendEncoded(w, r, ekaweb.StatusOK, &req)
}

func cvTypeUUID(v *validator.Validate) error {

	var validationFn = func(field reflect.Value) any {
		if id, ok := field.Interface().(uuid.UUID); ok && !id.IsNil() {
			return id.String()
		} else {
			return ""
		}
	}

	v.RegisterCustomTypeFunc(validationFn, uuid.UUID{})
	return nil
}

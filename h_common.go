package ekaweb

import (
	"context"
	"io"
	"net/http"
	"strings"

	"github.com/inaneverb/ekacore/ekaunsafe/v4"

	"github.com/inaneverb/ekaweb/private"
)

////////////////////////////////////////////////////////////////////////////////
///// HTTP Request headers methods /////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////

func Header(r *http.Request, key string) string {
	return HeaderWithDefault(r, key, "")
}

func HeaderWithDefault(r *http.Request, key, defaultValue string) string {
	if r != nil && r.Header != nil && key != "" {
		if value := r.Header.Get(key); value != "" {
			return value
		}
	}
	return defaultValue
}

func HeaderContain(r *http.Request, key, value string) bool {
	if value == "" {
		return false
	} else if gotValue := Header(r, key); gotValue == "" {
		return false
	} else {
		return strings.Contains(gotValue, value)
	}
}

//func Accept(r *http.Request, offer string) bool {
//	panic("TODO: Not implemented yet")
//}
//
//func Accepts(r *http.Request, offers ...string) string {
//	panic("TODO: Not implemented yet")
//}
//
//func AcceptCharset(r *http.Request, offer string) bool {
//	panic("TODO: Not implemented yet")
//}
//
//func AcceptsCharsets(r *http.Request, offers ...string) string {
//	panic("TODO: Not implemented yet")
//}
//
//func AcceptEncoding(r *http.Request, offer string) bool {
//	panic("TODO: Not implemented yet")
//}
//
//func AcceptsEncodings(r *http.Request, offers ...string) string {
//	panic("TODO: Not implemented yet")
//}
//
//func AcceptLanguage(r *http.Request, offer string) bool {
//	panic("TODO: Not implemented yet")
//}
//
//func AcceptsLanguages(r *http.Request, offers ...string) string {
//	panic("TODO: Not implemented yet")
//}

////////////////////////////////////////////////////////////////////////////////
///// HTTP Request user key value storage methods //////////////////////////////
////////////////////////////////////////////////////////////////////////////////

func UserVarLookup(r *http.Request, key any) (any, bool) {
	return UserVarLookupByContext(r.Context(), key)
}

func UserVarGet(r *http.Request, key any) any {
	return UserVarGetByContext(r.Context(), key)
}

func UserVarGetOrDefault(r *http.Request, key, defaultValue any) any {
	return UserVarGetOrDefaultByContext(r.Context(), key, defaultValue)
}

func UserVarInsert(r *http.Request, key, value any) {
	UserVarInsertByContext(r.Context(), key, value)
}

func UserVarLookupByContext(ctx context.Context, key any) (any, bool) {
	return ekaweb_private.UkvsLookup(ctx, key)
}

func UserVarGetByContext(ctx context.Context, key any) any {
	return ekaweb_private.UkvsGet(ctx, key)
}

func UserVarGetOrDefaultByContext(ctx context.Context, key, defaultValue any) any {
	return ekaweb_private.UkvsGetOrDefault(ctx, key, defaultValue)
}

func UserVarInsertByContext(ctx context.Context, key, value any) {
	ekaweb_private.UkvsInsert(ctx, key, value)
}

////////////////////////////////////////////////////////////////////////////////
///// HTTP Request user key value storage specific entity methods //////////////
////////////////////////////////////////////////////////////////////////////////

func URLVarGet(r *http.Request, key string) string {
	return ekaweb_private.UkvsGetOrDefault(r.Context(), key, "").(string)
}

func RoutePath(r *http.Request) string {
	return ekaweb_private.UkvsGetOriginalPath(r.Context())
}

////////////////////////////////////////////////////////////////////////////////
///// HTTP Response generators /////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////

func SendJSON(
	w http.ResponseWriter, r *http.Request,
	statusCode int, obj any) {

	var data, err = ekaweb_private.EncodeJSON(r.Context(), obj)
	if err != nil {
		ErrorApply(r, err)
		return
	}

	SendRaw(w, r, statusCode, MIMEApplicationJSONCharsetUTF8, data)
}

func SendString(
	w http.ResponseWriter, r *http.Request,
	statusCode int, data string) {

	dataBytes := ekaunsafe.StringToBytes(data)
	SendRaw(w, r, statusCode, MIMETextPlainCharsetUTF8, dataBytes)
}

func SendEmpty(w http.ResponseWriter, r *http.Request, statusCode int) {
	SendRaw(w, r, statusCode, "", nil)
}

func SendRaw(
	w http.ResponseWriter, r *http.Request,
	statusCode int, mimeType string, data []byte) {

	customResponseHeaders := ekaweb_private.UkvsGetResponseCustomHeaders(r.Context())
	CopyHeaders(w.Header(), customResponseHeaders)

	if mimeType != "" {
		w.Header().Set(HeaderContentType, mimeType)
	}

	w.WriteHeader(statusCode)
	if len(data) != 0 {
		_, _ = w.Write(data) // TODO: Shall we ignore err from Write here?
	}
}

func SendStream(
	w http.ResponseWriter, r *http.Request,
	statusCode int, mimeType string, stream io.Reader) {

	SendRaw(w, r, statusCode, mimeType, nil)
	_, _ = io.Copy(w, stream)
}

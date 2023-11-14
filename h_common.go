package ekaweb

import (
	"context"
	"io"
	"net/http"
	"net/textproto"
	"strings"

	"github.com/inaneverb/ekacore/ekaarr/v4"
	"github.com/inaneverb/ekacore/ekaunsafe/v4"

	"github.com/inaneverb/ekaweb/v2/private"
)

////////////////////////////////////////////////////////////////////////////////
///// HTTP Request headers methods /////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////

// Header returns value of HTTP header with given 'key' extracting it from
// passed http.Request. Returns empty string if no such header found.
func Header(r *http.Request, key string) string {
	return HeaderWithDefault(r, key, "")
}

// HeaderWithDefault is the same as just Header(), but returns fallback
// value if no such header found.
func HeaderWithDefault(r *http.Request, key, defaultValue string) string {
	if r != nil && r.Header != nil && key != "" {
		if value := r.Header.Get(key); value != "" {
			return value
		}
	}
	return defaultValue
}

// HeaderContain reports whether HTTP header with given 'key' is exists
// in provided http.Request and its value has given 'value' as substring.
// If given 'value' is empty, false is always returned.
func HeaderContain(r *http.Request, key, value string) bool {
	if value == "" {
		return false
	} else if gotValue := Header(r, key); gotValue == "" {
		return false
	} else {
		return strings.Contains(gotValue, value)
	}
}

// HeadersMerge merges HTTP headers and theirs values. Returns resulted set.
// You can enable 'checkDuplicates' to ensure there's no duplicated HTTP
// header's values for each key.
// It decreases performance, but increases stability.
//
// If any of values ('a' or 'b') is nil, an opposite one is returned
// (thus, if both of them is nil, nil is returned).
//
// If 'a' is not nil, it copies headers from 'b' to 'a', returning modified 'a'.
func HeadersMerge(a, b http.Header, checkDuplicates bool) http.Header {

	if len(b) == 0 {
		return a
	}

	// Both not empty for now. Extend 'a' by 'b'.

	for key, bValues := range b {
		key = textproto.CanonicalMIMEHeaderKey(key)

		var aValues = a[key]
		if len(aValues) == 0 {
			a[key] = bValues // just place new key with b's values
			continue
		}

		aValues = append(aValues, bValues...) // faster, than .Add()

		if checkDuplicates {
			aValues = ekaarr.Distinct(aValues)
		}

		a[key] = aValues
	}

	return a
}

// CopyHeaders is the legacy version of HeadersMerge().
// Deprecated: Use HeadersMerge() instead.
func CopyHeaders(to, from http.Header) http.Header {
	return HeadersMerge(to, from, true)
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

// UserVarLookup provides you access to the http.Request's specific
// custom key-value storage, allowing you to retrieve previously stored data
// for the given 'key'.
// If no data with such 'key' was stored, (nil, false) is returned.
// Thus, you may to distinct, whether nil was stored, or no value at all.
func UserVarLookup(r *http.Request, key any) (any, bool) {
	return ekaweb_private.UkvsLookup(r.Context(), key)
}

// UserVarGet is the same as UserVarLookup(), but returns nil if no data
// with given 'key' was stored.
func UserVarGet(r *http.Request, key any) any {
	return ekaweb_private.UkvsGet(r.Context(), key)
}

// UserVarGetOrDefault is the same as just UserVarGet(), but allows you
// to provide fallback value, that will be returned if no data with given 'key'
// was stored.
//
// NOTE.
// This method doesn't rely on UserVarGet(). Thus, if you store nil
// with some key, you will get nil for the same key by this method,
// not a fallback value.
func UserVarGetOrDefault(r *http.Request, key, defaultValue any) any {
	return ekaweb_private.UkvsGetOrDefault(r.Context(), key, defaultValue)
}

// UserVarInsert provides you access to the http.Request's specific
// custom key-value storage, allowing you to store some 'value' with given 'key'
// in the storage.
//
// Later, you can access it by UserVarLookup(), UserVarGet()
// or UserVarGetOrDefault().
func UserVarInsert(r *http.Request, key, value any) {
	ekaweb_private.UkvsInsert(r.Context(), key, value)
}

////////////////////////////////////////////////////////////////////////////////

// UserVarLookupByContext is the same as UserVarLookup(), but works directly
// with context.Context instead.
// Deprecated: Access values by UserVarLookup() is more recommended.
func UserVarLookupByContext(ctx context.Context, key any) (any, bool) {
	// WARNING!
	// DO NOT REMOVE THIS FUNCTION, ALTHOUGH DEPRECATION MARK!
	// IN SOME CASES THERE'S NO OTHER WAY THAN STORE context.Context
	// AND ACCESS VALUES THROUGH THEM!
	return ekaweb_private.UkvsLookup(ctx, key)
}

// UserVarGetByContext is the same as UserVarGet(), but works directly
// with context.Context instead.
// Deprecated: Access values by UserVarGet() is more recommended.
func UserVarGetByContext(ctx context.Context, key any) any {
	// WARNING!
	// DO NOT REMOVE THIS FUNCTION, ALTHOUGH DEPRECATION MARK!
	// IN SOME CASES THERE'S NO OTHER WAY THAN STORE context.Context
	// AND ACCESS VALUES THROUGH THEM!
	return ekaweb_private.UkvsGet(ctx, key)
}

// UserVarGetOrDefaultByContext the same as UserVarGetOrDefault(),
// but works directly with context.Context instead.
// Deprecated: Access values by UserVarGetOrDefault() is more recommended.
func UserVarGetOrDefaultByContext(ctx context.Context, key, defaultValue any) any {
	// WARNING!
	// DO NOT REMOVE THIS FUNCTION, ALTHOUGH DEPRECATION MARK!
	// IN SOME CASES THERE'S NO OTHER WAY THAN STORE context.Context
	// AND ACCESS VALUES THROUGH THEM!
	return ekaweb_private.UkvsGetOrDefault(ctx, key, defaultValue)
}

// UserVarInsertByContext the same as UserVarInsert(), but works directly
// with context.Context instead.
// Deprecated: Access values by UserVarInsert() is more recommended.
func UserVarInsertByContext(ctx context.Context, key, value any) {
	// WARNING!
	// DO NOT REMOVE THIS FUNCTION, ALTHOUGH DEPRECATION MARK!
	// IN SOME CASES THERE'S NO OTHER WAY THAN STORE context.Context
	// AND ACCESS VALUES THROUGH THEM!
	ekaweb_private.UkvsInsert(ctx, key, value)
}

////////////////////////////////////////////////////////////////////////////////
///// HTTP Request user key value storage specific entity methods //////////////
////////////////////////////////////////////////////////////////////////////////

// URLVarGet returns the value that is matched against given 'key' in URL
// from given http.Request.
//
// Like, if you did register an HTTP route "/users/{user_id}/register",
// you can obtain, what "user_id" was actually used by this method in your
// HTTP controller, using "user_id" as a 'key'.
func URLVarGet(r *http.Request, key string) string {
	return ekaweb_private.UkvsGetOrDefault(r.Context(), key, "").(string)
}

// RoutePath allows you to get registered (!!) HTTP route path, for which
// your HTTP controller was called.
//
// Thus, if you did register an HTTP route "/users/{user_id}/register",
// and you have an HTTP request with "/users/alice/register" URL,
// this method will return "/users/{user_id}/register".
//
// NOTE.
// In case it's not possible to retrieve registered HTTP route, this method
// returns an actual URL, that is used by the client in their request.
func RoutePath(r *http.Request) string {
	var path = ekaweb_private.UkvsGetOriginalPath(r.Context())
	if path == "" {
		path = r.URL.Path
	}
	return path
}

////////////////////////////////////////////////////////////////////////////////
///// HTTP Response generators /////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////

// SendEncoded is the same as just SendEncodedWithMIME(), but uses JSON MIME
// no matter what codec is stored (assuming, that ~90% users uses JSON).
// Thus, MIMEApplicationJSONCharsetUTF8 constant will be used.
//
// WARNING!
// YOU MAY GET UNEXPECTED BEHAVIOUR IF YOUR STORED CODEC WITH ANY PROBABILITY
// COULD BE SOMETHING ELSE BUT JSON. SendEncodedWithMIME() is recommended.
func SendEncoded(
	w http.ResponseWriter, r *http.Request, statusCode int, obj any) {

	SendEncodedWithMIME(w, r, statusCode, MIMEApplicationJSONCharsetUTF8, obj)
}

// SendEncodedWithMIME encodes given 'obj' with codec that is stored inside
// http.Request's context.Context, finalizing sending HTTP response.
//
// It also uses given 'statusCode' as HTTP response code and 'mimeType' as MIME.
// It's ok to pass empty MIME type (no HeaderContentEncoding will be passed).
//
// If given 'obj' is nil, the data that will be sent depends on stored codec.
// If you want an empty HTTP response body, use SendEmpty() explicitly.
func SendEncodedWithMIME(
	w http.ResponseWriter, r *http.Request,
	statusCode int, mimeType string, obj any) {

	// For applying headers (status code + MIME)
	SendRaw(w, statusCode, mimeType, nil)

	if err := ekaweb_private.EncodeStream(r.Context(), w, obj); err != nil {
		// NOTE:
		// We could, potentially just call ErrorApply() w/o "if" statement,
		// but this call is not free and has context-searching operation.
		ErrorApply(r, err)
	}
}

// SendJSON is just the legacy named version of SendEncoded().
// Deprecated: Use SendEncoded() instead.
func SendJSON(w http.ResponseWriter, r *http.Request, statusCode int, obj any) {
	SendEncoded(w, r, statusCode, obj)
}

// SendString is the same as just SendRaw(), but has zero-cost conversion
// string -> bytes ([]byte) and uses MIMETextPlainCharsetUTF8 as MIME type.
func SendString(
	w http.ResponseWriter, statusCode int, data string) {

	var dataBytes = ekaunsafe.StringToBytes(data)
	SendRaw(w, statusCode, MIMETextPlainCharsetUTF8, dataBytes)
}

// SendEmpty sends an empty HTTP response (no HTTP body) with given 'statusCode'.
// It also DOES NOT send any MIME type, because there's no HTTP body in response.
func SendEmpty(w http.ResponseWriter, statusCode int) {
	SendRaw(w, statusCode, "", nil)
}

// SendRaw sends an HTTP response with given HTTP 'statusCode' and 'mimeType'.
// It also uses given 'data' as HTTP response body (but only if it's not empty).
// If MIME is empty, it won't be sent to the client.
func SendRaw(
	w http.ResponseWriter, statusCode int, mimeType string, data []byte) {

	if mimeType != "" {
		w.Header().Set(HeaderContentType, mimeType)
	}

	w.WriteHeader(statusCode)
	if len(data) != 0 {
		_, _ = w.Write(data) // TODO: Shall we ignore err from Write() here?
	}
}

// SendStream is the same as just SendRaw(), but it uses given io.Reader instead.
// Also, because there's no way to check, whether data is provided, MIME header
// will be sent anyway, and the presence of HTTP body is completely on your side.
func SendStream(
	w http.ResponseWriter, statusCode int, mimeType string, stream io.Reader) {

	// For applying headers (status code + MIME)
	SendRaw(w, statusCode, mimeType, nil)

	_, _ = io.Copy(w, stream) // TODO: Shall we ignore err from Copy() here?
}

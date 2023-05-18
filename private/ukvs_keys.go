package ekaweb_private

import (
	"context"
	"net/http"
)

type (
	_UkvsKeyError                  struct{}
	_UkvsKeyErrorDetail            struct{}
	_UkvsKeyJSONEncoderDecoder     struct{}
	_UkvsKeyResponseCustomHeaders  struct{}
	_UkvsKeyHijackedConn           struct{}
	_UkvsKeyOriginalPath           struct{}
	_UkvsKeyPathNotFoundNotAllowed struct{}
)

func UkvsGetUserError(ctx context.Context) error {
	if elem, ok := UkvsGet(ctx, _UkvsKeyError{}).(error); ok && elem != nil {
		return elem
	} else {
		return nil
	}
}

func UkvsInsertUserError(ctx context.Context, err error) {
	if err != nil {
		UkvsInsertIfNone(ctx, _UkvsKeyError{}, err)
	}
}

func UkvsGetUserErrorDetail(ctx context.Context) string {
	if elem, ok := UkvsGet(ctx, _UkvsKeyErrorDetail{}).(string); ok {
		return elem
	} else {
		return ""
	}
}

func UkvsInsertUserErrorDetail(ctx context.Context, detail string) {
	UkvsInsertIfNone(ctx, _UkvsKeyErrorDetail{}, detail)
}

func UkvsRemoveUserErrorFull(ctx context.Context) {
	UkvsRemoveNoInfo(ctx, _UkvsKeyError{})
	UkvsRemoveNoInfo(ctx, _UkvsKeyErrorDetail{})
}

func UkvsGetJSONEncoderDecoder(ctx context.Context) *RouterOptionCustomJSON {
	return UkvsGetOrDefault(ctx, _UkvsKeyJSONEncoderDecoder{}, (*RouterOptionCustomJSON)(nil)).(*RouterOptionCustomJSON)
}

func UkvsInsertJSONEncoderDecoder(ctx context.Context, opt *RouterOptionCustomJSON) {
	UkvsInsert(ctx, _UkvsKeyJSONEncoderDecoder{}, opt)
}

func UkvsGetResponseCustomHeaders(ctx context.Context) http.Header {
	return UkvsGetOrDefault(ctx, _UkvsKeyResponseCustomHeaders{}, http.Header(nil)).(http.Header)
}

func UkvsInsertResponseCustomHeaders(ctx context.Context, headers http.Header) {
	UkvsInsert(ctx, _UkvsKeyResponseCustomHeaders{}, headers)
}

func UkvsIsConnectionHijacked(ctx context.Context) bool {
	return UkvsGetOrDefault(ctx, _UkvsKeyHijackedConn{}, false).(bool)
}

func UkvsMarkConnectionAsHijacked(ctx context.Context) {
	UkvsInsert(ctx, _UkvsKeyHijackedConn{}, true)
}

func UkvsGetOriginalPath(ctx context.Context) string {
	return UkvsGetOrDefault(ctx, _UkvsKeyOriginalPath{}, "").(string)
}

func UkvsInsertOriginalPath(ctx context.Context, originalPath string) {
	UkvsInsert(ctx, _UkvsKeyOriginalPath{}, originalPath)
}

func UkvsIsPathNotFoundOrNotAllowed(ctx context.Context) bool {
	return UkvsGetOrDefault(ctx, _UkvsKeyPathNotFoundNotAllowed{}, false).(bool)
}

func UkvsSetPathNotFoundOrNotAllowed(ctx context.Context, notFoundOrNotAllowed bool) {
	UkvsInsert(ctx, _UkvsKeyPathNotFoundNotAllowed{}, notFoundOrNotAllowed)
}

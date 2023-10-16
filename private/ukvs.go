package ekaweb_private

import (
	"context"
)

type (
	// _Ukvs is an object that is stored inside each context.Context
	// saving some user data + metadata correspondent to current HTTP context.
	_Ukvs struct {
		mgr   *UkvsManager      // ptr to parent manager
		m     UkvsMap           // for user-specific values
		err   error             // error as is
		errD  string            // error detail (description)
		uri   string            // original URI path (with variables)
		flags uint32            // state & behaviour of current context
		codec RouterOptionCodec // encoder + decoder that used to operate
	}
)

////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////

const (
	_UkvsFlagConnHijacked uint32 = 1 << iota
	_UkvsFlagNotFound
	_UkvsFlagNotAllowed
	_UkvsFlagIsStolen
)

////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////

func UkvsLookup(ctx context.Context, key any) (elem any, found bool) {
	return ukvsGet(ctx).m.Lookup(key)
}

func UkvsGet(ctx context.Context, key any) (elem any) {
	return ukvsGet(ctx).m.Get(key)
}

func UkvsGetOrDefault(ctx context.Context, key, defaultValue any) any {
	if value, ok := UkvsLookup(ctx, key); ok {
		return value
	} else {
		return defaultValue
	}
}

func UkvsSwap(ctx context.Context, key, value any) (prev any, was bool) {
	return ukvsGet(ctx).m.Swap(key, value)
}

func UkvsInsert(ctx context.Context, key, value any) {
	ukvsGet(ctx).m.Set(key, value, true)
}

func UkvsInsertIfNone(ctx context.Context, key, value any) {
	ukvsGet(ctx).m.Set(key, value, false)
}

func UkvsRemove(ctx context.Context, key any) (prev any, was bool) {
	return ukvsGet(ctx).m.Remove(key)
}

////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////

func UkvsGetUserError(ctx context.Context) error {
	return ukvsGet(ctx).err
}

func UkvsInsertUserError(ctx context.Context, err error) {
	ukvsGet(ctx).err = err
}

func UkvsGetUserErrorDetail(ctx context.Context) string {
	return ukvsGet(ctx).errD
}

func UkvsInsertUserErrorDetail(ctx context.Context, detail string) {
	ukvsGet(ctx).errD = detail
}

func UkvsRemoveUserErrorFull(ctx context.Context) {
	ukvsGet(ctx).err = nil
	ukvsGet(ctx).errD = ""
}

func UkvsGetCodec(ctx context.Context) RouterOptionCodec {
	return ukvsGet(ctx).codec
}

func UkvsInsertCodec(ctx context.Context, codec RouterOptionCodec) {
	ukvsGet(ctx).codec = codec
}

func UkvsIsConnectionHijacked(ctx context.Context) bool {
	return ukvsGet(ctx).flags&_UkvsFlagConnHijacked != 0
}

func UkvsMarkConnectionAsHijacked(ctx context.Context) {
	ukvsGet(ctx).flags |= _UkvsFlagConnHijacked
}

func UkvsGetOriginalPath(ctx context.Context) string {
	return ukvsGet(ctx).uri
}

func UkvsInsertOriginalPath(ctx context.Context, originalPath string) {
	ukvsGet(ctx).uri = originalPath
}

func UkvsIsPathNotFoundOrNotAllowed(ctx context.Context) bool {
	return ukvsGet(ctx).flags&(_UkvsFlagNotFound|_UkvsFlagNotAllowed) != 0
}

func UkvsSetPathNotFoundOrNotAllowed(ctx context.Context, notFoundOrNotAllowed bool) {
	if notFoundOrNotAllowed {
		ukvsGet(ctx).flags |= _UkvsFlagNotFound | _UkvsFlagNotAllowed
	} else {
		ukvsGet(ctx).flags &^= _UkvsFlagNotFound | _UkvsFlagNotAllowed
	}
}

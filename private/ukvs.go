package ekaweb_private

import (
	"context"

	"github.com/inaneverb/ekacore/ekaunsafe/v4"
	"github.com/xaionaro-go/atomicmap"
)

type _UkvsContextKey struct{}

func UkvsInit(ctx context.Context) context.Context {
	return context.WithValue(ctx, _UkvsContextKey{}, atomicmap.NewWithArgs(16))
}

func UkvsPropagate(original, base context.Context) context.Context {
	return context.WithValue(base, _UkvsContextKey{}, ukvsMap(original))
}

func UkvsDestroy(_ context.Context) {}

func UkvsLookup(ctx context.Context, key any) (elem any, found bool) {
	return ukvsHandleErr(ukvsMap(ctx).Get(ukvsPrepareUserKey(key)))
}

func UkvsGet(ctx context.Context, key any) (elem any) {
	return ukvsIgnoreErr(ukvsMap(ctx).Get(ukvsPrepareUserKey(key)))
}

func UkvsGetOrDefault(ctx context.Context, key, defaultValue any) any {
	if value, ok := UkvsLookup(ctx, key); ok {
		return value
	} else {
		return defaultValue
	}
}

func UkvsSwap(ctx context.Context, key, value any) (prev any, was bool) {
	return ukvsHandleErr(ukvsMap(ctx).Swap(ukvsPrepareUserKey(key), value))
}

func UkvsInsert(ctx context.Context, key, value any) {
	_ = ukvsMap(ctx).Set(ukvsPrepareUserKey(key), value)
}

func UkvsInsertIfNone(ctx context.Context, key, value any) {
	var m = ukvsMap(ctx)

	key = ukvsPrepareUserKey(key)
	if _, err := m.Get(key); err == atomicmap.NotFound {
		_ = m.Set(key, value)
	}
}

func UkvsRemove(ctx context.Context, key any) (prev any, was bool) {
	_ = ukvsMap(ctx).UnsetIf(ukvsPrepareUserKey(key), func(value any) bool {
		prev = value
		was = true
		return true
	})
	return
}

func UkvsRemoveNoInfo(ctx context.Context, key any) (removed bool) {
	return ukvsMap(ctx).Unset(ukvsPrepareUserKey(key)) == nil
}

////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////

func ukvsMap(ctx context.Context) atomicmap.Map {
	return ctx.Value(_UkvsContextKey{}).(atomicmap.Map)
}

func ukvsHandleErr(elem any, err error) (any, bool) {
	return elem, err == nil
}

func ukvsIgnoreErr(elem any, _ error) any {
	return elem
}

////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////

func ukvsPrepareUserKey(key any) any {
	switch i := ekaunsafe.UnpackInterface(key); i.Type {
	case ekaunsafe.RTypeString():
		return key
	default:
		return i.Type
	}
}

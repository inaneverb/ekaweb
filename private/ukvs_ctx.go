package ekaweb_private

import (
	"context"

	"github.com/inaneverb/ekacore/ekaunsafe/v4"
)

type (
	_UkvsContext struct {
		context.Context
		kvs *_Ukvs
	}

	// _UkvsContextKey is a key for context.Context to store _Ukvs.
	_UkvsContextKey struct{}
)

var (
	rtypeContext    = ekaunsafe.RTypeOf(_UkvsContext{})
	rtypeContextKey = ekaunsafe.RTypeOf((*_UkvsContextKey)(nil))
)

func (c _UkvsContext) Value(key any) any {
	if ekaunsafe.UnpackInterface(key).Type == rtypeContextKey {
		return c.kvs
	} else {
		return c.Context.Value(key)
	}
}

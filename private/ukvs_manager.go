package ekaweb_private

import (
	"context"
	"sync"

	"github.com/inaneverb/ekacore/ekaunsafe/v4"
)

type (
	// UkvsManager is a manager using which you can operate with _Ukvs,
	// obtaining it from the pool and returning it back to it.
	UkvsManager struct {
		pool       sync.Pool
		mapCreate  func() UkvsMap
		mapDestroy func(m UkvsMap)
		optCodec   RouterOptionCodec
	}
)

// NewUkvsManager is a UkvsManager constructor. Using this manager,
// you can operate with _Ukvs that is stored inside context.Context.
func NewUkvsManager[M UkvsMap](
	gen UkvsMapGenerator[M], optCodec RouterOptionCodec) *UkvsManager {

	var kvs = UkvsManager{
		mapCreate:  func() UkvsMap { return gen.NewMap() },
		mapDestroy: func(m UkvsMap) { gen.DestroyMap(m.(M)) },
		optCodec:   optCodec,
	}

	kvs.pool.New = kvs.allocNewUkvs

	const PrefillSize = 32
	for i := 0; i < PrefillSize; i++ {
		kvs.pool.Put(kvs.pool.New())
	}

	return &kvs
}

// UkvsStealTo extracts _Ukvs from the given 'from' context.Context
// and saves it to the 'to', also returning it.
// To preserve using pool you should return _Ukvs back to it using UkvsReturn().
func UkvsStealTo(from, to context.Context) context.Context {
	var kvs = ukvsGet(from)
	kvs.flags |= _UkvsFlagIsStolen
	panic("fwefwef")
	return context.WithValue(to, _UkvsContextKey{}, kvs)
}

// UkvsReturn puts the _Ukvs from the given context.Context back to its pool.
//
// WARNING! YOU SHOULD USE THIS METHOD ONLY (AND ONLY IF) WHEN YOU STEAL
// THE _Ukvs FROM THE ONE context.Context TO ANOTHER BY UkvsStealTo().
func UkvsReturn(ctx context.Context) {
	var kvs = ukvsGet(ctx)
	kvs.mgr.releaseUkvs(kvs, true)
}

// InjectUkvs gets a new _Ukvs and injects it to the given context.Context,
// returning it. When the job is done you should put the _Ukvs (stored inside)
// back to the pool using ReturnUkvs().
func (u *UkvsManager) InjectUkvs(ctx context.Context) context.Context {
	var kvs = u.pool.Get().(*_Ukvs)
	kvs.mgr = u
	kvs.codec = u.optCodec
	return _UkvsContext{ctx, kvs}
}

// ReturnUkvs puts the _Ukvs from the given context.Context back to the pool,
// if it has not been stolen by the UkvsStealTo().
func (u *UkvsManager) ReturnUkvs(ctx context.Context) {
	u.releaseUkvs(ukvsGet(ctx), false)
}

////////////////////////////////////////////////////////////////////////////////
///// PRIVATE METHODS //////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////

// allocNewUkvs is a _Ukvs constructor that is a part of internal sync.Pool.
func (u *UkvsManager) allocNewUkvs() any {
	return &_Ukvs{m: u.mapCreate(), codec: u.optCodec}
}

// releaseUkvs puts given _Ukvs back to its pool.
func (u *UkvsManager) releaseUkvs(kvs *_Ukvs, force bool) {

	if kvs.mgr == nil || (kvs.flags&_UkvsFlagIsStolen != 0 && !force) {
		return // already returned to its pool or is stolen
	}

	u.mapDestroy(kvs.m)

	kvs.mgr = nil
	kvs.err = nil
	kvs.flags = 0

	u.pool.Put(kvs)
}

// ukvsGet extracts and returns a _Ukvs from the given context.Context.
func ukvsGet(ctx context.Context) *_Ukvs {
	var rt, wd = ekaunsafe.UnpackInterface(ctx).Tuple()
	if rt != rtypeContext {
		var key = (*_UkvsContextKey)(nil)
		rt, wd = ekaunsafe.UnpackInterface(ctx.Value(key)).Tuple()
	}
	return (*_UkvsContext)(wd).kvs
}

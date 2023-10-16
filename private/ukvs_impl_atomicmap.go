package ekaweb_private

/*
================================================================================
 DEPRECATED DEPRECATED DEPRECATED DEPRECATED DEPRECATED DEPRECATED DEPRECATED
================================================================================

 THIS IMPLEMENTATION IS CONSIDERED DEPRECATED AND FULLY COMMENTED TO AVOID
 UNNECESSARY DEPENDENCY IN GO.MOD FILE.
 THE MAIN REASONS OF FULLY REJECTED THIS IMPLEMENTATION IS BENCHMARKING RESULT:

 ATOMIC MAP LOSE ON ALMOST ALL BENCHMARK CASES.
 IT STARTS DOING WELL IN THE CASES WHEN THE QUANTITY OF STORED VALUES IS HIGH
 (LIKELY > 40-50) AND ONLY FOR RANDOM-ACCESS OPERATION.

 DESPITE THE FACT IT HAS THREAD-SAFETY OPPOSITE TO ITS COMPETITORS,
 THREAD-SAFETY NOT REALLY A REQUIRED FEATURE DURING ACCESSING KEY-VALUE STORAGE
 USING GOLANG'S context.Context.

 ANOTHER PROBLEM IS THAT WE CANNOT ZERO THE WHOLE MAP (REMOVE ALL VALUES)
 USING SOME EASY OPERATION. WE HAVE TO ITERATE OVER ALL KEYS AND REMOVE THEM
 ONE BY ONE.

================================================================================
 DEPRECATED DEPRECATED DEPRECATED DEPRECATED DEPRECATED DEPRECATED DEPRECATED
================================================================================
*/

/*

import (
	"github.com/xaionaro-go/atomicmap"
)

type (
	// _UkvsImplAtomicMap implements UkvsMap based on external thread-safe
	// lock-free (partially) map: https://github.com/xaionaro-go/atomicmap .
	_UkvsImplAtomicMap struct{ atomicmap.Map }

	// _UkvsImplAtomicMapGenerator implements UkvsMapGenerator,
	// working with _UkvsImplAtomicMap.
	_UkvsImplAtomicMapGenerator struct{}
)

////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////

// NewUkvsMapGeneratorAtomicMap returns a new UkvsMapGenerator that operates
// with UkvsMap as an atomic map from https://github.com/xaionaro-go/atomicmap .
func NewUkvsMapGeneratorAtomicMap() UkvsMapGenerator[_UkvsImplAtomicMap] {
	return _UkvsImplAtomicMapGenerator{}
}

// NewMap creates and returns a new _UkvsImplAtomicMap with pre-allocated
// 16 map slots.
func (_ _UkvsImplAtomicMapGenerator) NewMap() _UkvsImplAtomicMap {
	return _UkvsImplAtomicMap{atomicmap.NewWithArgs(16)}
}

// DestroyMap removes all values from given _UkvsImplAtomicMap, preserving
// allocated map slots, preparing it for being reused later.
func (_ _UkvsImplAtomicMapGenerator) DestroyMap(m _UkvsImplAtomicMap) {
	for _, key := range m.Keys() {
		_ = m.Unset(key)
	}
}

////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////

func (u _UkvsImplAtomicMap) Get(key any) any {
	return u.ie(u.Map.Get(ukvsPrepareUserKey(key)))
}

func (u _UkvsImplAtomicMap) Lookup(key any) (value any, found bool) {
	return u.he(u.Map.Get(ukvsPrepareUserKey(key)))
}

func (u _UkvsImplAtomicMap) Swap(key, value any) (prev any, was bool) {
	return u.he(u.Map.Swap(ukvsPrepareUserKey(key), value))
}

func (u _UkvsImplAtomicMap) Set(key, value any, overwrite bool) {
	key = ukvsPrepareUserKey(key)

	if overwrite {
		_ = u.Map.Set(key, value)
		return
	}

	if _, err := u.Map.Get(key); err == atomicmap.NotFound {
		_ = u.Map.Set(key, value)
	}
}

func (u _UkvsImplAtomicMap) Remove(key any) (prev any, was bool) {
	_ = u.Map.UnsetIf(ukvsPrepareUserKey(key), func(value any) bool {
		prev = value
		was = true
		return true
	})
	return
}

////////////////////////////////////////////////////////////////////////////////
///// PRIVATE METHODS //////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////

// he stands for "handle error", and returns given 'elem' as 1st arg and true
// when given 'err' is nil, false otherwise as 2nd arg.
func (_ _UkvsImplAtomicMap) he(elem any, err error) (any, bool) {
	return elem, err == nil
}

// ie stands for "ignore err", just returning given 'elem'.
func (_ _UkvsImplAtomicMap) ie(elem any, _ error) any { return elem }

*/

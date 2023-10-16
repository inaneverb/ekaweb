package ekaweb_private

import (
	"unsafe"

	"github.com/inaneverb/ekacore/ekaunsafe/v4"
)

type (
	// _UkvsImplSlice implements UkvsMap based on Golang slice.
	// It's the fastest implementation for small number of elems in UKVS.
	// NOT THREAD SAFETY! (But it's not a requirement, right?)
	_UkvsImplSlice []_UkvsImplSliceValue

	// _UkvsImplSliceValue represents one item of _UkvsImplSlice.
	// Key is xxhash of given key and Value is a stored value.
	_UkvsImplSliceValue = struct {
		Key   uint64
		Value any
	}

	// _UkvsImplAtomicMapGenerator implements UkvsMapGenerator,
	// working with _UkvsImplSlice.
	_UkvsImplSliceGenerator struct{}
)

////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////

// NewUkvsMapGeneratorSlice returns a new UkvsMapGenerator that operates
// with UkvsMap as a Golang's slice.
func NewUkvsMapGeneratorSlice() UkvsMapGenerator[*_UkvsImplSlice] {
	return _UkvsImplSliceGenerator{}
}

// NewMap creates and returns a new _UkvsImplSlice with pre-allocated
// 16 map slots.
func (_ _UkvsImplSliceGenerator) NewMap() *_UkvsImplSlice {
	var arr = make(_UkvsImplSlice, 0, 16)
	return &arr
}

// DestroyMap removes all values from given _UkvsImplSlice, preserving
// allocated map slots, preparing it for being reused later.
func (_ _UkvsImplSliceGenerator) DestroyMap(m *_UkvsImplSlice) {
	for i := range *m {
		(*m)[i].Value = nil
	}
	*m = (*m)[:0]
}

////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////

func (u *_UkvsImplSlice) Get(key any) any {
	var v, _ = u.Lookup(key)
	return v
}

func (u *_UkvsImplSlice) Lookup(key any) (value any, found bool) {

	var idx = u.getIndex(ukvsPrepareUserKey(key))
	if idx == -1 {
		return nil, false
	}

	var rt, wdRaw = u.getValueUnsafe(idx)
	return ekaunsafe.PackInterface(rt, unsafe.Add(nil, wdRaw)), true
}

func (u *_UkvsImplSlice) Swap(key, value any) (prev any, was bool) {

	var k = ukvsPrepareUserKey(key)
	var idx = u.getIndex(k)

	if was = idx != -1; was {
		var rt, wdRaw = u.getValueUnsafe(idx)
		prev = ekaunsafe.PackInterface(rt, unsafe.Add(nil, wdRaw))
		(*u)[idx].Value = value
	} else {
		*u = append(*u, _UkvsImplSliceValue{Key: k, Value: value})
	}

	return prev, was
}

func (u *_UkvsImplSlice) Set(key, value any, overwrite bool) {

	var k = ukvsPrepareUserKey(key)
	var idx = u.getIndex(k)

	switch {
	case idx != -1 && overwrite:
		(*u)[idx].Value = value // already exist this key

	case idx == -1:
		*u = append(*u, _UkvsImplSliceValue{Key: k, Value: value}) // not exist
	}
}

func (u *_UkvsImplSlice) Remove(key any) (prev any, was bool) {

	var idx = u.getIndex(ukvsPrepareUserKey(key))
	if was = idx != -1; was {
		var rt, wdRaw = u.getValueUnsafe(idx)
		prev = ekaunsafe.PackInterface(rt, unsafe.Add(nil, wdRaw))

		if n := len(*u); n == 1 {
			(*u)[idx].Value = nil // TODO: Do we really have to nil it?
			*u = (*u)[:0]         // only one elem total, just shrink to zero
		} else {
			(*u)[idx] = (*u)[n-1] // replace item with the last one
			(*u)[n-1].Value = nil // TODO: Do we really have to nil it?
			*u = (*u)[:n-1]       // shrink slice
		}
	}

	return prev, was
}

////////////////////////////////////////////////////////////////////////////////
///// PRIVATE METHODS //////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////

// getValueUnsafe returns stored value by given 'key' as its index in slice
// and Go interface's rtype and word (as uintptr).
//
// Using uintptr as a word (instead of unsafe.Pointer) is kinda trick to confuse
// escape analysis in a compiler to avoid unnecessary heap copy and additional
// allocation.
//
// NOTE! THIS METHOD EXPECTS (AND ACCEPTS) SLICE INDEX, NOT A "MAP" KEY!
func (u *_UkvsImplSlice) getValueUnsafe(idx int) (rt, wd uintptr) {

	if idx == -1 {
		return 0, 0
	}

	var wdPtr unsafe.Pointer
	rt, wdPtr = ekaunsafe.UnpackInterface((*u)[idx].Value).Tuple()
	return rt, uintptr(wdPtr)
}

// getIndex returns an index of element with the given 'key' from the slice.
// Returns -1 if no such element.
func (u *_UkvsImplSlice) getIndex(key uint64) int {
	for i, v := range *u {
		if v.Key == key {
			return i
		}
	}
	return -1
}

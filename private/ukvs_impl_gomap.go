package ekaweb_private

type (
	// _UkvsImplGoMap implements UkvsMap based on Golang map. Just a map.
	// NOT THREAD SAFETY! (But it's not a requirement, right?)
	_UkvsImplGoMap map[uint64]any

	// _UkvsImplGoMapGenerator implements UkvsMapGenerator,
	// working with _UkvsImplGoMap.
	_UkvsImplGoMapGenerator struct{}
)

////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////

// NewUkvsMapGeneratorGoMap returns a new UkvsMapGenerator that operates
// with UkvsMap as a Golang's maps.
func NewUkvsMapGeneratorGoMap() UkvsMapGenerator[_UkvsImplGoMap] {
	return _UkvsImplGoMapGenerator{}
}

// NewMap creates and returns a new _UkvsImplSlice with pre-allocated
// 16 map slots.
func (_ _UkvsImplGoMapGenerator) NewMap() _UkvsImplGoMap {
	return make(_UkvsImplGoMap, 16)
}

// DestroyMap removes all values from given _UkvsImplSlice, preserving
// allocated map slots, preparing it for being reused later.
func (_ _UkvsImplGoMapGenerator) DestroyMap(m _UkvsImplGoMap) {
	for k := range m {
		delete(m, k)
	}
}

////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////

func (u _UkvsImplGoMap) Get(key any) any {
	var k = ukvsPrepareUserKey(key)
	return u[k]
}

func (u _UkvsImplGoMap) Lookup(key any) (value any, found bool) {
	var k = ukvsPrepareUserKey(key)
	value, found = u[k]
	return value, found
}

func (u _UkvsImplGoMap) Swap(key, value any) (prev any, was bool) {
	var k = ukvsPrepareUserKey(key)

	prev, was = u[k]
	u[k] = value

	return prev, was
}

func (u _UkvsImplGoMap) Set(key, value any, overwrite bool) {
	var k = ukvsPrepareUserKey(key)

	if !overwrite {
		if _, found := u[k]; found {
			return
		}
	}

	u[k] = value
}

func (u _UkvsImplGoMap) Remove(key any) (prev any, was bool) {
	var k = ukvsPrepareUserKey(key)

	prev, was = u[k]
	if was {
		delete(u, k)
	}

	return prev, was
}

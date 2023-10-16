package ekaweb_private

import (
	"encoding/binary"

	"github.com/cespare/xxhash/v2"
	"github.com/inaneverb/ekacore/ekaunsafe/v4"
)

type (
	// UkvsMap is a core of UKVS - a storage itself.
	// It's an interface, so there's multiple implementations.
	UkvsMap interface {
		Get(key any) any
		Lookup(key any) (value any, found bool)
		Swap(key, value any) (prev any, was bool)
		Set(key, value any, overwrite bool)
		Remove(key any) (prev any, was bool)
	}

	// UkvsMapGenerator is a UkvsMap getter/destructor.
	UkvsMapGenerator[M UkvsMap] interface {
		NewMap() M
		DestroyMap(m M)
	}
)

////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////

// ukvsPrepareUserKey returns hash sum for given 'key'.
func ukvsPrepareUserKey(key any) uint64 {
	var i = ekaunsafe.UnpackInterface(key)

	switch i.Type {
	case ekaunsafe.RTypeString():
		return xxhash.Sum64String(*(*string)(i.Word))

	case ekaunsafe.RTypeInt():
		return ukvsPrepareUserKeyInt(uint64(*(*int)(i.Word)))

	case ekaunsafe.RTypeInt8():
		return ukvsPrepareUserKeyInt(uint64(*(*int8)(i.Word)))

	case ekaunsafe.RTypeInt16():
		return ukvsPrepareUserKeyInt(uint64(*(*int16)(i.Word)))

	case ekaunsafe.RTypeInt32():
		return ukvsPrepareUserKeyInt(uint64(*(*int32)(i.Word)))

	case ekaunsafe.RTypeInt64():
		return ukvsPrepareUserKeyInt(uint64(*(*int64)(i.Word)))

	case ekaunsafe.RTypeUint():
		return ukvsPrepareUserKeyInt(uint64(*(*uint)(i.Word)))

	case ekaunsafe.RTypeUint8():
		return ukvsPrepareUserKeyInt(uint64(*(*uint8)(i.Word)))

	case ekaunsafe.RTypeUint16():
		return ukvsPrepareUserKeyInt(uint64(*(*uint16)(i.Word)))

	case ekaunsafe.RTypeUint32():
		return ukvsPrepareUserKeyInt(uint64(*(*uint32)(i.Word)))

	case ekaunsafe.RTypeUint64():
		return ukvsPrepareUserKeyInt(*(*uint64)(i.Word))

	default:
		return uint64(i.Type)
	}
}

// ukvsPrepareUserKeyInt returns a hash sum for given 'v'.
func ukvsPrepareUserKeyInt(v uint64) uint64 {
	var d [8]byte
	binary.BigEndian.PutUint64(d[:], v)
	return xxhash.Sum64(d[:])
}

package ekaweb_respondent

import (
	"reflect"
)

func isGoHashableObject(kind reflect.Kind) bool {
	switch kind {
	case reflect.Array, reflect.Slice, reflect.Map, reflect.Chan, reflect.Func:
		return false
	default:
		return true
	}
}

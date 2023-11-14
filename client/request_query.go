package ekaweb_client

import (
	"fmt"
	"net/url"
	"reflect"
	"strings"

	"github.com/inaneverb/ekaweb/v2/private"
)

// TODO: This is not optimized version of URI query encoding.
//  You should reject usage of Golang reflect package and rewrite all this code.

type requestQuery struct {
	source any
}

func RequestQuery(source any) ekaweb_private.ClientRequest {
	return &requestQuery{source}
}

func (w *requestQuery) Data() ([]byte, error) {

	var v = reflect.ValueOf(w.source)

	switch {
	case !v.IsValid():
		return nil, fmt.Errorf("invalid source (%s)", v.String())

	case v.Kind() == reflect.Pointer && v.IsNil():
		return nil, fmt.Errorf("nil source (type: %s)", v.Type().String())

	case v.Kind() == reflect.Pointer:
		v = v.Elem()
	}

	switch {
	case v.Kind() == reflect.Struct:
		return w.encodeStruct(v)

	case v.Kind() == reflect.Map && v.Type().Key().Kind() != reflect.String:
		const E = "unsupported map key type (%s); only string is allowed"
		return nil, fmt.Errorf(E, v.Type().Key().String())

	case v.Kind() == reflect.Map:
		return w.encodeMap(v)

	default:
		return nil, fmt.Errorf("unsupported type (%s)", v.Type().String())
	}
}

func (w *requestQuery) ContentType() string {
	return ""
}

////////////////////////////////////////////////////////////////////////////////

func (w *requestQuery) encodeStruct(v reflect.Value) ([]byte, error) {

	const TagName1 = "uri"
	const TagName2 = "form"

	var buf = make([]byte, 0, 64)

	var vf = reflect.VisibleFields(v.Type())
	for i, n := 0, len(vf); i < n; i++ {

		switch {
		case !vf[i].IsExported():
			continue

		case !w.isAtom(vf[i].Type.Kind()):
			const E = "non atomic nested field (offset: %d, name: %s)"
			return nil, fmt.Errorf(E, i, vf[i].Name)
		}

		var name string
		if name = vf[i].Tag.Get(TagName1); name == "" {
			name = vf[i].Tag.Get(TagName2)
		}

		var omitempty bool

		if omitempty = strings.HasSuffix(name, ",omitempty"); omitempty {
			name = name[:len(name)-10]
		}

		switch {
		case name == "" && vf[i].Anonymous:
			continue
		case name == "":
			name = strings.ToLower(vf[i].Name)
		}

		buf = w.appendAtom(buf, name, v.Field(i), omitempty)
	}

	return w.cleanUpBuffer(buf)
}

func (w *requestQuery) encodeMap(v reflect.Value) ([]byte, error) {

	var buf = make([]byte, 0, 64)

	for iter := v.MapRange(); iter.Next(); {
		var name = iter.Key().String()
		var v = iter.Value()

		switch {
		case name == "":
			continue

		case w.isAtom(v.Kind()):
			return nil, fmt.Errorf("non atomic map value (name: %s)", name)
		}

		buf = w.appendAtom(buf, name, v, false)
	}

	return w.cleanUpBuffer(buf)
}

func (w *requestQuery) appendAtom(
	to []byte, name string, v reflect.Value, omitempty bool) []byte {

	if v.IsZero() && omitempty {
		return to
	}

	to = append(to, name...)
	to = append(to, '=')
	to = append(to, url.QueryEscape(fmt.Sprintf("%v", v.Interface()))...)
	to = append(to, '&')

	return to
}

func (w *requestQuery) cleanUpBuffer(buf []byte) ([]byte, error) {
	if buf[len(buf)-1] == '&' {
		buf = buf[:len(buf)-1]
	}
	if len(buf) == 0 {
		buf = nil
	}
	return buf, nil
}

func (_ *requestQuery) isAtom(k reflect.Kind) bool {
	return k == reflect.Bool || k == reflect.Int || k == reflect.Uint ||
		k == reflect.Int8 || k == reflect.Int16 ||
		k == reflect.Int32 || k == reflect.Int64 ||
		k == reflect.Uint8 || k == reflect.Uint16 ||
		k == reflect.Uint32 || k == reflect.Uint64 ||
		k == reflect.Float32 || k == reflect.Float64 ||
		k == reflect.String
}

package ekaweb_bind

import (
	"errors"
	"io"
	"net/http"

	"github.com/inaneverb/ekaweb"
)

//goland:noinspection GoErrorStringFormat
var (
	ErrBodyEmpty    = errors.New("Data: Empty body")
	ErrBodyTooLarge = errors.New("Data: Body too large")
)

// ReadBody reads the payload from the HTTP request if it's not bigger
// than provided `maxMemory`. The method do not return an error, because if so,
// it's already saved to the HTTP request context. Just check the returned data:
// if it's nil - an error is occurred (and already saved).
//
// Rules:
//   - You don't have to worry about allocations, check max limit, etc.
//     Just provide max allowed memory, and all of that will be encapsulated.
//
//   - Only one call. You even don't have to check errors. All of them already saved,
//     using ErrorApply() and ErrorDetailApply(). Just check returned arguments.
//     If it's nil, an error occurred, and it's already saved to user context.
//
// Returned errors:
// - ErrBodyEmpty: Empty body.
// - ErrBodyTooLarge: Body size exceeded the maximum limit.
func ReadBody(r *http.Request, maxMemory int64) []byte {

	// TODO: Add bytes buffer pool, because they should be re-used
	//  instead of allocated each time.

	var buf = make([]byte, maxMemory+1)
	var n, err = r.Body.Read(buf)
	switch {
	case err != nil && err != io.EOF:
		ekaweb.ErrorApply(r, err)
		return nil

	case n == 0:
		ekaweb.ErrorApply(r, ErrBodyEmpty)
		return nil

	case int64(n) == maxMemory+1:
		ekaweb.ErrorApply(r, ErrBodyTooLarge)
		return nil
	}

	return buf[:n]
}

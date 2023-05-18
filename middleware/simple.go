package ekaweb_middleware

import (
	"github.com/inaneverb/ekaweb/private"
)

// AbortWith generates a middleware that is always returns passed error
// failing the whole process of handling HTTP request.
func AbortWith(err error) ekaweb_private.Middleware {
	return ekaweb_private.AbortWith(err)
}

// Recover generates a middleware that prevents panics from next HTTP controllers.
// It wraps these calls by the deferring recover() call and if panic is recovered,
// it transforms it to the error and returns it.
func Recover() ekaweb_private.Middleware {
	return ekaweb_private.Recover()
}

package ekaweb_noop

import (
	"github.com/inaneverb/ekaweb/v2"
	"github.com/inaneverb/ekaweb/v2/private"
)

type nopeRouterSimple struct{}

func (r *nopeRouterSimple) Reg(_ string, _ ...any) ekaweb.RouterSimple {
	return r
}

func (r *nopeRouterSimple) Build() ekaweb_private.Handler {
	return ekaweb_private.NewEmptyHandler()
}

// NewRouterSimple returns an ekaweb.RouterSimple that does nothing
// and returns an empty handler (no-op) as build result.
func NewRouterSimple(_ ...ekaweb.RouterOption) ekaweb.RouterSimple {
	return (*nopeRouterSimple)(nil)
}

var _ ekaweb.RouterSimple = (*nopeRouterSimple)(nil)

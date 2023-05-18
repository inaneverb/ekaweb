package ekaweb_nope

import (
	"github.com/inaneverb/ekaweb"
	"github.com/inaneverb/ekaweb/private"
)

type nopeRouter struct{}

func (r *nopeRouter) Use(_ ...any) ekaweb.Router               { return r }
func (r *nopeRouter) Group(_ string, _ ...any) ekaweb.Router   { return r }
func (r *nopeRouter) Get(_ string, _ ...any) ekaweb.Router     { return r }
func (r *nopeRouter) Head(_ string, _ ...any) ekaweb.Router    { return r }
func (r *nopeRouter) Post(_ string, _ ...any) ekaweb.Router    { return r }
func (r *nopeRouter) Put(_ string, _ ...any) ekaweb.Router     { return r }
func (r *nopeRouter) Delete(_ string, _ ...any) ekaweb.Router  { return r }
func (r *nopeRouter) Connect(_ string, _ ...any) ekaweb.Router { return r }
func (r *nopeRouter) Options(_ string, _ ...any) ekaweb.Router { return r }
func (r *nopeRouter) Trace(_ string, _ ...any) ekaweb.Router   { return r }
func (r *nopeRouter) Patch(_ string, _ ...any) ekaweb.Router   { return r }
func (r *nopeRouter) NotFound(_ any) ekaweb.Router             { return r }
func (r *nopeRouter) MethodNotAllowed(_ any) ekaweb.Router     { return r }

func (r *nopeRouter) Build() ekaweb.Handler {
	return ekaweb_private.NewEmptyHandler()
}

// NewRouter returns an ekaweb.Router that does nothing
// and returns an empty handler (no-op) as build result.
func NewRouter(_ ...ekaweb.RouterOption) ekaweb.Router {
	return (*nopeRouter)(nil)
}

var _ ekaweb.Router = (*nopeRouter)(nil)

package ekaweb_noop

import (
	"github.com/inaneverb/ekaweb/v2"
)

type nopeServer struct{}

func (_ *nopeServer) AsyncStart() error { return nil }
func (_ *nopeServer) Stop() error       { return nil }

// NewServer returns an ekaweb.Server that does nothing at all.
// Its methods always returns nil as error.
func NewServer(_ ...ekaweb.ServerOption) ekaweb.Server {
	return (*nopeServer)(nil)
}

var _ ekaweb.Server = (*nopeServer)(nil)

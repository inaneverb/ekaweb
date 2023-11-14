package ekaweb_nbio

import (
	"strings"

	"github.com/inaneverb/ekacore/ekaunsafe/v4"
	"github.com/inaneverb/ekaweb/v2"
	"github.com/inaneverb/ekaweb/v2/private"

	"github.com/lesismal/nbio/logging"
	"github.com/lesismal/nbio/nbhttp"
)

type Server struct {
	origin *nbhttp.Server
}

////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////

func (s *Server) AsyncStart() error {
	return s.origin.Start()
}

func (s *Server) Stop() error {
	s.origin.Stop()
	return nil
}

////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////

func NewServer(options ...ekaweb.ServerOption) ekaweb.Server {

	var nbioConfig nbhttp.Config

	for i, n := 0, len(options); i < n; i++ {
		if ekaunsafe.UnpackInterface(options[i]).Word == nil {
			continue
		}

		switch option := options[i].(type) {
		case *ekaweb_private.ServerOptionHandler:
			if ekaunsafe.UnpackInterface(option.Handler).Word != nil {
				nbioConfig.Handler = option.Handler
			}

		case *ekaweb_private.ServerOptionListenAddr:
			if option.Addr = strings.TrimSpace(option.Addr); option.Addr != "" {
				nbioConfig.Addrs = append(nbioConfig.Addrs, option.Addr)
			}

		case *ekaweb_private.ClientServerOptionLogger:
			if option.Log != nil {
				logging.SetLogger(newLoggingBridge(option.Log))
			}

		case *ekaweb_private.ClientServerOptionTimeout:
			// NOT SUPPORTED
		}
	}

	return &Server{nbhttp.NewServer(nbioConfig)}
}

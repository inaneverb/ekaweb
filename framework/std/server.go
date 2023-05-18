package ekaweb_std

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/inaneverb/ekacore/ekaunsafe/v4"
	"github.com/inaneverb/ekaweb"
	"github.com/inaneverb/ekaweb/private"
)

type Server struct {
	origin *http.Server
}

////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////

func (s *Server) AsyncStart() error {
	go func() { _ = s.origin.ListenAndServe() }()
	return nil
}

func (s *Server) Stop() error {
	var ctx, cancelFunc = context.WithTimeout(context.TODO(), 5*time.Second)
	defer cancelFunc()
	return s.origin.Shutdown(ctx)
}

////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////

func NewServer(options ...ekaweb.ServerOption) ekaweb.Server {

	var server Server
	server.origin = new(http.Server)

	server.origin.IdleTimeout = 10 * time.Second
	server.origin.SetKeepAlivesEnabled(true)

	for i, n := 0, len(options); i < n; i++ {
		if ekaunsafe.UnpackInterface(options[i]).Word == nil {
			continue
		}

		switch option := options[i].(type) {

		case *ekaweb_private.ServerOptionKeepAlive:
			server.origin.SetKeepAlivesEnabled(option.Enabled)

			if option.Duration >= 1*time.Second {
				server.origin.IdleTimeout = option.Duration
			}

		case *ekaweb_private.ServerOptionHandler:
			if ekaunsafe.UnpackInterface(option.Handler).Word != nil {
				server.origin.Handler = option.Handler
			}

		case *ekaweb_private.ServerOptionListenAddr:
			if option.Addr = strings.TrimSpace(option.Addr); option.Addr != "" {
				server.origin.Addr = option.Addr
			}

		case *ekaweb_private.ClientServerOptionLogger:
			// UNSUPPORTED

		case *ekaweb_private.ClientServerOptionTimeout:
			if option.ReadTimeout > 0 {
				server.origin.ReadTimeout = option.ReadTimeout
			}
			if option.WriteTimeout > 0 {
				server.origin.WriteTimeout = option.WriteTimeout
			}
		}
	}

	return &server
}

package ekaweb_fasthttp

import (
	"strings"
	"time"

	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttpadaptor"

	"github.com/inaneverb/ekacore/ekaunsafe/v4"
	"github.com/inaneverb/ekaweb"
	"github.com/inaneverb/ekaweb/private"
)

type Server struct {
	origin *fasthttp.Server
	addr   string
}

////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////

func (s *Server) AsyncStart() error {
	go func() { _ = s.origin.ListenAndServe(s.addr) }()
	return nil
}

func (s *Server) Stop() error {
	return s.origin.Shutdown()
}

////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////

func NewServer(options ...ekaweb.ServerOption) ekaweb.Server {

	var server Server
	server.origin = new(fasthttp.Server)

	server.origin.IdleTimeout = 10 * time.Second
	server.origin.TCPKeepalive = true

	for i, n := 0, len(options); i < n; i++ {
		if ekaunsafe.UnpackInterface(options[i]).Word == nil {
			continue
		}

		switch option := options[i].(type) {

		case *ekaweb_private.ServerOptionKeepAlive:
			server.origin.DisableKeepalive = !option.Enabled
			server.origin.TCPKeepalive = option.Enabled

			if option.Duration >= 1*time.Second {
				server.origin.IdleTimeout = option.Duration
			}

		case *ekaweb_private.ServerOptionHandler:
			if ekaunsafe.UnpackInterface(option.Handler).Word != nil {
				server.origin.Handler =
					fasthttpadaptor.NewFastHTTPHandler(option.Handler)

			}

		case *ekaweb_private.ServerOptionListenAddr:
			if option.Addr = strings.TrimSpace(option.Addr); option.Addr != "" {
				server.addr = option.Addr
			}

		case *ekaweb_private.ClientServerOptionLogger:
			if option.Log != nil {
				server.origin.Logger = newLoggingBridge(option.Log)
			}

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

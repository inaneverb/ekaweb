package ekaweb_private

import (
	"net/http"
	"time"
)

type ClientOption interface {
	Name() string
	noOneCanImplementClientOptionInterface()
}

type ClientOptionHostAddr struct {
	Addr string
}

type ClientOptionUserAgent struct {
	UserAgent string
}

func (o *ClientOptionHostAddr) Name() string {
	return "WithHostAddr"
}

func (o *ClientOptionUserAgent) Name() string {
	return "WithUserAgent"
}

func (o *ClientOptionHostAddr) noOneCanImplementClientOptionInterface()  {}
func (o *ClientOptionUserAgent) noOneCanImplementClientOptionInterface() {}

////////////////////////////////////////////////////////////////////////////////

type ServerOption interface {
	Name() string
	noOneCanImplementServerOptionInterface()
}

type ServerOptionTransport struct {
	Transport http.RoundTripper
}

type ServerOptionHandler struct {
	Handler Handler
}

type ServerOptionListenAddr struct {
	Addr string
}

type ServerOptionKeepAlive struct {
	Enabled  bool
	Duration time.Duration
}

func (o *ServerOptionTransport) Name() string {
	return "WithTransport"
}

func (o *ServerOptionHandler) Name() string {
	return "WithHandler"
}

func (o *ServerOptionListenAddr) Name() string {
	return "WithListenAddr"
}

func (o *ServerOptionKeepAlive) Name() string {
	return "WithKeepAlive"
}

func (o *ServerOptionTransport) noOneCanImplementServerOptionInterface()  {}
func (o *ServerOptionHandler) noOneCanImplementServerOptionInterface()    {}
func (o *ServerOptionListenAddr) noOneCanImplementServerOptionInterface() {}
func (o *ServerOptionKeepAlive) noOneCanImplementServerOptionInterface()  {}

////////////////////////////////////////////////////////////////////////////////

type ClientServerOption interface {
	Name() string
	noOneCanImplementClientOptionInterface()
	noOneCanImplementServerOptionInterface()
}

type ClientServerOptionTimeout struct {
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

type ClientServerOptionLogger struct {
	Log Logger
}

func (o *ClientServerOptionTimeout) Name() string {
	return "WithTimeouts"
}

func (o *ClientServerOptionLogger) Name() string {
	return "WithLogger"
}

func (o *ClientServerOptionTimeout) noOneCanImplementClientOptionInterface() {}
func (o *ClientServerOptionTimeout) noOneCanImplementServerOptionInterface() {}
func (o *ClientServerOptionLogger) noOneCanImplementClientOptionInterface()  {}
func (o *ClientServerOptionLogger) noOneCanImplementServerOptionInterface()  {}

package ekaweb

import (
	"time"

	"github.com/inaneverb/ekaweb/private"
)

////////////////////////////////////////////////////////////////////////////////
///// ROUTER OPTIONS ///////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////

func WithCoreInit(do bool) RouterOption {
	return &ekaweb_private.RouterOptionCoreInit{Enable: do}
}

func WithErrorHandler(cb ekaweb_private.ErrorHandlerHTTP) RouterOption {
	return &ekaweb_private.RouterOptionErrorHandler{Handler: cb}
}

func WithCustomJSON(enc ekaweb_private.MarshalCallback, dec ekaweb_private.UnmarshalCallback) RouterOption {
	return &ekaweb_private.RouterOptionCustomJSON{Encoder: enc, Decoder: dec}
}

func WithServerName(name string) RouterOption {
	return &ekaweb_private.RouterOptionServerName{ServerName: name}
}

////////////////////////////////////////////////////////////////////////////////
///// CLIENT & SERVER OPTIONS //////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////

func WithLogger(log ekaweb_private.Logger) ClientServerOption {
	return &ekaweb_private.ClientServerOptionLogger{Log: log}
}

func WithTimeouts(readTimeout, writeTimeout time.Duration) ClientServerOption {
	return &ekaweb_private.ClientServerOptionTimeout{
		ReadTimeout: readTimeout, WriteTimeout: writeTimeout,
	}
}

////////////////////////////////////////////////////////////////////////////////
///// SERVER OPTIONS ///////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////

func WithHandler(handler ekaweb_private.Handler) ServerOption {
	return &ekaweb_private.ServerOptionHandler{Handler: handler}
}

func WithListenAddr(addr string) ServerOption {
	return &ekaweb_private.ServerOptionListenAddr{Addr: addr}
}

func WithKeepAlive(enabled bool, duration time.Duration) ServerOption {
	return &ekaweb_private.ServerOptionKeepAlive{Enabled: enabled, Duration: duration}
}

////////////////////////////////////////////////////////////////////////////////
///// CLIENT OPTIONS ///////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////

func WithHostAddr(addr string) ClientOption {
	return &ekaweb_private.ClientOptionHostAddr{Addr: addr}
}

func WithUserAgent(userAgent string) ClientOption {
	return &ekaweb_private.ClientOptionUserAgent{UserAgent: userAgent}
}

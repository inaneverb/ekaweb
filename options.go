package ekaweb

import (
	"encoding/json"
	"io"
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

// WithCodec returns an Option, that overwrites default codec (encoder or/and
// decoder), which are defaulted to "encoding/json" Golang package.
func WithCodec[E ekaweb_private.Encoder, D ekaweb_private.Decoder](
	encGetter func(w io.Writer) E, decGetter func(r io.Reader) D) RouterOption {

	var encGetterTransformed = wrapEncGetter(json.NewEncoder)
	if encGetter != nil {
		encGetterTransformed = wrapEncGetter(encGetter)
	}

	var decGetterTransformed = wrapDecGetter(json.NewDecoder)
	if decGetter != nil {
		decGetterTransformed = wrapDecGetter(decGetter)
	}

	return &ekaweb_private.RouterOptionCodec{
		EncoderGetter: encGetterTransformed,
		DecoderGetter: decGetterTransformed,
	}
}

// wrapEncGetter returns _EncoderGetter from its generic variant.
func wrapEncGetter[E ekaweb_private.Encoder](
	encGetter func(w io.Writer) E) ekaweb_private.EncoderGetter {

	return func(w io.Writer) ekaweb_private.Encoder { return encGetter(w) }
}

// wrapDecGetter returns _DecoderGetter from its generic variant.
func wrapDecGetter[D ekaweb_private.Decoder](
	decGetter func(r io.Reader) D) ekaweb_private.DecoderGetter {

	return func(r io.Reader) ekaweb_private.Decoder { return decGetter(r) }
}

func WithServerName(name string) RouterOption {
	return &ekaweb_private.RouterOptionServerName{ServerName: name}
}

func WithTrailingSlash(redirect, strip bool) RouterOption {
	return &ekaweb_private.RouterOptionTrailingSlash{Redirect: redirect, Strip: strip}
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

package ekaweb_private

import (
	"io"
)

type RouterSimple interface {
	Reg(path string, middlewaresAndHandler ...any) RouterSimple
	Build() Handler
}

type Router interface {
	Use(middlewares ...any) Router
	Group(prefix string, middlewares ...any) Router

	Get(path string, middlewaresAndHandler ...any) Router
	Head(path string, middlewaresAndHandler ...any) Router
	Post(path string, middlewaresAndHandler ...any) Router
	Put(path string, middlewaresAndHandler ...any) Router
	Delete(path string, middlewaresAndHandler ...any) Router
	Connect(path string, middlewaresAndHandler ...any) Router
	Options(path string, middlewaresAndHandler ...any) Router
	Trace(path string, middlewaresAndHandler ...any) Router
	Patch(path string, middlewaresAndHandler ...any) Router

	NotFound(handler any) Router
	MethodNotAllowed(handler any) Router

	Build() Handler
}

type RouterOption interface {
	Name() string
	noOneCanImplementRouterOptionInterface()
}

type (
	// Encoder shall encode by some codec given object writing it to himself.
	// Usually, the type that implements it, writes output to io.Writer.
	Encoder interface {
		Encode(e any) error
	}

	// Decoder shall decode data by some codec from itself to given object.
	// Usually, the type that implements it, reads data from io.Reader.
	Decoder interface {
		Decode(to any) error
	}

	EncoderGetter = func(w io.Writer) Encoder // generic-less aliases
	DecoderGetter = func(r io.Reader) Decoder // generic-less aliases
)

////////////////////////////////////////////////////////////////////////////////
///// Option argument holders //////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////

type (
	RouterOptionCoreInit struct {
		Enable bool
	}

	RouterOptionErrorHandler struct {
		Handler ErrorHandlerHTTP
	}

	RouterOptionCodec struct {
		EncoderGetter EncoderGetter
		DecoderGetter DecoderGetter
	}

	RouterOptionServerName struct {
		ServerName string
	}

	RouterOptionTrailingSlash struct {
		Redirect bool
		Strip    bool
	}
)

////////////////////////////////////////////////////////////////////////////////
///// Option interface implementations /////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////

func (o *RouterOptionCoreInit) Name() string {
	return "WithCoreInit"
}

func (o *RouterOptionErrorHandler) Name() string {
	return "WithErrorHandler"
}

func (o *RouterOptionCodec) Name() string {
	return "WithCodec"
}

func (o *RouterOptionServerName) Name() string {
	return "WithServerName"
}

func (o *RouterOptionTrailingSlash) Name() string {
	return "WithTrailingSlash"
}

func (o *RouterOptionErrorHandler) noOneCanImplementRouterOptionInterface()  {}
func (o *RouterOptionCodec) noOneCanImplementRouterOptionInterface()         {}
func (o *RouterOptionServerName) noOneCanImplementRouterOptionInterface()    {}
func (o *RouterOptionCoreInit) noOneCanImplementRouterOptionInterface()      {}
func (o *RouterOptionTrailingSlash) noOneCanImplementRouterOptionInterface() {}

////////////////////////////////////////////////////////////////////////////////
///// PRIVATE FUNCTIONS ////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////

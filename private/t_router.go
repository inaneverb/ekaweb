package ekaweb_private

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
	MarshalCallback   = func(v any) ([]byte, error)
	UnmarshalCallback = func(data []byte, v any) error
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

	RouterOptionCustomJSON struct {
		Encoder MarshalCallback
		Decoder UnmarshalCallback
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

func (o *RouterOptionCustomJSON) Name() string {
	return "WithCustomJSON"
}

func (o *RouterOptionServerName) Name() string {
	return "WithServerName"
}

func (o *RouterOptionTrailingSlash) Name() string {
	return "WithTrailingSlash"
}

func (o *RouterOptionErrorHandler) noOneCanImplementRouterOptionInterface()  {}
func (o *RouterOptionCustomJSON) noOneCanImplementRouterOptionInterface()    {}
func (o *RouterOptionServerName) noOneCanImplementRouterOptionInterface()    {}
func (o *RouterOptionCoreInit) noOneCanImplementRouterOptionInterface()      {}
func (o *RouterOptionTrailingSlash) noOneCanImplementRouterOptionInterface() {}

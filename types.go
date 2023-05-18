package ekaweb

import (
	"github.com/inaneverb/ekaweb/private"
)

type Handler = ekaweb_private.Handler
type HandlerExtended = ekaweb_private.HandlerExtended

type HandlerFunc = ekaweb_private.HandlerFunc
type HandlerFuncNoErrorCheck = ekaweb_private.HandlerFuncNoErrorCheck

type Middleware = ekaweb_private.Middleware
type MiddlewareExtended = ekaweb_private.MiddlewareExtended

type MiddlewareFunc = ekaweb_private.MiddlewareFunc
type MiddlewareFuncNoErrorCheck = ekaweb_private.MiddlewareFuncNoErrorCheck

type Logger = ekaweb_private.Logger
type Client = ekaweb_private.Client
type Router = ekaweb_private.Router
type RouterSimple = ekaweb_private.RouterSimple
type Server = ekaweb_private.Server

type ClientOption = ekaweb_private.ClientOption
type RouterOption = ekaweb_private.RouterOption
type ServerOption = ekaweb_private.ServerOption
type ClientServerOption = ekaweb_private.ClientServerOption

type ErrorHandler = ekaweb_private.ErrorHandler
type ErrorHandlerHTTP = ekaweb_private.ErrorHandlerHTTP

type ClientRequest = ekaweb_private.ClientRequest
type ClientResponse = ekaweb_private.ClientResponse

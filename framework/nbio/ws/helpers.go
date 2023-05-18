package ekaweb_nbiows

import (
	"context"

	"github.com/inaneverb/ekaweb/private"
	"github.com/inaneverb/ekaweb/websocket"
)

// These aliases are exist just because I hate named imports
// but nbio WebSocket package and internal/extension one has the same names.

type Handler = ekaweb_socket.Handler
type Option = ekaweb_socket.Option
type Options = ekaweb_socket.Options
type MessageType = ekaweb_socket.MessageType
type CloseCode = ekaweb_socket.CloseCode

const (
	MessageTypeDataText       = ekaweb_socket.MessageTypeDataText
	MessageTypeDataBinary     = ekaweb_socket.MessageTypeDataBinary
	MessageTypeControlClose   = ekaweb_socket.MessageTypeControlClose
	MessageTypeControlPing    = ekaweb_socket.MessageTypeControlPing
	MessageTypeControlPong    = ekaweb_socket.MessageTypeControlPong
	MessageTypeControlInvalid = ekaweb_socket.MessageTypeInvalid
)

const (
	CloseCodeNoStatusReceived = ekaweb_socket.CloseCodeNoStatusReceived
	CloseCodeInternalError    = ekaweb_socket.CloseCodeInternalError
	CloseCodeAbnormal         = ekaweb_socket.CloseCodeAbnormal
)

var (
	ErrHandshakeBadMethod           = ekaweb_socket.ErrHandshakeBadMethod
	ErrHandshakeBadHeaderUpgrade    = ekaweb_socket.ErrHandshakeBadHeaderUpgrade
	ErrHandshakeBadHeaderSecKey     = ekaweb_socket.ErrHandshakeBadHeaderSecKey
	ErrHandshakeBadHeaderSecVersion = ekaweb_socket.ErrHandshakeBadHeaderSecVersion
	ErrHandshakeBadOrigin           = ekaweb_socket.ErrHandshakeBadOrigin
)

func makeOptions(options []Option) *Options {
	return ekaweb_socket.PrepareOptions(options)
}

func applyErrorHandler(conn *Conn, err error) {

	ekaweb_private.UkvsInsertUserError(conn.ctx, err)

	if conn.options.ErrorHandler != nil {
		conn.options.ErrorHandler(ekaweb_socket.Conn(conn), err)
	} else {
		conn.CloseWithCode(CloseCodeAbnormal)
	}
}

func duplicateHTTPRequestContext(ctx context.Context) (context.Context, func()) {
	return context.WithCancel(ekaweb_private.UkvsPropagate(ctx, context.Background()))
}

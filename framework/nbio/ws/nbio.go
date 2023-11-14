package ekaweb_nbiows

import (
	"context"
	"encoding/binary"
	"errors"
	"net/http"

	"github.com/lesismal/nbio/mempool"
	"github.com/lesismal/nbio/nbhttp/websocket"

	"github.com/inaneverb/ekaweb/v2"
	"github.com/inaneverb/ekaweb/v2/private"
)

type _NbioWebSocketOnOpenCallback = func(conn *websocket.Conn)

type _NbioWebSocketOnCloseCallback = func(conn *websocket.Conn, err error)

type _NbioWebSocketOnMessageCallback = func(conn *websocket.Conn, typ websocket.MessageType, data []byte)

////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////

func NewHandler(handler Handler, options ...Option) ekaweb.Handler {
	optionsSet := makeOptions(options)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx, cancelFunc := duplicateHTTPRequestContext(r.Context())

		var customUpgradeHeaders = w.Header()
		var optionsHeadersPresent = len(optionsSet.ResponseHeaders) > 0
		var customUpgradeHeadersPresent = len(customUpgradeHeaders) > 0

		switch {

		case optionsHeadersPresent && customUpgradeHeadersPresent:
			h1 := optionsSet.ResponseHeaders.Clone()
			customUpgradeHeaders = ekaweb.HeadersMerge(customUpgradeHeaders, h1, true)

		case optionsHeadersPresent:
			customUpgradeHeaders = optionsSet.ResponseHeaders
		}

		upgrader := websocket.NewUpgrader()
		upgrader.OnOpen(makeOnOpenCallback(ctx, handler, optionsSet))
		upgrader.OnClose(makeOnCloseCallback(handler, cancelFunc))
		upgrader.OnMessage(makeOnMessageCallback(handler))

		upgrader.SetCloseHandler(nbioCloseMessageHandler)

		if optionsSet.CheckOrigin != nil {
			upgrader.CheckOrigin = optionsSet.CheckOrigin
		}

		if _, err := upgrader.Upgrade(w, r, customUpgradeHeaders); err != nil {
			if err = mapNbioUpgradeError(err); err != nil {
				ekaweb_private.UkvsInsertUserError(r.Context(), err)
			}
		}

		ekaweb_private.UkvsMarkConnectionAsHijacked(r.Context())
	})
}

////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////

func makeOnOpenCallback(
	ctx context.Context, handler Handler, options *Options,
) _NbioWebSocketOnOpenCallback {

	return _NbioWebSocketOnOpenCallback(func(originConn *websocket.Conn) {
		conn := makeConn(ctx, options, originConn)

		if err := handler.OnOpen(conn); err != nil {
			applyErrorHandler(conn, err)
		}
	})
}

func makeOnCloseCallback(
	handler Handler, ctxCancelFunc func()) _NbioWebSocketOnCloseCallback {

	return func(originConn *websocket.Conn, err error) {

		defer ctxCancelFunc()

		var wrappedConn = connFromOrigin(originConn)
		var cc, ccDetail = wrappedConn.getCloseData()

		// The response close message is sent already. So it doesn't matter
		// whether user return true or false from OnClose.

		_ = handler.OnClose(wrappedConn, cc, ccDetail)
	}
}

func makeOnMessageCallback(
	handler Handler) _NbioWebSocketOnMessageCallback {

	return func(conn *websocket.Conn, typ websocket.MessageType, data []byte) {
		var wrappedConn = connFromOrigin(conn)
		if wrappedConn == nil {
			return
		}
		var err = handler.OnMessage(wrappedConn, MessageTypeFromNbio(typ), data)
		if err != nil {
			applyErrorHandler(wrappedConn, err)
		}
	}
}

////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////

func nbioCloseMessageHandler(conn *websocket.Conn, code int, detail string) {
	processCloseMessage(conn, uint16(code), detail)
}

func processCloseMessage(conn *websocket.Conn, code uint16, detail string) {
	const MaxControlFramePayloadSize = 125

	if len(detail)+2 > MaxControlFramePayloadSize {
		detail = detail[:MaxControlFramePayloadSize-2]
	}

	buf := mempool.Malloc(len(detail) + 2)
	binary.BigEndian.PutUint16(buf[:2], code)
	copy(buf[2:], detail)
	_ = conn.WriteMessage(MessageTypeToNbio(MessageTypeControlClose), buf)
	mempool.Free(buf)

	connFromOrigin(conn).setCloseData(CloseCode(code), detail)
}

func mapNbioUpgradeError(nbioErr error) error {
	switch {
	case errors.Is(nbioErr, websocket.ErrUpgradeTokenNotFound):
		return ErrHandshakeBadHeaderUpgrade

	case errors.Is(nbioErr, websocket.ErrUpgradeMethodIsGet):
		return ErrHandshakeBadMethod

	case errors.Is(nbioErr, websocket.ErrUpgradeInvalidWebsocketVersion):
		return ErrHandshakeBadHeaderSecVersion

	case errors.Is(nbioErr, websocket.ErrUpgradeOriginNotAllowed):
		return ErrHandshakeBadOrigin

	case errors.Is(nbioErr, websocket.ErrUpgradeMissingWebsocketKey):
		return ErrHandshakeBadHeaderSecKey

	default:
		return nbioErr
	}
}

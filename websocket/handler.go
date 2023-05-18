package ekaweb_socket

type Handler interface {
	OnOpen(c Conn) error

	// OnMessage registers a callback that will be called when a new data message
	// is received from the WebSocket client.
	// The callback should return an error if connection should be closed.
	// You wouldn't catch close message by this handler. Use OnClose() instead.
	//
	// Calling more than once or calling in any callback but Handler is prohibited
	// and may lead to UB.
	OnMessage(c Conn, typ MessageType, payload []byte) error

	// OnClose registers a callback that will be called when a new close message
	// is received from the WebSocket client.
	// The callback should return yes if it wants the same close message be sent
	// to the client after callback is done. False should be used otherwise.
	// Being inside that callback you should consider that connection (logically)
	// is closed already and you should use only CloseWithCode() method inside.
	// Calling CloseWithCode() won't trigger this callback.
	//
	// Calling more than once or calling in any callback but Handler is prohibited
	// and may lead to UB.
	OnClose(c Conn, cc CloseCode, detail string) bool
}

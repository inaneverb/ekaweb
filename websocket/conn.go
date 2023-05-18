package ekaweb_socket

import (
	"context"
)

// Conn is an interface that represents WebSocket client.
type Conn interface {
	ID() string

	Context() context.Context

	// WriteMessage writes a new message to the Client.
	WriteMessage(typ MessageType, payload []byte)

	// CloseWithCode writes a new close message to the Client and then closes
	// the connection. You must NOT use the connection after usage of this method.
	// Calling this method won't trigger the OnClose() callback.
	CloseWithCode(cc CloseCode)
}

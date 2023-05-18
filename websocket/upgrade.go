package ekaweb_socket

import (
	"errors"
)

var (
	ErrHandshakeBadMethod           = errors.New("Extension.WebSocket: Bad HTTP method (must be GET)")
	ErrHandshakeBadProtocol         = errors.New("Extension.WebSocket: Bad HTTP protocol (must be 1.1 or greater)")
	ErrHandshakeBadHeaderUpgrade    = errors.New("Extension.WebSocket: Bad HTTP Upgrade header (must be 'websocket')")
	ErrHandshakeBadHeaderConnection = errors.New("Extension.WebSocket: Bad HTTP Connection header (must contain 'Upgrade')")
	ErrHandshakeBadHeaderSecKey     = errors.New("Extension.WebSocket: Bad HTTP Sec-Websocket-Key header (must have valid length)")
	ErrHandshakeBadHeaderSecVersion = errors.New("Extension.WebSocket: Bad HTTP Sec-Websocket-Version header (must be '13')")
	ErrHandshakeBadOrigin           = errors.New("Extension.WebSocket: Bad HTTP origin (maybe need less strict policy?)")
)

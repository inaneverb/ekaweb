package ekaweb_nbiows

import (
	"github.com/lesismal/nbio/nbhttp/websocket"
)

func MessageTypeFromNbio(nt websocket.MessageType) MessageType {
	switch nt {
	case websocket.TextMessage:
		return MessageTypeDataText
	case websocket.BinaryMessage:
		return MessageTypeDataBinary
	case websocket.CloseMessage:
		return MessageTypeControlClose
	case websocket.PingMessage:
		return MessageTypeControlPing
	case websocket.PongMessage:
		return MessageTypeControlPong
	default:
		return MessageTypeControlInvalid
	}
}

func MessageTypeToNbio(t MessageType) websocket.MessageType {
	switch t {
	case MessageTypeDataText:
		return websocket.TextMessage
	case MessageTypeDataBinary:
		return websocket.BinaryMessage
	case MessageTypeControlClose:
		return websocket.CloseMessage
	case MessageTypeControlPing:
		return websocket.PingMessage
	case MessageTypeControlPong:
		return websocket.PongMessage
	default:
		return websocket.BinaryMessage
	}
}

package ekaweb_socket

// MessageType is a type that describes what kind of some Websocket message is.
//
// If some function returns a binary data with this value, the data is expected
// to be a payload of correspondent websocket frame.
type MessageType uint8

const (
	// All variants, their values, reserved parts and other are taken from:
	// https://datatracker.ietf.org/doc/html/rfc6455#section-5.2
	// https://www.iana.org/assignments/websocket/websocket.xhtml

	// ------------------------- NON CONTROL FRAMES ------------------------- //

	MessageTypeContinuation MessageType = 0x00
	MessageTypeDataText     MessageType = 0x01
	MessageTypeDataBinary   MessageType = 0x02
	_                       MessageType = 0x03 // reserved for further frames
	_                       MessageType = 0x04 // reserved for further frames
	_                       MessageType = 0x05 // reserved for further frames
	_                       MessageType = 0x06 // reserved for further frames
	_                       MessageType = 0x07 // reserved for further frames

	// --------------------------- CONTROL FRAMES --------------------------- //

	MessageTypeControlClose MessageType = 0x08
	MessageTypeControlPing  MessageType = 0x09
	MessageTypeControlPong  MessageType = 0x0A
	_                       MessageType = 0x0B // reserved for further frames
	_                       MessageType = 0x0C // reserved for further frames
	_                       MessageType = 0x0D // reserved for further frames
	_                       MessageType = 0x0E // reserved for further frames
	MessageTypeInvalid      MessageType = 0x0F // reserved, but who cares?
)

// String returns text representation of the MessageType. Useful for debugging.
func (mt MessageType) String() string {
	switch mt {
	case MessageTypeContinuation:
		return "Continuation"
	case MessageTypeDataText:
		return "Text"
	case MessageTypeDataBinary:
		return "Binary"
	case MessageTypeControlClose:
		return "Close"
	case MessageTypeControlPing:
		return "Ping"
	case MessageTypeControlPong:
		return "Pong"
	case MessageTypeInvalid:
		return "Invalid (non part of RFC6455)"
	default:
		return "Unknown"
	}
}

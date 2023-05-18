package ekaweb_socket

// CloseCode is a type that describes why websocket connection is closing.
//
// Note.
// You may use String() to get a description of the CloseCode and even send it
// to the client. If you need to register custom close code and its description,
// you should use RegisterCustomCloseCode() function.
type CloseCode uint16

const (
	// All variants, their values, reserved parts and other are taken from:
	// https://datatracker.ietf.org/doc/html/rfc6455#section-7.4.1
	// https://www.iana.org/assignments/websocket/websocket.xhtml

	// ---------------------- PART OF RFC6455 STANDARD ---------------------- //

	CloseCodeNormal           CloseCode = 1000
	CloseCodeGoingAway        CloseCode = 1001
	_                         CloseCode = 1002 // https://stackoverflow.com/questions/22587438
	CloseCodeUnsupportedData  CloseCode = 1003
	_                         CloseCode = 1004 // Reserved
	CloseCodeNoStatusReceived CloseCode = 1005 // MUST NOT BE SENT TO THE CLIENT!
	CloseCodeAbnormal         CloseCode = 1006 // MUST NOT BE SENT TO THE CLIENT!
	CloseCodeInvalidPayload   CloseCode = 1007
	CloseCodePolicyViolation  CloseCode = 1008
	CloseCodeMessageTooBig    CloseCode = 1009
	CloseCodeMandatoryExt     CloseCode = 1010 // MUST NOT BE SENT TO THE CLIENT!
	CloseCodeInternalError    CloseCode = 1011
	CloseCodeTLSHandshake     CloseCode = 1015 // MUST NOT BE SENT TO THE CLIENT!

	// -------------------- NON PART OF RFC6455 STANDARD -------------------- //

	CloseCodeServiceRestart CloseCode = 1012 // http://www.ietf.org/mail-archive/web/hybi/current/msg09670.html
	CloseCodeTryAgainLater  CloseCode = 1013 // http://www.ietf.org/mail-archive/web/hybi/current/msg09670.html
	CloseCodeBadGateway     CloseCode = 1014 // http://www.ietf.org/mail-archive/web/hybi/current/msg10748.html

	CloseCodeUnauthorized CloseCode = 3000

	// ------------------------ APPLICATION DEFINED ------------------------- //

	CloseCodeApplicationDefinedMin CloseCode = 4000
	CloseCodeApplicationDefinedMax CloseCode = 4999
)

var (
	// Keeps the user registered CloseCode and their descriptions.
	//
	// Thread unsafety is made by design and because of performance requirements.
	// In 100% cases you firstly need to register all custom codes and then use them.
	// You will likely never need to update CloseCode description dynamically
	// at the runtime of your workers. If so it's your headache, not mine.
	customCloseCodes = map[CloseCode]string{
		CloseCodePolicyViolation: "ClosePolicyViolation (You should not have done that)",
		CloseCodeUnauthorized:    "CloseUnauthorized (Check the credentials and access)",
	}
)

// String returns a description (detail) of the current CloseCode.
// It may be used along with the value of CloseCode to be sent to the client.
//
// If you want to get a description of user defined CloseCode, it must be registered
// using RegisterCustomCloseCode() or empty string is returned otherwise.
func (cc CloseCode) String() string {
	switch cc {
	case CloseCodeNormal:
		return "CloseNormal (Goodbye)"
	case CloseCodeGoingAway:
		return "CloseGoingAway (Server is going shutdown)"
	case CloseCodeUnsupportedData:
		return "CloseUnsupportedData (Unsupported WebSocket message)"
	//case CloseCodeNoStatusReceived:
	//	return "" // Read RFC6455 why this code should not be used
	//case CloseCodeAbnormal:
	//	return "" // Read RFC6455 why this code should not be used
	case CloseCodeInvalidPayload:
		return "CloseInvalidFrame (Decode error or illegal internal data)"
	case CloseCodeMessageTooBig:
		return "CloseTooBig (Message is too big)"
	//case CloseCodeMandatoryExt:
	//	return "" // Read RFC6455 why this code should not be used
	case CloseCodeInternalError:
		return "CloseInternalError (Unrecoverable)"
	//case CloseCodeTLSHandshake:
	//	return "" // Read RFC6455 why this code should not be used
	case CloseCodeServiceRestart:
		return "CloseServiceRestart (Try re-connect in [5..30] seconds)"
	case CloseCodeTryAgainLater:
		return "CloseTryAgainLater (There is no available resources or you exceeded your limit)"
	case CloseCodeBadGateway:
		return "CloseBadGateway (Upstream WebSocket service is unavailable)"
	default:
		return customCloseCodes[cc]
	}
}

// IsAllowedForTransmission reports whether current CloseCode can be send
// from the server to the client and vice-versa according with RFC6455.
func (cc CloseCode) IsAllowedForTransmission() bool {
	switch cc {
	case 1002, 1004, CloseCodeNoStatusReceived, CloseCodeAbnormal, CloseCodeTLSHandshake:
		return false
	default:
		return true
	}
}

// RegisterCustomCloseCode allows you to register custom (user defined) CloseCode
// and specify its description that will be available using CloseCode.String().
//
// WARNING! NOT THREAD-SAFETY! DATA RACE OTHERWISE!
// This function is setter and CloseCode.String() is getter. You MUST NOT use them
// concurrently. Call setter to set all custom codes you want to use and only then
// run your workers that can extract descriptions.
// You also cannot use 2 or more concurrent "write" goroutines to set custom CloseCode
// but you may use any number of concurrent "read" goroutines to call CloseCode.String().
//
// WARNING! ONLY USER-DEFINED CLOSE CODE DESCRIPTIONS MAY BE OVERWRITTEN!
// You cannot overwrite (specify) descriptions for CloseCode which are not
// in the range [4000..4999], except 1008, 3000. Calling this func for other CloseCode
// will do nothing.
func RegisterCustomCloseCode(cc CloseCode, detail string) {
	if detail != "" &&
		(cc == CloseCodePolicyViolation || cc == CloseCodeUnauthorized ||
			(cc >= CloseCodeApplicationDefinedMin && cc <= CloseCodeApplicationDefinedMax)) {

		customCloseCodes[cc] = detail
	}
}

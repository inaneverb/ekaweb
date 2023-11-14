package ekaweb_jrpc

import (
	"fmt"
)

// ResponseError represents a server-defined jRPC error. It should contain
// a meaningful info about occurred error, that will be sent to the client.
// Read more: https://www.jsonrpc.org/specification#error_object
type ResponseError struct {
	Code    int    `json:"code"`           // -32000 to -32099
	Message string `json:"message"`        // single sentence
	Data    any    `json:"data,omitempty"` // detailed info about occurred error
}

var (
	// ErrRequestMalformed is when given jRPC request is not valid JSON,
	// has no jRPC method or considered incorrect because of other reason.
	//
	// WARNING! Use deep error check (using errors.Is()) to check error
	// against this one, because it can be wrapped to more detailed one.
	ErrRequestMalformed = fmt.Errorf("jRPC: malformed request")

	// ErrMethodNotRegistered is when jRPC request want to call the method
	// that is not registered in jRPC router.
	// NOTE. errors.Is() is recommended, but just equality check is fine.
	ErrMethodNotRegistered = fmt.Errorf("jRPC: method not registered")
)

// FillMissedFields fills fields that are required by jRPC standard,
// but are missed in current ResponseError.
func (re *ResponseError) FillMissedFields() {
	if re.Code == 0 {
		re.Code = -32000
	}
	if re.Message == "" {
		re.Message = "internal server error"
	}
}

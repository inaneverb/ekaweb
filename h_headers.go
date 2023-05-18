package ekaweb

import (
	"net/http"
)

func CopyHeaders(to, from http.Header) http.Header {
	if to == nil || from == nil {
		return to
	}

	for fromHeaderKey, fromHeaderValues := range from {
		for _, fromHeaderValue := range fromHeaderValues {
			to.Add(fromHeaderKey, fromHeaderValue)
		}
	}

	return to
}

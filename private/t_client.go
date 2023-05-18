package ekaweb_private

import (
	"context"
	"net/http"
)

type ClientRequest interface {
	Data() ([]byte, error)
	ContentType() string
}

type ClientResponse interface {
	FromData(statusCode int, data []byte) error
}

type Client interface {
	Do(ctx context.Context, method, path string, headers http.Header,
		req ClientRequest, resp ClientResponse) error
}

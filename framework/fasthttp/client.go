package ekaweb_fasthttp

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/valyala/fasthttp"

	"github.com/inaneverb/ekacore/ekaunsafe/v4"
	"github.com/inaneverb/ekaweb/v2"
	"github.com/inaneverb/ekaweb/v2/private"
)

type Client struct {
	origin *fasthttp.Client
	log    ekaweb.Logger
	path   string
}

////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////

func (c *Client) Do(
	ctx context.Context, method, path string, headers http.Header,
	req ekaweb.ClientRequest, resp ekaweb.ClientResponse) error {

	// Perform early encoding.
	// It allows us to skip performing operations if encoding is failed.

	var data []byte
	var err error

	if req != nil {
		if data, err = req.Data(); err != nil {
			//c.e(err, "Failed to encode request.", method, nil, nil, nil)
			return fmt.Errorf("failed to encode request: %w", err)
		}
	}

	// Allocate necessary objects, fill them.

	var fhUri = fasthttp.AcquireURI()
	defer fasthttp.ReleaseURI(fhUri)

	var fhReq = fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(fhReq)

	var fhResp = fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(fhResp)

	path = strings.Trim(path, "\\/")
	if c.path != "" {
		path = c.path + "/" + path
	}

	if err = fhUri.Parse(nil, ekaunsafe.StringToBytes(path)); err != nil {
		//c.e(err, "Failed to parse URL.", method, nil, nil, nil)
		return fmt.Errorf("failed to parse URL (%s): %w", path, err)
	}

	// Apply request's data if it was generated.
	// Apply it as query parameters for GET & DELETE methods;
	// apply as HTTP body for other methods.

	switch {
	case len(data) == 0:
		// Skip setting request data.

	case method == ekaweb.MethodGet || method == ekaweb.MethodDelete:
		fhUri.SetQueryStringBytes(data)

	default:
		if mimeType := req.ContentType(); mimeType != "" {
			fhReq.Header.SetContentType(mimeType)
		}
		fhReq.SetBody(data)
	}

	for headerKey, headerValue := range headers {
		switch len(headerValue) {
		case 0:
		case 1:
			fhReq.Header.Set(headerKey, headerValue[0])
		default:
			for i, n := 0, len(headerValue); i < n; i++ {
				fhReq.Header.Add(headerKey, headerValue[i])
			}
		}
	}

	fhReq.Header.SetMethod(method)
	fhReq.SetRequestURI(fhUri.String())

	// Ok, we're ready to perform HTTP request.

	if deadLine, ok := ctx.Deadline(); ok && !deadLine.IsZero() {
		err = c.origin.DoDeadline(fhReq, fhResp, deadLine)
	} else {
		err = c.origin.Do(fhReq, fhResp)
	}

	if err != nil {
		//c.e(err, "", method, fhUri, fhReq, nil)
		return fmt.Errorf("failed to perform HTTP request: %w", err)
	}

	// Analyze and decode response.

	var statusCode = fhResp.StatusCode()
	var isOK = statusCode >= 200 && statusCode <= 299
	var respBody = fhResp.Body()

	switch {
	case !isOK && resp == nil:
		const E = "HTTP status code is %d, but response is not declared"
		err = fmt.Errorf(E, statusCode)

	case resp != nil:
		err = resp.FromData(statusCode, respBody)
	}

	{
		//const D = "Failed to decode or analyze response."
		//c.e(err, D, method, fhUri, fhReq, respBody)
	}

	return err
}

//func (c *Client) e(
//	err error, description, method string,
//	uri *fasthttp.URI, req *fasthttp.Request, respBody []byte) {
//
//	const s0 = "HTTP request successfully completed."
//	const s1 = "Failed to perform HTTP request using client."
//
//	if c.log == nil {
//		return
//	}
//
//	var descField zap.Field = zap.Skip()
//	var uriField zap.Field = zap.Skip()
//	var reqField zap.Field = zap.Skip()
//	var respField zap.Field = zap.Skip()
//
//	if err != nil && description != "" {
//		descField = zap.String("description", description)
//	}
//
//	if uri != nil {
//		uriField = zap.String("req_uri", method+" "+uri.String())
//	}
//
//	if req != nil {
//		reqField = zap.Stringer("req_dump", req)
//	}
//
//	switch {
//	case len(respBody) > 64:
//		respBody = respBody[:64]
//		fallthrough
//
//	case len(respBody) > 0:
//		respField = zap.String("resp_part_dump", ekaunsafe.BytesToString(respBody))
//	}
//
//	var f = helpers.If(err == nil, c.log.Debug, c.log.Error)
//	var s = helpers.If(err == nil, s0, s1)
//
//	f(s, zap.Error(err), descField, uriField, reqField, respField)
//}

////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////

func NewClient(options ...ekaweb.ClientOption) ekaweb.Client {

	var client Client
	client.origin = new(fasthttp.Client)

	client.origin.MaxIdleConnDuration = 30 * time.Second

	for i, n := 0, len(options); i < n; i++ {
		if ekaunsafe.UnpackInterface(options[i]).Word == nil {
			continue
		}

		switch option := options[i].(type) {

		case *ekaweb_private.ClientOptionHostAddr:
			var _, err = url.Parse(strings.Trim(option.Addr, "/\\ "))
			if err == nil {
				client.path = option.Addr
			}

		case *ekaweb_private.ClientOptionUserAgent:
			client.origin.Name = option.UserAgent
			client.origin.NoDefaultUserAgentHeader = option.UserAgent == ""

		case *ekaweb_private.ClientServerOptionTimeout:
			if option.ReadTimeout > 0 {
				client.origin.ReadTimeout = option.ReadTimeout
			}
			if option.WriteTimeout > 0 {
				client.origin.WriteTimeout = option.WriteTimeout
			}

		case *ekaweb_private.ClientServerOptionLogger:
			if option.Log != nil {
				client.log = option.Log
			}
		}
	}

	return &client
}

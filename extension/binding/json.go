// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package ekaweb_bind

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/goccy/go-json"

	"github.com/inaneverb/ekaweb/private"
)

type jsonBinding struct{}

func (jsonBinding) Name() string {
	return "JSON"
}

func (jsonBinding) Bind(req *http.Request, obj any) error {
	if req == nil || req.Body == nil {
		return fmt.Errorf("invalid request")
	}
	return decodeJSON(req.Context(), req.Body, obj)
}

func (jsonBinding) BindBody(body []byte, obj any) error {
	return decodeJSON(nil, bytes.NewReader(body), obj)
}

func decodeJSON(ctx context.Context, r io.Reader, obj any) error {

	var data, err = io.ReadAll(r)
	if err != nil && err != io.EOF {
		return err
	}

	if err = ekaweb_private.DecodeJSON(ctx, data, obj); err != nil {
		return err
	}

	var jsonEncDec *ekaweb_private.RouterOptionCustomJSON = nil
	if ctx != nil {
		jsonEncDec = ekaweb_private.UkvsGetJSONEncoderDecoder(ctx)
	}

	if jsonEncDec != nil && jsonEncDec.Decoder != nil {
		err = jsonEncDec.Decoder(data, obj)
	} else {
		err = json.Unmarshal(data, obj)
	}

	if err != nil && err != io.EOF {
		return err
	}

	return validate(obj)
}

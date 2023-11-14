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

	"github.com/inaneverb/ekaweb/v2/private"
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
	var err = ekaweb_private.DecodeStream(ctx, r, obj)
	if err != nil {
		return err
	}
	return validate(obj)
}

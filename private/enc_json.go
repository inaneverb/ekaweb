package ekaweb_private

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
)

// EncodeJSON ...
// Deprecated: Use EncodeStream() instead.
func EncodeJSON(ctx context.Context, v any) ([]byte, error) {

	var b = bytes.NewBuffer(nil)
	var err = EncodeStream(ctx, b, v)

	if err != nil {
		return nil, err
	}

	return b.Bytes(), err
}

func EncodeStream(ctx context.Context, to io.Writer, v any) error {
	if encoderGetter := UkvsGetCodec(ctx).EncoderGetter; encoderGetter != nil {
		return encoderGetter(to).Encode(v)
	} else {
		return json.NewEncoder(to).Encode(v)
	}
}

// DecodeJSON ...
// Deprecated: Use DecodeStream() instead.
func DecodeJSON(ctx context.Context, data []byte, dest any) error {
	return DecodeStream(ctx, bytes.NewReader(data), dest)
}

func DecodeStream(ctx context.Context, r io.Reader, to any) error {
	var decoderGetter DecoderGetter
	if ctx != nil {
		decoderGetter = UkvsGetCodec(ctx).DecoderGetter
	}

	var err error
	if decoderGetter != nil {
		err = decoderGetter(r).Decode(to)
	} else {
		err = json.NewDecoder(r).Decode(to)
	}

	if err == io.EOF {
		err = nil
	}

	return err
}

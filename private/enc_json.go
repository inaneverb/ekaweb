package ekaweb_private

import (
	"context"
	"io"

	"github.com/goccy/go-json"
)

func EncodeJSON(ctx context.Context, v any) ([]byte, error) {
	var jsonEncDec = UkvsGetJSONEncoderDecoder(ctx)
	if jsonEncDec != nil && jsonEncDec.Encoder != nil {
		return jsonEncDec.Encoder(v)
	}
	return json.Marshal(v)
}

func DecodeJSON(ctx context.Context, data []byte, dest any) error {

	var jsonEncDec *RouterOptionCustomJSON = nil
	if ctx != nil {
		jsonEncDec = UkvsGetJSONEncoderDecoder(ctx)
	}

	var err error
	if jsonEncDec != nil && jsonEncDec.Decoder != nil {
		err = jsonEncDec.Decoder(data, dest)
	}
	err = json.Unmarshal(data, dest)

	if err == io.EOF {
		err = nil
	}

	return err
}

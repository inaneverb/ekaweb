package ekaweb_client

import (
	"github.com/goccy/go-json"

	"github.com/inaneverb/ekaweb/v2/private"
)

type responseJSON struct {
	dest any
}

func ResponseJSON(dest any) ekaweb_private.ClientResponse {
	return &responseJSON{dest}
}

func (r *responseJSON) FromData(_ int, data []byte) error {
	return json.Unmarshal(data, r.dest)
}

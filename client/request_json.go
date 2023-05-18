package ekaweb_client

import (
	"github.com/goccy/go-json"

	"github.com/inaneverb/ekaweb"
	"github.com/inaneverb/ekaweb/private"
)

type requestJSON struct {
	source any
}

func RequestJSON(source any) ekaweb_private.ClientRequest {
	return &requestJSON{source}
}

func (w *requestJSON) Data() ([]byte, error) {
	return json.Marshal(w.source)
}

func (w *requestJSON) ContentType() string {
	return ekaweb.MIMEApplicationJSON
}

package ekaweb_respondent

import (
	"fmt"
	"net/http"

	"github.com/inaneverb/ekaweb"
	"github.com/inaneverb/ekaweb/private"
)

// CommonApplicator is an error "applicator". Implements Applicator interface.
// The core of the logic is Apply() method. It allows you to use some Manifest
// as an HTTP response.
//
// It's ready-to-use object after instantiating this type but feel free to use
// its constructor: NewCommonApplicator().
type CommonApplicator struct{}

// Apply tries to encode provided Manifest using JSON encoder and then use it
// as an HTTP response writing Manifest's status code, correspondent HTTP headers
// and the encoded JSON data as an HTTP body to the socket.
//
// Because it utilizes provided types.Context to make an HTTP response
// it propagates an error of any context's method that is used
// while sending a response.
func (*CommonApplicator) Apply(ctx any, manifest *Manifest) {

	//goland:noinspection GoVetStructTag
	type ManifestJSON struct {
		Status        int                    `json:"-"`
		Error         string                 `json:"error"`
		ErrorID       string                 `json:"error_id,omitempty"`
		ErrorCode     int                    `json:"error_code"`
		ErrorDetail   string                 `json:"error_detail,omitempty"`
		ErrorDetails  []string               `json:"error_details,omitempty"`
		customFillers []ManifestCustomFiller `json:"-"`
	}

	var w http.ResponseWriter
	var r *http.Request

	if httpCtx, ok := ctx.(HttpContext); ok {
		w, r = httpCtx.W, httpCtx.R
	} else {
		return
	}

	// Connection may be hijacked (WebSocket, for example).
	// In that case it should be handled anyhow else.

	if ekaweb_private.UkvsIsConnectionHijacked(r.Context()) {
		return
	}

	// Since go1.16 we can do this type conversion without unsafe.
	// Moreover, without unsafe this check also guarantees
	// that structures are the same.

	if len(manifest.customFillers) != 0 {
		manifest = manifest.Clone()
	}

	for i, n := 0, len(manifest.customFillers); i < n; i++ {
		manifest.customFillers[i](r, manifest)
	}

	jsonManifest := (*ManifestJSON)(manifest)
	ekaweb.SendJSON(w, r, jsonManifest.Status, jsonManifest)
}

// NewCommonApplicator is a constructor of CommonApplicator. It initializes all
// internal fields and returns a ready-to-use object.
func NewCommonApplicator() *CommonApplicator {
	return new(CommonApplicator)
}

////////////////////////////////////////////////////////////////////////////////
///// PUBLIC HELPERS ///////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////

// CACF_AutoErrorDetail is a ManifestCustomFiller that extracts error detail
// from the HTTP context being a part of CommonApplicator. You can register
// this function for CommonExpander for auto-filling an error detail field.
//
//goland:noinspection GoSnakeCaseUsage
func CACF_AutoErrorDetail(ctx any, manifest *Manifest) {
	if r, ok := ctx.(*http.Request); ok {
		manifest.ErrorDetail = ekaweb_private.UkvsGetUserErrorDetail(r.Context())
	} else {
		fmt.Println("ekaweb.CACF_AutoErrorDetail(): Not in a HTTP context")
	}
}

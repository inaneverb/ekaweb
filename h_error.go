package ekaweb

import (
	"net/http"

	"github.com/inaneverb/ekaweb/private"
)

func ErrorGet(r *http.Request) error {
	return ekaweb_private.UkvsGetUserError(r.Context())
}

func ErrorApply(r *http.Request, err error) {
	ekaweb_private.UkvsInsertUserError(r.Context(), err)
}

func ErrorDetailGet(r *http.Request) string {
	return ekaweb_private.UkvsGetUserErrorDetail(r.Context())
}

func ErrorDetailApply(r *http.Request, detail string) {
	ekaweb_private.UkvsInsertUserErrorDetail(r.Context(), detail)
}

func ErrorRemove(r *http.Request) {
	ekaweb_private.UkvsRemoveUserErrorFull(r.Context())
}

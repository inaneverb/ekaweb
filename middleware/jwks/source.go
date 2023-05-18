package ekaweb_jwks

import (
	"context"

	"github.com/lestrrat-go/jwx/jwk"
)

// Source is an interface that provides a way to get a jwk.Set instance.
type Source interface {
	GetJwks() (jwk.Set, error)
}

// SourceCommon implements Source interface.
// It is instantiated by jwk.Set or by jwk.Key (or the resource that will be used
// to built a jwk.Key) based on which a new jwk.Set will be created.
//
// The GetJwks() method has O(1) constant access.
type SourceCommon struct {
	s jwk.Set
}

// SourceAutoRefresh implements source interface.
// It is more complicated way to get a jwk.Set. It contains jwk.AutoRefresh
// object and the arguments to call its methods.
type SourceAutoRefresh struct {
	ar   *jwk.AutoRefresh
	addr string
	ctx  context.Context
}

// GetJwks just returns underlying jwk.Set. Returned error is always nil.
func (sc SourceCommon) GetJwks() (jwk.Set, error) {
	return sc.s, nil
}

// GetJwks calls Fetch() of underlying jwk.AutoRefresh and returns a jwk.Set.
// It's underlying jwk.AutoRefresh responsibility to answer the question
// how fresh returned jwk.Set is.
//
// Technically, if all refresh calls are failed there's only one first version
// of jwk.Set is exist. Under long time executed services it may lead
// to unexpected behaviour.
func (sar *SourceAutoRefresh) GetJwks() (jwk.Set, error) {
	return sar.ar.Fetch(sar.ctx, sar.addr)
}

// NewSourceCommon returns a new initialized SourceCommon.
func NewSourceCommon(s jwk.Set) (Source, error) {
	return SourceCommon{s}, nil
}

// NewSourceAutoRefresh returns a new initialized SourceAutoRefresh.
//
// If delay is not provided, the HTTP headers of response will be used
// to calculate a new time of jwk.Set refreshing.
// If those headers are not provided, default refresh time is 1h.
func NewSourceAutoRefresh(
	ctx context.Context, url string, errChan chan jwk.AutoRefreshError,
	options ...jwk.AutoRefreshOption,
) (Source, error) {

	ar := jwk.NewAutoRefresh(ctx)

	if errChan != nil {
		ar.ErrorSink(errChan)
	}

	ar.Configure(url, options...)

	if _, err := ar.Refresh(ctx, url); err != nil {
		return nil, err
	}

	return &SourceAutoRefresh{ar, url, ctx}, nil
}

package ekaweb_sbac

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/inaneverb/ekaweb"
	"github.com/inaneverb/ekaweb/private"

	"github.com/inaneverb/ekaweb/middleware/jwks"
)

var (
	// ErrScopesNotFound is returned when no scopes are present in somewhere:
	// neither JWT nor HTTP context after or the set is empty.
	ErrScopesNotFound = errors.New("Middleware.SBAC: Scopes not found")

	// ErrAccessDenied is returned when user doesn't have required scopes.
	ErrAccessDenied = errors.New("Middleware.SBAC: Access denied")
)

////////////////////////////////////////////////////////////////////////////////
///// Extract middleware (extracts scopes from JWT) ////////////////////////////
////////////////////////////////////////////////////////////////////////////////

// SearchScopesInJwt is a middleware generator. Returns a middleware that will
// extract JWT from the HTTP context, find JWT field with the name "scopes"
// (or with the name you pass in this generator if non-empty)
// and then saves that scopes to the HTTP context.
//
// Thus returned middleware DO NOT performs scopes validation. It only extracts
// them and saves to the HTTP context.

type ExtractMiddleware interface {
	ekaweb.Middleware
	ManualByContext(ctx context.Context)
	Manual(token ekaweb_jwks.Token) ([]string, error)
}

type extractMiddleware struct {
	JwtClaimKey string
}

func (m *extractMiddleware) Callback(next ekaweb.Handler) ekaweb.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		m.ManualByContext(r.Context())
		next.ServeHTTP(w, r)
	})
}

func (m *extractMiddleware) ManualByContext(ctx context.Context) {

	scopes, err := m.Manual(ekaweb_jwks.GetTokenByContext(ctx))
	if err != nil {
		ekaweb_private.UkvsInsertUserError(ctx, err)
		return
	}

	SetScopesByContext(ctx, scopes)
}

func (m *extractMiddleware) Manual(token ekaweb_jwks.Token) ([]string, error) {

	if token == nil {
		return nil, ErrScopesNotFound
	}

	scopesUntyped, found := token.Get(m.JwtClaimKey)
	if !found || scopesUntyped == nil {
		return nil, ErrScopesNotFound
	}

	scopes, ok := scopesUntyped.([]string)
	if !ok || len(scopes) == 0 {
		return nil, ErrScopesNotFound
	}

	return scopes, nil
}

func NewExtractMiddleware(jwtClaimKey string) ExtractMiddleware {
	const DefaultJwtClaimKey = "scopes"

	if jwtClaimKey = strings.TrimSpace(jwtClaimKey); jwtClaimKey == "" {
		jwtClaimKey = DefaultJwtClaimKey
	}

	ekaweb_jwks.RegisterCustomField(jwtClaimKey, []string{"_"})
	return &extractMiddleware{jwtClaimKey}
}

////////////////////////////////////////////////////////////////////////////////
///// Check middleware (makes sure valid scopes are provided) //////////////////
////////////////////////////////////////////////////////////////////////////////

type CheckMiddleware interface {
	ekaweb.Middleware
	ManualByContext(ctx context.Context)
	Manual(providedScopes []string) error
}

type checkMiddleware struct {
	RequiredScopes []string
	NeedAll        bool
}

func (m *checkMiddleware) Callback(next ekaweb.Handler) ekaweb.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		m.ManualByContext(r.Context())
		next.ServeHTTP(w, r)
	})
}

func (m *checkMiddleware) ManualByContext(ctx context.Context) {
	err := m.Manual(GetScopesByContext(ctx))
	if err != nil {
		ekaweb_private.UkvsInsertUserError(ctx, err)
		return
	}
}

func (m *checkMiddleware) Manual(providedScopes []string) error {
	switch {
	case len(providedScopes) == 0:
		return ErrScopesNotFound

	case !TestScopes(providedScopes, m.RequiredScopes, m.NeedAll):
		return ErrAccessDenied

	default:
		return nil
	}
}

func NewCheckMiddleware(requiredScopes []string, needAll bool) CheckMiddleware {
	return &checkMiddleware{requiredScopes, needAll}
}

////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////

// TestScopes returns true if scope verification is passed.
// It's the main method of the middleware.
func TestScopes(presentedScopes, requiredScopes []string, needAll bool) bool {

	if r0, p0 := len(requiredScopes) == 0, len(presentedScopes) == 0; r0 || p0 {
		return r0
	}

	// There's kinda bool logic magic in the if-statement inside the loop
	// and in the return statement.
	//
	// Look.
	// If after inner loop, needAll != approvedScope it handles 2 cases:
	//
	// 1. At least one scope must be present (needAll == false) and the current
	//    required scope is found in the present ones (approvedScope == true);
	//
	// 2. All required scopes must be present (needAll == true) and the current
	//    required scope not found in the present ones (approvedScope == false);
	//
	// In all that cases there's no necessary to continue processing.
	// The result whether scope check is passed or not depends on the approvedScope.
	//
	// -----
	//
	// The outer last return statement also is simple.
	// If we're before that statement it means that early return didn't fired.
	// It's possible also in 2 cases:
	//
	// 1. At least one scope must be present (needAll == false) and no one from
	//    required scopes found in the presented (otherwise it would be early return);
	//
	// 2. All required scopes must be present (needAll == true) and
	//    there were no early return (all required scopes are approved).
	//
	// So, returned bool value (the scope check) depends on needAll value.

	for i, n := 0, len(requiredScopes); i < n; i++ {

		approvedScope := false
		for j, m := 0, len(presentedScopes); j < m && !approvedScope; j++ {
			approvedScope = requiredScopes[i] == presentedScopes[j]
		}

		if needAll != approvedScope {
			return approvedScope
		}
	}

	return needAll
}

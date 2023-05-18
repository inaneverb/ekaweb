package ekaweb_sbac

import (
	"context"
	"net/http"
)

/*
SBAC is Scopes based access control middleware.
It allows you to check whether an initiator of HTTP request has required scopes.

First of all you need to register middleware that is returned by SbacExtractFromJwt()
function. That middleware will extract scopes from the JWT that is stored
in the HTTP context as user's value.
Keep in mind that since that middleware extracts scopes from JWT token that is stored,
you need a middleware that will extract JWT from the HTTP request, parse, validate,
verify it and then store scopes to the context. The JWKS middleware can do that.

Then you need to prepend SbacRequireScopes() or SbacRequireScopesAny()
to the handler's list of your specific route that you want to guard
requesting user to have some scopes.

Also you can use SbacContainsScopes() or SbacContainsScopesAny() inside your HTTP
handlers (controllers) to perform scopes check.

You also can use SbacGetScopes() extractor to get low-level access to the scopes,
getting the list of scopes that are extracted from the JWT.
*/

////////////////////////////////////////////////////////////////////////////////
///// Middleware generators ////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////

// ExtractFromJwt generates a SBAC middleware
// that tries to extract a scopes from the JWT token and save it
// to the user's context for being able to validate these scopes later.
func ExtractFromJwt(customScopesKey ...string) ExtractMiddleware {

	scopesKey := ""
	if len(customScopesKey) > 0 {
		scopesKey = customScopesKey[0]
	}

	return NewExtractMiddleware(scopesKey)
}

// RequireScopes generates a middleware that retrieves extracted scopes
// (if any) from the user's scope and validates them. The scopes are stored
// in user's context MUST contain ALL scopes that are passed to this function.
func RequireScopes(scopes ...string) CheckMiddleware {
	return NewCheckMiddleware(scopes, true)
}

// RequireScopesAny is almost the same as SbacRequireScopes()
// but returned middleware ensures that at least any (not all)
// of required scope is present.
func RequireScopesAny(scopes ...string) CheckMiddleware {
	return NewCheckMiddleware(scopes, false)
}

////////////////////////////////////////////////////////////////////////////////
///// Extractors from the HTTP Request or context.Context //////////////////////
////////////////////////////////////////////////////////////////////////////////

// GetScopes is the same as SbacGetScopesByContext
// but extracts context.Context from http.Request by itself.
func GetScopes(r *http.Request) []string {
	return GetScopesByContext(r.Context())
}

// ContainsScopesByContext is a helper function that allows code to check
// whether some scope is present.
//
// It's useful when you have e.g. 5 scopes for 3 API methods but 1 API method
// in some cases require also a different scope(s).
// In that case you require scopes that are common for all API methods
// by the middleware and in the API handler (controller) checks whether on-demand
// scope is provided using this method.
//
// If required scopes are empty, the true is returned.
func ContainsScopesByContext(ctx context.Context, scopes ...string) bool {
	return TestScopes(GetScopesByContext(ctx), scopes, true)
}

// ContainsScopesAnyByContext is almost the same as SbacContainsScopesByContext()
// but the check ensures that at least one (not all) required scope is provided.
func ContainsScopesAnyByContext(ctx context.Context, scopes ...string) bool {
	return TestScopes(GetScopesByContext(ctx), scopes, false)
}

// ContainsScopes is the same as SbacContainsScopesByContext
// but extracts context.Context from http.Request by itself.
func ContainsScopes(r *http.Request, scopes ...string) bool {
	return TestScopes(GetScopesByContext(r.Context()), scopes, true)
}

// ContainsScopesAny is the same as SbacContainsScopesAnyByContext
// but extracts context.Context from http.Request by itself.
func ContainsScopesAny(r *http.Request, scopes ...string) bool {
	return TestScopes(GetScopesByContext(r.Context()), scopes, false)
}

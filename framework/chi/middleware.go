package ekaweb_chi

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/inaneverb/ekaweb"
	"github.com/inaneverb/ekaweb/private"
)

type MiddlewareCleanPathAndVariables struct{}

func (*MiddlewareCleanPathAndVariables) CheckErrorBefore() bool {
	return false
}

func (*MiddlewareCleanPathAndVariables) Callback(next ekaweb.Handler) ekaweb.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		ctx := r.Context()
		chiRouteContext := chi.RouteContext(ctx)

		path := chiRouteContext.RoutePattern()
		if path == "" {
			path = r.URL.Path
		}

		ekaweb_private.UkvsInsertOriginalPath(ctx, path)

		urlVariablesKeys := chiRouteContext.URLParams.Keys
		urlVariablesValues := chiRouteContext.URLParams.Values

		for i, n := 0, len(urlVariablesKeys); i < n; i++ {
			ekaweb_private.UkvsInsert(ctx, urlVariablesKeys[i], urlVariablesValues[i])
		}

		// Treat all paths as invalid ones and overwrite it
		// by the path specific middleware (registered for each path).
		//
		// Why such complexity?
		// Well, some underlying routers (not Chi) may handle 404/405
		// by the different way.
		// Or if I'd set it to the true by default everywhere,
		// and you will write adapters for your underlying router,
		// you will forget to add this thing and will spend a lot of time
		// thinking why you're getting "not found or not allowed"
		// when everything is ok. So, just thoughts about futured me or you.

		ekaweb_private.UkvsSetPathNotFoundOrNotAllowed(ctx, true)

		next.ServeHTTP(w, r)
	})
}

func NewCleanPathAndVariablesMiddleware() ekaweb.Middleware {
	return &MiddlewareCleanPathAndVariables{}
}

////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////

type MiddlewareInvalidatePath struct{}

func (*MiddlewareInvalidatePath) CheckErrorBefore() bool {
	return false
}

func (*MiddlewareInvalidatePath) Callback(next ekaweb.Handler) ekaweb.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ekaweb_private.UkvsSetPathNotFoundOrNotAllowed(r.Context(), false)
		next.ServeHTTP(w, r)
	})
}

func NewInvalidatePathMiddleware() ekaweb.Middleware {
	return &MiddlewareInvalidatePath{}
}

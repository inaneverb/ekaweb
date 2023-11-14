package ekaweb_sbac

import (
	"context"

	"github.com/inaneverb/ekaweb/v2/private"
)

type (
	ukvsKeyScopes struct{} // Ukvs key to store scopes ([]string)
)

// GetScopesByContext retrieves the scopes from the user's scope
// of the context.Context. If no scopes were set the empty array or nil
// is returned.
func GetScopesByContext(ctx context.Context) []string {
	return ekaweb_private.UkvsGetOrDefault(ctx, ukvsKeyScopes{}, []string(nil)).([]string)
}

// SetScopesByContext saves the presented scopes into the user's scope
// of context.Context. If scopes were set already it will be overwritten.
func SetScopesByContext(ctx context.Context, scopes []string) {
	ekaweb_private.UkvsInsert(ctx, ukvsKeyScopes{}, scopes)
}

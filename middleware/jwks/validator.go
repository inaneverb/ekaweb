package ekaweb_jwks

import (
	"context"
)

// TokenAdditionalValidator is a func alias that take Token and reports whether
// token valid or not (or any other error if any). It allows user
// to perform additional token check. User should return nil if token is valid.
type TokenAdditionalValidator = func(ctx context.Context, tok Token) error

// tokenAdditionalValidatorEmpty is an empty TokenAdditionalValidator.
// Just always returns nil as error.
func tokenAdditionalValidatorEmpty(_ context.Context, _ Token) error {
	return nil
}

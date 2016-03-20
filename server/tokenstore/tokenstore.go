package tokenstore

import (
	"github.com/clawio/service-auth/server/spec"
)

// TokenStore represents a store that can handle creation of tokens from an identity
// and also verify issued tokens.
// Implementations are: JWT
type TokenStore interface {
	Create(identity *spec.Identity) (string, error)
	Verify(token string) (*spec.Identity, error)
}

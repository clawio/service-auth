package tokenstore

import (
	"github.com/clawio/service-auth/server/spec"
)

type TokenStore interface {
	Create(identity *spec.Identity) (string, error)
	Verify(token string) (*spec.Identity, error)
}

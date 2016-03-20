package userstore

import "github.com/clawio/service-auth/server/spec"

// UserStore represents a store that can find users by their credentials.
// Implementations of UserStore are: SQL and LDAP.
type UserStore interface {
	FindByCredentials(username, password string) (*spec.Identity, error)
}
